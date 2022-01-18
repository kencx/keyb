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
  -p, --print	Print to stdout
  -e, --export	Export to file
  -s, --strip   Strip ANSI chars (only for print/export)
  -k, --key     Key bindings at custom path
  -c, --config	Config file at custom path

  -h, --help	help for keyb
`

func main() {

	var (
		strip      bool
		output     bool
		exportFile string
		keybFile   string
		configFile string
	)
	var defaultConfig = path.Join(os.Getenv("HOME"), ".config/keyb/config")
	var defaultKeyb = path.Join(os.Getenv("HOME"), ".config/keyb/keyb.yaml")

	flag.BoolVar(&strip, "s", false, "strip ANSI chars")
	flag.BoolVar(&strip, "strip", false, "strip ANSI chars")
	flag.BoolVar(&output, "p", false, "print to stdout")
	flag.BoolVar(&output, "print", false, "print to stdout")
	flag.StringVar(&exportFile, "e", "", "export to file")
	flag.StringVar(&exportFile, "export", "", "export to file")
	flag.StringVar(&keybFile, "k", defaultKeyb, "keybindings file")
	flag.StringVar(&keybFile, "key", defaultKeyb, "keybindings file")
	flag.StringVar(&configFile, "c", defaultConfig, "config file")
	flag.StringVar(&configFile, "config", defaultConfig, "config file")

	flag.Usage = func() { os.Stdout.Write([]byte(help)) }
	flag.Parse()

	prog, config, err := handleFlags(keybFile, configFile)
	if err != nil {
		log.Fatal(err)
	}
	m := New(prog, config)

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

func handleFlags(keybFile, configFile string) (map[string]Program, *Config, error) {

	if err := createConfigFile(); err != nil {
		return nil, nil, err
	}
	config, err := GetConfig(configFile)
	if err != nil {
		return nil, nil, err
	}

	customKeybFile := config.KeybPath
	if customKeybFile != "" {
		keybFile = customKeybFile
	}
	prog, err := GetPrograms(keybFile)
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
