package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	tea "github.com/charmbracelet/bubbletea"
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

// TODO support diff OS
var (
	defaultConfig = path.Join(os.Getenv("HOME"), ".config/keyb/config")
	defaultKeyb   = path.Join(os.Getenv("HOME"), ".config/keyb/keyb.yaml")
)

type Categories map[string]Program

func main() {

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
	m := NewModel(program, config)

	if output {
		if err := m.OutputBodyToStdout(strip); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	if exportFile != "" {
		if err := m.OutputBodyToFile(exportFile, strip); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	if err := startProgram(m); err != nil {
		log.Fatal(err)
	}
}

func handleFlags(keybPath, configPath string) (Categories, *Config, error) {

	if err := createConfigFile(); err != nil {
		return nil, nil, fmt.Errorf("error: could not locate config file: %w", err)
	}
	config, err := GetConfig(configPath)
	if err != nil {
		return nil, nil, err
	}

	customKeybPath := config.KeybPath

	if keybPath != "" { // flag takes priority
		customKeybPath = keybPath
	} else if customKeybPath == "" && keybPath == "" {
		return nil, nil, fmt.Errorf("ERROR: no keyb file found")
	}

	prog, err := GetPrograms(customKeybPath)
	if err != nil {
		return nil, nil, err
	}
	return prog, config, nil
}

func startProgram(m *model) error {

	p := tea.NewProgram(m)
	p.EnableMouseCellMotion()
	p.EnterAltScreen()
	defer p.ExitAltScreen()

	if err := p.Start(); err != nil {
		return fmt.Errorf("failed to start program: %w", err)
	}
	return nil
}
