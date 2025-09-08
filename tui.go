package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// actionMsg is sent when an action is completed.
type actionMsg struct {
	message string
}

// clearStatusMsg is sent to clear the status message after a delay.
type clearStatusMsg struct{}


type model struct {
	buddy         *BitBuddy
	spinner       spinner.Model
	loading       bool
	statusMessage string
	cursor        int
	choices       []string
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return model{
		buddy:   NewBitBuddy("BitBuddy"),
		spinner: s,
		choices: []string{"Feed", "Play", "Sleep"},
	}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if m.loading {
			return m, nil // Don't allow input while loading
		}
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
		case "enter":
			m.loading = true
			choice := m.choices[m.cursor]

			// Define the action command with a delay
			actionCmd := func() tea.Msg {
				time.Sleep(time.Second * 1)
				switch choice {
				case "Feed":
					m.buddy.Feed()
					return actionMsg{"Yum, that was tasty! (^o^)"}
				case "Play":
					m.buddy.Play()
					return actionMsg{"Weee, that was fun! (>^ω^<)"}
				case "Sleep":
					m.buddy.Sleep()
					return actionMsg{"Zzzz... (u_u)"}
				}
				return nil
			}

			return m, tea.Sequence(m.spinner.Tick, actionCmd)
		}

	case actionMsg:
		m.loading = false
		m.statusMessage = msg.message
		// Define the command to clear the status message after a delay
		clearMsgCmd := func() tea.Msg {
			time.Sleep(time.Second * 2)
			return clearStatusMsg{}
		}
		return m, clearMsgCmd

	case clearStatusMsg:
		m.statusMessage = ""
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		if m.loading { // Only tick the spinner when loading
			m.spinner, cmd = m.spinner.Update(msg)
		}
		return m, cmd

	default:
		return m, nil
	}
	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	if m.loading {
		fmt.Fprintf(&b, "%s Performing action...", m.spinner.View())
		return b.String()
	}

	if m.statusMessage != "" {
		fmt.Fprintf(&b, "BitBuddy says: %s\n\n", m.statusMessage)
	} else {
		fmt.Fprintf(&b, "Your BitBuddy! (^_^)\n\n")
	}

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
	// Clamp value between 0 and 100
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
	return bar
}