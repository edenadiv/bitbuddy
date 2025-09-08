package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
)

type model struct {
	buddy   *BitBuddy
	cursor  int
	choices []string
}

func initialModel() model {
	return model{
		buddy:   NewBitBuddy("BitBuddy"),
		choices: []string{"Feed", "Play", "Sleep"},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	fmt.Fprintf(&b, "Your BitBuddy!\n\n")
	fmt.Fprintf(&b, "Hunger:    %s\n", renderBar(m.buddy.Hunger))
	fmt.Fprintf(&b, "Happiness: %s\n", renderBar(m.buddy.Happiness))
	fmt.Fprintf(&b, "Energy:    %s\n\n", renderBar(m.buddy.Energy))

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		fmt.Fprintf(&b, "%s %s\n", cursor, choice)
	}

	fmt.Fprintf(&b, "\nPress 'q' to quit.\n")

	return b.String()
}

func renderBar(value int) string {
	bar := ""
	for i := 0; i < 10; i++ {
		if i < value/10 {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	return bar
}
