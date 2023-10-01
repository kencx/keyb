package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kencx/keyb/config"
	"github.com/kencx/keyb/output"
	"github.com/kencx/keyb/ui"
)

const (
	help = `usage: keyb [options] <command>

  Options:
    -p, --print	    Print to stdout
    -e, --export    Export to file
    -k, --key       Key bindings at custom path
    -c, --config    Config file at custom path
    -v, --version   Version info
    -h, --help	    Show help

  Commands:
    a, add          Add keybind to keyb file
`

	addHelp = `usage: keyb [-k file] add [app; name; key]

  Options:
    -k, --key      Key bindings file at custom path
    -b, --binding  Key binding
    -p, --prefix   Ignore prefix
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

		addBind   string
		addPrefix bool
	)

	shortVersion := flag.Bool("v", false, "version information")
	longVersion := flag.Bool("version", false, "version information")

	flag.BoolVar(&stdout, "p", false, "print to stdout")
	flag.BoolVar(&stdout, "print", false, "print to stdout")

	flag.StringVar(&exportFile, "e", "", "export to file")
	flag.StringVar(&exportFile, "export", "", "export to file")

	flag.StringVar(&keybFile, "k", "", "keybindings file")
	flag.StringVar(&keybFile, "key", "", "keybindings file")

	flag.StringVar(&configFile, "c", "", "config file")
	flag.StringVar(&configFile, "config", "", "config file")

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addCmd.StringVar(&addBind, "b", "", "keybind")
	addCmd.StringVar(&addBind, "binding", "", "keybind")
	addCmd.BoolVar(&addPrefix, "p", false, "prefix")
	addCmd.BoolVar(&addPrefix, "prefix", false, "prefix")

	flag.Usage = func() { os.Stdout.Write([]byte(help)) }
	flag.Parse()

	if *shortVersion || *longVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	keys, cfg, err := config.Parse(configFile, keybFile)
	if err != nil {
		log.Fatal(err)
	}

	args := flag.Args()
	if len(args) > 1 {
		switch args[0] {
		case "add", "a":
			addCmd.Usage = func() { os.Stdout.Write([]byte(addHelp)) }
			addCmd.Parse(args[1:])

			var addFile string
			if keybFile != "" {
				// use flag -k path
				addFile = keybFile
			} else {
				// use default path in config
				addFile = cfg.KeybPath
			}
			if err := config.AddEntry(addFile, addBind, addPrefix); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s added to %s", addBind, addFile)
			os.Exit(0)
		default:
			fmt.Print(help)
			os.Exit(1)
		}
	}

	m := ui.NewModel(keys, cfg)

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

	p := tea.NewProgram(m, tea.WithMouseCellMotion(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to start: %w", err)
	}
	return nil
}
