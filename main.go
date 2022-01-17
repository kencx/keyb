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
  -k, --key		Key bindings at custom path
  -c, --config	Config file at custom path

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
	flag.BoolVar(&output, "print", false, "print to stdout")
	flag.StringVar(&export, "e", "", "export to file")
	flag.StringVar(&export, "export", "", "export to file")
	flag.StringVar(&keybindings, "k", "", "keybindings file at custom path")
	flag.StringVar(&keybindings, "key", "", "keybindings file at custom path")
	flag.StringVar(&keybrc, "c", "", "config file at custom path")
	flag.StringVar(&keybrc, "config", "", "config file at custom path")

	flag.Usage = func() { os.Stdout.Write([]byte(help)) }
	flag.Parse()

	if err := CreateConfigFile(); err != nil {
		log.Fatal(err)
	}

	// TODO set path
	if keybindings == "" {
		keybindings = "examples/test.yml"
	}

	// TODO set path
	if keybrc == "" {
		keybrc = path.Join(os.Getenv("HOME"), ".config/keyb/.keybrc")
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
	p.EnableMouseCellMotion()
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
