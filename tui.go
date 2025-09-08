package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// -- STYLES ---
var (
	// General
	docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	// Title
	titleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color("#00BFFF")). // DeepSkyBlue
		Padding(0, 1)
	// Status Message
	statusMessageStyle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}). // Mint Green
		Bold(true)
	// UI Panel
	uiPanelStyle = lipgloss.NewStyle().PaddingLeft(2)
	// Menu
	menuChoiceStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("240")) // Gray
	selectedChoiceStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00BFFF"))
	// Quitting
	quitStyle = lipgloss.NewStyle().MarginTop(1).Foreground(lipgloss.Color("240"))
)

// -- ASCII ART --
const (
	artNeutral = `
  /\_/\
 ( o.o )
  > ^ <
`
	artHappy = `
  /\_/\
 ( ^.^ )
  > ^ <
`
	artEating = `
  /\_/\
 ( o.o )
  > w <
`
	artSleeping = `
  /\_/\
 ( -.- )
  > z <
`
)

// -- MESSAGES --
type actionMsg struct{ message string }
type clearStatusMsg struct{}

// -- MODEL --
type model struct {
	buddy         *BitBuddy
	spinner       spinner.Model
	loading       bool
	statusMessage string
	currentAction string // To know which art to display
	cursor        int
	choices       []string
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#00BFFF"))
	return model{
		buddy:   NewBitBuddy("BitBuddy"),
		spinner: s,
		choices: []string{"Feed", "Play", "Sleep"},
	}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

// -- UPDATE --
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.loading {
			return m, nil
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
			m.currentAction = m.choices[m.cursor]
			actionCmd := func() tea.Msg {
				time.Sleep(time.Second * 2)
				switch m.currentAction {
				case "Feed":
					m.buddy.Feed()
					return actionMsg{"Yum, that was tasty!"}
				case "Play":
					m.buddy.Play()
					return actionMsg{"Weee, that was fun!"}
				case "Sleep":
					m.buddy.Sleep()
					return actionMsg{"Zzzz..."}
				}
				return nil
			}
			return m, tea.Sequence(m.spinner.Tick, actionCmd)
		}

	case actionMsg:
		m.loading = false
		m.statusMessage = msg.message
		m.currentAction = ""
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
		if m.loading {
			m.spinner, cmd = m.spinner.Update(msg)
		}
		return m, cmd
	}
	return m, nil
}

// -- VIEW --
func (m model) View() string {
	var art string
	if m.loading {
		switch m.currentAction {
		case "Feed":
			art = artEating
		case "Sleep":
			art = artSleeping
		default:
			art = artHappy
		}
	} else {
		art = artNeutral
	}

	// Left side (ASCII Art)
	artPanel := lipgloss.NewStyle().
		Width(20).
		Height(5).
		Align(lipgloss.Center, lipgloss.Center).
		Render(art)

	// Right side (UI)
	var ui strings.Builder

	// Title
	ui.WriteString(titleStyle.Render("BitBuddy!") + "\n\n")

	// Status or Bars
	if m.loading {
		ui.WriteString(fmt.Sprintf("%s %s...", m.spinner.View(), m.currentAction))
	} else if m.statusMessage != "" {
		ui.WriteString(statusMessageStyle.Render(m.statusMessage))
	} else {
		ui.WriteString(renderBar("Hunger", m.buddy.Hunger) + "\n")
		ui.WriteString(renderBar("Happiness", m.buddy.Happiness) + "\n")
		ui.WriteString(renderBar("Energy", m.buddy.Energy))
	}
	ui.WriteString("\n\n")

	// Menu
	for i, choice := range m.choices {
		style := menuChoiceStyle
		cursor := " "
		if m.cursor == i {
			style = selectedChoiceStyle
			cursor = ">"
		}
		ui.WriteString(style.Render(fmt.Sprintf("%s %s", cursor, choice)) + "\n")
	}

	ui.WriteString(quitStyle.Render("Press 'q' to quit."))
	uiPanel := uiPanelStyle.Render(ui.String())

	return docStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top, artPanel, uiPanel))
}

func renderBar(label string, value int) string {
	if value < 0 {
		value = 0
	}
	if value > 100 {
		value = 100
	}

	barColor := lipgloss.Color("#04B575") // Green
	if value < 50 {
		barColor = lipgloss.Color("#FFB347") // Orange
	}
	if value < 25 {
		barColor = lipgloss.Color("#FF6961") // Red
	}

	barStyle := lipgloss.NewStyle().Foreground(barColor)
	labelStyle := lipgloss.NewStyle().Width(10)

	bar := strings.Repeat("█", value/10) + strings.Repeat("░", 10-value/10)

	return fmt.Sprintf("%s %s", labelStyle.Render(label), barStyle.Render(bar))
}
