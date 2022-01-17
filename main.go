package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	prog := GetConfig()
	m := New(prog)

	p := tea.NewProgram(m)
	p.EnterAltScreen()
	defer p.ExitAltScreen()

	if err := p.Start(); err != nil {
		log.Fatalln(err)
	}
}
