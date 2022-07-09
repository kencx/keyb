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

const help = `usage: keyb [options] [file]

  Flags:
    -p, --print	    Print to stdout
    -e, --export    Export to file
    -k, --key       Key bindings at custom path
    -c, --config    Config file at custom path

    -h, --help	    help for keyb
`

// TODO support diff OS
var (
	defaultConfig = path.Join(os.Getenv("HOME"), ".config/keyb/config.yml")
	defaultKeyb   = path.Join(os.Getenv("HOME"), ".config/keyb/keyb.yml")
)

func main() {
	log.SetPrefix("keyb: ")
	log.SetFlags(log.Lshortfile)

	var (
		stdout     bool
		exportFile string
		keybFile   string
		configFile string
	)

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

	keys, config, err := parseFiles(keybFile, configFile)
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

func parseFiles(keyb, configPath string) (ui.Apps, *config.Config, error) {
	if err := config.CreateConfigFile(); err != nil {
		return nil, nil, fmt.Errorf("no config file found: %w", err)
	}

	cfg, err := config.Parse(configPath)
	if err != nil {
		return nil, nil, err
	}

	defaultKeyb := cfg.KeybPath
	if defaultKeyb == "" && keyb == "" {
		return nil, nil, fmt.Errorf("no keyb file found")
	}
	if keyb != "" {
		// overwrite default path with flag
		defaultKeyb = keyb
	}

	keys, err := ui.ParseApps(defaultKeyb)
	if err != nil {
		return nil, nil, err
	}
	return keys, cfg, nil
}

func start(m *ui.Model) error {

	p := tea.NewProgram(m)
	p.EnableMouseCellMotion()
	p.EnterAltScreen()
	defer p.ExitAltScreen()

	if err := p.Start(); err != nil {
		return fmt.Errorf("failed to start: %w", err)
	}
	return nil
}
