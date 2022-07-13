package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kencx/keyb/config"
	"github.com/kencx/keyb/output"
	"github.com/kencx/keyb/ui"
)

const (
	help = `usage: keyb [options] [file]

  Flags:
    -p, --print	    Print to stdout
    -e, --export    Export to file
    -k, --key       Key bindings at custom path
    -c, --config    Config file at custom path
    -v, --version   Version info
    -h, --help	    help for keyb
`
)

var version string

func main() {
	log.SetPrefix("keyb: ")
	log.SetFlags(log.Lshortfile)

	var (
		stdout     bool
		exportFile string
		keybFile   string
		configFile string
	)

	baseDir, err := config.GetBaseDir()
	if err != nil {
		log.Fatalf("os not supported: %v", err)
	}
	defaultConfig := path.Join(baseDir, "keyb", "config.yml")

	shortVersion := flag.Bool("v", false, "version information")
	longVersion := flag.Bool("version", false, "version information")
	flag.BoolVar(&stdout, "p", false, "print to stdout")
	flag.BoolVar(&stdout, "print", false, "print to stdout")
	flag.StringVar(&exportFile, "e", "", "export to file")
	flag.StringVar(&exportFile, "export", "", "export to file")
	flag.StringVar(&keybFile, "k", "", "keybindings file")
	flag.StringVar(&keybFile, "key", "", "keybindings file")
	flag.StringVar(&configFile, "c", defaultConfig, "config file")
	flag.StringVar(&configFile, "config", defaultConfig, "config file")

	flag.Usage = func() { os.Stdout.Write([]byte(help)) }
	flag.Parse()

	if *shortVersion || *longVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	keys, config, err := config.Parse(keybFile, configFile)
	if err != nil {
		log.Fatal(err)
	}

	m := ui.NewModel(keys, config)

	if stdout {
		if err := output.ToStdout(m); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}
	if exportFile != "" {
		if err := output.ToFile(m, exportFile); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	if err := start(m); err != nil {
		log.Fatal(err)
	}
}

func start(m *ui.Model) error {

	p := tea.NewProgram(m)
	p.EnableMouseCellMotion()

	if err := p.Start(); err != nil {
		return fmt.Errorf("failed to start: %w", err)
	}
	return nil
}
