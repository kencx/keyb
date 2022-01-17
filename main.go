package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const help = `usage: keyb [options] [file]

Flags:
  -p		Print to stdout
  -e		Export to file
  -k		Key bindings at custom path
  -c		Config file at custom path

  -h, --help	help for keyb
`

var (
	output      bool
	export      string
	keybindings string
	keybrc      string
)

func main() {

	flag.BoolVar(&output, "p", false, "print to stdout")
	flag.StringVar(&export, "e", "", "export to file")
	flag.StringVar(&keybindings, "k", "", "keybindings file at custom path")
	flag.StringVar(&keybrc, "c", "", "config file at custom path")

	flag.Usage = func() { os.Stdout.Write([]byte(help)) }
	flag.Parse()

	if keybindings == "" {
		keybindings = "examples/config.yml"
	}

	if keybrc == "" {
		keybrc = "examples/.keybrc"
	}

	// initialize model
	prog, err := GetPrograms(keybindings)
	if err != nil {
		log.Fatal(err)
	}
	config, err := GetConfig(keybrc)
	if err != nil {
		log.Fatal(err)
	}
	m := New(prog, config)

	if output {
		if err := m.OutputBodyToStdout(false); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	if export != "" {
		if err := m.OutputBodyToFile(export, false); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	if err := startProgram(m); err != nil {
		log.Fatal(err)
	}
}

func startProgram(m *model) error {

	p := tea.NewProgram(m)
	p.EnterAltScreen()
	defer p.ExitAltScreen()

	if err := p.Start(); err != nil {
		return fmt.Errorf("failed to start program: %w", err)
	}
	return nil
}

func isFlagPassed(name string) bool {
	var found bool
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
