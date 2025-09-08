package main

import (
	"fmt"
	"strings"
	tea "github.com/charmbracelet/bubbletea"
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
	b.WriteString("BitBuddy! (Simplified UI)\n\n")
	b.WriteString(renderBar("Hunger", m.buddy.Hunger) + "\n")
	b.WriteString(renderBar("Happiness", m.buddy.Happiness) + "\n")
	b.WriteString(renderBar("Energy", m.buddy.Energy) + "\n\n")

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		b.WriteString(fmt.Sprintf("%s %s\n", cursor, choice))
	}

	b.WriteString("\nPress 'q' to quit.\n")
	return b.String()
}

// A simplified renderBar that doesn't use lipgloss
func renderBar(label string, value int) string {
	if value < 0 {
		value = 0
	}
	if value > 100 {
		value = 100
	}
	bar := ""
	for i := 0; i < 10; i++ {
		if i < value/10 {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	return fmt.Sprintf("% -10s %s", label, bar)
}