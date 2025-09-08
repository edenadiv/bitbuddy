package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	ea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// -- STYLES ---
var (
	// General
	docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	// Title
	titleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#0A0A0B")).
		Background(lipgloss.Color("#7DD3FC")). // Sky
		Padding(0, 2).
		Bold(true).
		MarginBottom(1)
	// Status Message
	statusMessageStyle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}). // Mint Green
		Bold(true)
	// UI Panel
	uiPanelStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#38BDF8")).
		Padding(1, 2)
	// Menu
	menuChoiceStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("240")) // Gray
	selectedChoiceStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#38BDF8")).Bold(true)
	// Quitting
	quitStyle = lipgloss.NewStyle().MarginTop(1).Foreground(lipgloss.Color("240"))
)

// -- ASCII ART FRAMES --
const (
	idle1 = "" +
		"  /\\_/\\\n" +
		" ( •.• )\n" +
		"  > ^ <\n"
	idle2 = "" +
		"  /\\_/\\\n" +
		" ( •_• )\n" +
		"  > ^ <\n"
	idle3 = "" +
		"  /\\_/\\\n" +
		" ( -_- )\n" +
		"  > ^ <\n"

	eat1 = "" +
		"  /\\_/\\\n" +
		" ( •.• )\n" +
		"  > w <\n"
	eat2 = "" +
		"  /\\_/\\\n" +
		" ( •.• )\n" +
		"  > o <\n"

	play1 = "" +
		"  /\\_/\\\n" +
		" ( ^.^ )\n" +
		"  > ^ <\n"
	play2 = "" +
		"  /\\_/\\\n" +
		" ( ^o^ )\n" +
		"  > ^ <\n"

	sleep1 = "" +
		"  /\\_/\\\n" +
		" ( -.- ) z\n" +
		"  > ^ <\n"
	sleep2 = "" +
		"  /\\_/\\\n" +
		" ( -.- ) zz\n" +
		"  > ^ <\n"
)

// -- MESSAGES --
type actionMsg struct{ message string }
type clearStatusMsg struct{}
type tickMsg struct{}
type animTickMsg struct{}

// -- MODEL --
type model struct {
	buddy         *BitBuddy
	spinner       spinner.Model
	loading       bool
	statusMessage string
	currentAction string // To know which art to display
	cursor        int
	choices       []string

	// Animation / layout
	width  int
	height int
	frame  int
	stars  []star
}

type star struct {
	x, y int
	on   bool
}

func initialModel(buddy *BitBuddy) model {
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#00BFFF"))
	return model{
		buddy:   buddy,
		spinner: s,
		choices: []string{"Feed", "Play", "Sleep"},
	}
}

func (m model) Init() tea.Cmd {
	return tea.Sequence(m.spinner.Tick, tick(), animTick())
}

// -- UPDATE --
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width, m.height = msg.Width, msg.Height
        m.initStars()
        return m, nil
    case tea.KeyMsg:
        if m.loading {
            return m, nil
        }
		switch msg.String() {
		case "ctrl+c", "q":
			_ = save(m.buddy) // Save on quit
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

	case tickMsg:
		m.buddy.UpdateStats()
		return m, tea.Sequence(tick(), m.spinner.Tick)

    case spinner.TickMsg:
        var cmd tea.Cmd
        if m.loading {
            m.spinner, cmd = m.spinner.Update(msg)
        }
        return m, cmd
    case animTickMsg:
        m.frame++
        // twinkle some stars randomly
        for i := range m.stars {
            if rand.Intn(12) == 0 {
                m.stars[i].on = !m.stars[i].on
            }
        }
        return m, animTick()
    }
    return m, nil
}

// -- VIEW --
func (m model) View() string {
	// Animated buddy & starfield panel
	art := m.renderBuddy()
	canvas := m.renderStars(24, 7)
	artLines := strings.Split(strings.TrimRight(art, "\n"), "\n")
	for i := range canvas {
		if i >= len(artLines) {
			break
		}
		line := artLines[i]
		pad := 0
		if w := len(canvas[i]); w > len(line) {
			pad = (w - len(line)) / 2
		}
		if pad < 0 {
			pad = 0
		}
		if pad+len(line) <= len(canvas[i]) {
			canvas[i] = canvas[i][:pad] + line + canvas[i][pad+len(line):]
		}
	}
	artPanel := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#334155")).
		Padding(1, 2).
		Render(strings.Join(canvas, "\n"))

	// Right side (UI)
	var ui strings.Builder

	// Title
	ui.WriteString(titleStyle.Render("BitBuddy ✨") + "\n")

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

	barColor := lipgloss.Color("#10B981") // Emerald
	if value < 50 {
		barColor = lipgloss.Color("#F59E0B") // Amber
	}
	if value < 25 {
		barColor = lipgloss.Color("#EF4444") // Red
	}

	barStyle := lipgloss.NewStyle().Foreground(barColor).Bold(true)
	labelStyle := lipgloss.NewStyle().Width(10)

	bar := strings.Repeat("█", value/10) + strings.Repeat("░", 10-value/10)

	return fmt.Sprintf("%s %s", labelStyle.Render(label), barStyle.Render(bar))
}

// tick is a command that sends a tickMsg every 5 seconds.
func tick() tea.Cmd {
    return tea.Tick(time.Second*5, func(t time.Time) tea.Msg {
        return tickMsg{}
    })
}

// animTick is a faster tick for UI animations
func animTick() tea.Cmd {
    return tea.Tick(time.Millisecond*120, func(t time.Time) tea.Msg {
        return animTickMsg{}
    })
}

func (m *model) initStars() {
    width := 24
    height := 7
    if len(m.stars) == 0 {
        for i := 0; i < 30; i++ {
            m.stars = append(m.stars, star{
                x:  rand.Intn(width),
                y:  rand.Intn(height),
                on: rand.Intn(2) == 0,
            })
        }
    }
}

func (m model) renderStars(width, height int) []string {
    rows := make([]string, height)
    for i := 0; i < height; i++ {
        rows[i] = strings.Repeat(" ", width)
    }
    dot := lipgloss.NewStyle().Foreground(lipgloss.Color("#94A3B8")).Render("·")
    bright := lipgloss.NewStyle().Foreground(lipgloss.Color("#E5E7EB")).Bold(true).Render("•")
    for _, s := range m.stars {
        if s.y >= height || s.x >= width {
            continue
        }
        ch := dot
        if s.on {
            ch = bright
        }
        if s.x < len(rows[s.y]) {
            left := rows[s.y][:s.x]
            right := rows[s.y][s.x+1:]
            rows[s.y] = left + ch + right
        }
    }
    return rows
}

func (m model) renderBuddy() string {
    if m.loading {
        switch m.currentAction {
        case "Feed":
            if m.frame%2 == 0 {
                return eat1
            }
            return eat2
        case "Play":
            if m.frame%2 == 0 {
                return play1
            }
            return play2
        case "Sleep":
            if m.frame%2 == 0 {
                return sleep1
            }
            return sleep2
        default:
            switch m.frame % 3 {
            case 0:
                return idle1
            case 1:
                return idle2
            default:
                return idle3
            }
        }
    }
    switch m.frame % 3 {
    case 0:
        return idle1
    case 1:
        return idle2
    default:
        return idle3
    }
}
