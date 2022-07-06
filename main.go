package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kencx/keyb/config"
	"github.com/kencx/keyb/ui"
)

const help = `usage: keyb [options] [file]

  Flags:
    -p, --print	    Print to stdout
    -e, --export    Export to file
    -s, --strip     Strip ANSI chars (only for print/export)
    -k, --key       Key bindings at custom path
    -c, --config    Config file at custom path

    -h, --help	    help for keyb
`

//go:embed examples/config
var configFs string

// TODO support diff OS
var (
	defaultConfig = path.Join(os.Getenv("HOME"), ".config/keyb/config")
	defaultKeyb   = path.Join(os.Getenv("HOME"), ".config/keyb/keyb.yaml")
)

func main() {
	config.ConfigFs = configFs

	log.SetPrefix("keyb: ")
	log.SetFlags(0)

	var (
		strip      bool
		output     bool
		exportFile string
		keybFile   string
		configFile string
	)

	flag.BoolVar(&strip, "s", false, "strip ANSI chars")
	flag.BoolVar(&strip, "strip", false, "strip ANSI chars")
	flag.BoolVar(&output, "p", false, "print to stdout")
	flag.BoolVar(&output, "print", false, "print to stdout")
	flag.StringVar(&exportFile, "e", "", "export to file")
	flag.StringVar(&exportFile, "export", "", "export to file")
	flag.StringVar(&keybFile, "k", "", "keybindings file")
	flag.StringVar(&keybFile, "key", "", "keybindings file")
	flag.StringVar(&configFile, "c", defaultConfig, "config file")
	flag.StringVar(&configFile, "config", defaultConfig, "config file")

	flag.Usage = func() { os.Stdout.Write([]byte(help)) }
	flag.Parse()

	program, config, err := handleFlags(keybFile, configFile)
	if err != nil {
		log.Fatal(err)
	}
	m := ui.NewModel(program, config)

	if output {
		if err := OutputBodyToStdout(m, strip); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	if exportFile != "" {
		if err := OutputBodyToFile(m, exportFile, strip); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	if err := start(m); err != nil {
		log.Fatal(err)
	}
}

func handleFlags(keybPath, configPath string) (config.Bindings, *config.Config, error) {

	if err := config.CreateConfigFile(); err != nil {
		return nil, nil, fmt.Errorf("error: could not locate config file: %w", err)
	}
	cfg, err := config.GetConfig(configPath)
	if err != nil {
		return nil, nil, err
	}

	customKeybPath := cfg.KeybPath

	if keybPath != "" { // flag takes priority
		customKeybPath = keybPath
	} else if customKeybPath == "" && keybPath == "" {
		return nil, nil, fmt.Errorf("ERROR: no keyb file found")
	}

	bindings, err := config.GetBindings(customKeybPath)
	if err != nil {
		return nil, nil, err
	}
	return bindings, cfg, nil
}

func start(m *ui.Model) error {

	p := tea.NewProgram(m)
	p.EnableMouseCellMotion()
	p.EnterAltScreen()
	defer p.ExitAltScreen()

	if err := p.Start(); err != nil {
		return fmt.Errorf("failed to start program: %w", err)
	}
	return nil
}
