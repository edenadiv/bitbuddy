package main

import (
	"fmt"
	"math/rand"
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

// -- ASCII ART FRAMES (ASCII-only for stability) --
const (
    // Cat frames
    catIdle1 = "" +
        "  /\\_/\\    \n" +
        " ( o.o )    \n" +
        "  > ^ <     \n"
    catIdle2 = "" +
        "  /\\_/\\    \n" +
        " ( o_o )    \n" +
        "  > ^ <     \n"
    catIdle3 = "" +
        "  /\\_/\\    \n" +
        " ( -_- )    \n" +
        "  > ^ <     \n"

    catEat1 = "" +
        "  /\\_/\\    \n" +
        " ( o.o )    \n" +
        "  > w <     \n"
    catEat2 = "" +
        "  /\\_/\\    \n" +
        " ( o.o )    \n" +
        "  > o <     \n"

    catPlay1 = "" +
        "  /\\_/\\    \n" +
        " ( ^.^ )    \n" +
        "  > ^ <     \n"
    catPlay2 = "" +
        "  /\\_/\\    \n" +
        " ( ^o^ )    \n" +
        "  > ^ <     \n"

    catSleep1 = "" +
        "  /\\_/\\    \n" +
        " ( -.- ) z  \n" +
        "  > ^ <     \n"
    catSleep2 = "" +
        "  /\\_/\\    \n" +
        " ( -.- ) zz \n" +
        "  > ^ <     \n"

    // Corgi frames (distinct muzzle and tiny legs)
    dogIdle1 = "" +
        "  /\\_/\\    \n" +
        " ( o.o )>   \n" +
        "  |_ _|     \n"
    dogIdle2 = "" +
        "  /\\_/\\    \n" +
        " ( o_o )>   \n" +
        "  |_ _|     \n"
    dogIdle3 = "" +
        "  /\\_/\\    \n" +
        " ( -_- )>   \n" +
        "  |_ _|     \n"

    dogEat1 = "" +
        "  /\\_/\\    \n" +
        " ( o.o )>   \n" +
        "  |\\_/|     \n"
    dogEat2 = "" +
        "  /\\_/\\    \n" +
        " ( o.o )>   \n" +
        "  | o |     \n"

    dogPlay1 = "" +
        "  /\\_/\\    \n" +
        " ( ^.^ )>   \n" +
        "  |_ _|     \n"
    dogPlay2 = "" +
        "  /\\_/\\    \n" +
        " ( ^o^ )>   \n" +
        "  |_ _|     \n"

    dogSleep1 = "" +
        "  /\\_/\\    \n" +
        " ( -.- )> z \n" +
        "  |_ _|     \n"
    dogSleep2 = "" +
        "  /\\_/\\    \n" +
        " ( -.- )>zz \n" +
        "  |_ _|     \n"

    // Bunny frames (tall ears)
    bunIdle1 = "" +
        "  (\\_/ )    \n" +
        "  ( o.o)     \n" +
        "  / > <\\    \n"
    bunIdle2 = "" +
        "  (\\_/ )    \n" +
        "  ( o_o)     \n" +
        "  / > <\\    \n"
    bunIdle3 = "" +
        "  (\\_/ )    \n" +
        "  ( -_-)     \n" +
        "  / > <\\    \n"

    bunEat1 = "" +
        "  (\\_/ )    \n" +
        "  ( o.o)     \n" +
        "  / w w\\    \n"
    bunEat2 = "" +
        "  (\\_/ )    \n" +
        "  ( o.o)     \n" +
        "  / o o\\    \n"

    bunPlay1 = "" +
        "  (\\_/ )    \n" +
        "  ( ^.^)     \n" +
        "  / > <\\    \n"
    bunPlay2 = "" +
        "  (\\_/ )    \n" +
        "  ( ^o^)     \n" +
        "  / > <\\    \n"

    bunSleep1 = "" +
        "  (\\_/ )    \n" +
        "  ( -.-) z   \n" +
        "  / > <\\    \n"
    bunSleep2 = "" +
        "  (\\_/ )    \n" +
        "  ( -.-) zz  \n" +
        "  / > <\\    \n"
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

    // Action overlays
    confetti []confettiParticle
    zzzs     []zzzParticle

    // UI toggles
    showHelp bool
    dark     bool

    // Rename flow
    renaming  bool
    nameInput string
}

type star struct {
    x, y int
    on   bool
}

type confettiParticle struct {
    x, y int
    dx   int
    dy   int
    life int
    ch   string // ASCII glyph
}

type zzzParticle struct {
    x, y int
    life int
    text string // "z", "zz", "zzz"
}

func initialModel(buddy *BitBuddy) model {
    s := spinner.New()
    s.Spinner = spinner.Points
    s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#00BFFF"))
    m := model{
        buddy:   buddy,
        spinner: s,
        choices: []string{"Feed", "Play", "Sleep", "Rename"},
        dark:    true,
    }
    setTheme(m.dark)
    return m
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
        // Handle rename input mode first
        if m.renaming {
            switch msg.Type {
            case tea.KeyEnter:
                trimmed := strings.TrimSpace(m.nameInput)
                if trimmed != "" {
                    m.buddy.Name = trimmed
                    _ = save(m.buddy)
                    m.statusMessage = "Renamed to: " + trimmed
                }
                m.renaming = false
                m.nameInput = ""
                return m, nil
            case tea.KeyEsc:
                m.renaming = false
                m.nameInput = ""
                return m, nil
            case tea.KeyBackspace, tea.KeyCtrlH:
                if len(m.nameInput) > 0 {
                    m.nameInput = m.nameInput[:len(m.nameInput)-1]
                }
                return m, nil
            default:
                // Append regular characters
                s := msg.String()
                if len(s) == 1 && s >= " " && s <= "~" { // printable ASCII
                    m.nameInput += s
                }
                return m, nil
            }
        }
        if m.loading {
            return m, nil
        }
        switch msg.String() {
        case "ctrl+c", "q":
            _ = save(m.buddy) // Save on quit
            return m, tea.Quit
        case "?", "h":
            m.showHelp = !m.showHelp
            return m, nil
        case "t":
            m.dark = !m.dark
            setTheme(m.dark)
            return m, nil
        case "p":
            // Cycle pets: Cat -> Corgi -> Bunny -> Cat
            switch m.buddy.PetType {
            case "Cat":
                m.buddy.PetType = "Corgi"
            case "Corgi":
                m.buddy.PetType = "Bunny"
            default:
                m.buddy.PetType = "Cat"
            }
            m.statusMessage = "Pet: " + m.buddy.PetType
            return m, nil
        case "up", "k":
            if m.cursor > 0 {
                m.cursor--
            }
        case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
        case "enter":
            m.currentAction = m.choices[m.cursor]
            if m.currentAction == "Rename" {
                m.renaming = true
                m.nameInput = m.buddy.Name
                return m, nil
            }
            m.loading = true
            // Start action-specific overlays
            m.startEffectsForAction()
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
        // Clear overlays when action completes
        m.confetti = nil
        m.zzzs = nil
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
        // Update overlays for current action
        if m.loading {
            switch m.currentAction {
            case "Play":
                m.updateConfetti()
            case "Sleep":
                m.updateZzz()
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
        // Place buddy near the left edge with a small inset
        inset := 1
        if inset+len(line) <= len(canvas[i]) {
            canvas[i] = canvas[i][:inset] + line + canvas[i][inset+len(line):]
        }
    }
    // Overlays: confetti and Zzz over the art
    for _, p := range m.confetti {
        if p.y >= 0 && p.y < len(canvas) {
            row := canvas[p.y]
            if p.x >= 0 && p.x < len(row) {
                left := row[:p.x]
                right := row[p.x+1:]
                canvas[p.y] = left + p.ch + right
            }
        }
    }
    for _, z := range m.zzzs {
        if z.y >= 0 && z.y < len(canvas) {
            row := canvas[z.y]
            if z.x >= 0 && z.x < len(row) {
                left := row[:z.x]
                right := row[z.x+1:]
                canvas[z.y] = left + z.text + right
            }
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
    title := "BitBuddy ✨"
    if m.dark {
        title += " · Dark"
    } else {
        title += " · Light"
    }
    title += " · " + m.buddy.PetType
    ui.WriteString(titleStyle.Render(title) + "\n")

    if m.renaming {
        ui.WriteString("Rename Pet (Enter to save, Esc to cancel)\n\n")
        ui.WriteString("> " + m.nameInput + "\n\n")
        ui.WriteString("Tip: Names are saved to bitbuddy.json")
    } else if m.showHelp {
        // Help overlay
        ui.WriteString("Keys:\n")
        ui.WriteString("  ↑/k, ↓/j  Navigate\n")
        ui.WriteString("  Enter     Do action\n")
        ui.WriteString("  ?         Toggle help\n")
        ui.WriteString("  t         Toggle theme\n")
        ui.WriteString("  p         Switch pet (Cat/Corgi/Bunny)\n")
        ui.WriteString("  q         Quit\n\n")
        ui.WriteString("Legend:\n")
        ui.WriteString("  Hunger/Happiness/Energy bars update over time.\n\n")
        ui.WriteString("Files:\n")
        ui.WriteString("  bitbuddy.json — saved state (ignored by git)\n")
    } else {
        // Status or Bars
        if m.loading {
            ui.WriteString(fmt.Sprintf("%s %s...", m.spinner.View(), m.currentAction))
        } else if m.statusMessage != "" {
            ui.WriteString(statusMessageStyle.Render(m.statusMessage))
        } else {
            // Mood indicator
            mood, face := computeMood(m.buddy)
            ui.WriteString(fmt.Sprintf("Mood: %s %s\n\n", mood, face))
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
        ui.WriteString(quitStyle.Render("Press '?' for help · 't' theme · 'q' quit"))
    }
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

// computeMood derives a simple mood label from stats
func computeMood(b *BitBuddy) (string, string) {
    if b == nil {
        return "--", "--"
    }
    score := b.Happiness + b.Energy - b.Hunger
    switch {
    case score >= 120:
        return "Ecstatic", ":D"
    case score >= 60:
        return "Happy", ":)"
    case score >= 20:
        return "Okay", ":|"
    case score >= -20:
        return "Tired", "-_-"
    default:
        return "Grumpy", ":("
    }
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
    for _, s := range m.stars {
        if s.y >= height || s.x >= width {
            continue
        }
        ch := "."
        if s.on { ch = "*" }
        left := rows[s.y][:s.x]
        right := rows[s.y][s.x+1:]
        rows[s.y] = left + ch + right
    }
    return rows
}

func (m model) renderBuddy() string {
    isDog := m.buddy != nil && strings.EqualFold(m.buddy.PetType, "Corgi")
    isBun := m.buddy != nil && strings.EqualFold(m.buddy.PetType, "Bunny")
    if m.loading {
        switch m.currentAction {
        case "Feed":
            if isDog {
                if m.frame%2 == 0 { return dogEat1 }
                return dogEat2
            } else if isBun {
                if m.frame%2 == 0 { return bunEat1 }
                return bunEat2
            }
            if m.frame%2 == 0 { return catEat1 }
            return catEat2
        case "Play":
            if isDog {
                if m.frame%2 == 0 { return dogPlay1 }
                return dogPlay2
            } else if isBun {
                if m.frame%2 == 0 { return bunPlay1 }
                return bunPlay2
            }
            if m.frame%2 == 0 { return catPlay1 }
            return catPlay2
        case "Sleep":
            if isDog {
                if m.frame%2 == 0 { return dogSleep1 }
                return dogSleep2
            } else if isBun {
                if m.frame%2 == 0 { return bunSleep1 }
                return bunSleep2
            }
            if m.frame%2 == 0 { return catSleep1 }
            return catSleep2
        default:
            // fall through to idle
        }
    }
    switch m.frame % 3 {
    case 0:
        if isDog { return dogIdle1 }
        if isBun { return bunIdle1 }
        return catIdle1
    case 1:
        if isDog { return dogIdle2 }
        if isBun { return bunIdle2 }
        return catIdle2
    default:
        if isDog { return dogIdle3 }
        if isBun { return bunIdle3 }
        return catIdle3
    }
}

// -- THEME --
func setTheme(dark bool) {
    if dark {
        titleStyle = lipgloss.NewStyle().
            Foreground(lipgloss.Color("#0A0A0B")).
            Background(lipgloss.Color("#7DD3FC")).
            Padding(0, 2).Bold(true).MarginBottom(1)
        uiPanelStyle = lipgloss.NewStyle().
            Border(lipgloss.RoundedBorder()).
            BorderForeground(lipgloss.Color("#38BDF8")).
            Padding(1, 2)
        menuChoiceStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
        selectedChoiceStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#38BDF8")).Bold(true)
        quitStyle = lipgloss.NewStyle().MarginTop(1).Foreground(lipgloss.Color("240"))
        docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)
        statusMessageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Bold(true)
    } else {
        titleStyle = lipgloss.NewStyle().
            Foreground(lipgloss.Color("#0A0A0B")).
            Background(lipgloss.Color("#FDE68A")). // warm light
            Padding(0, 2).Bold(true).MarginBottom(1)
        uiPanelStyle = lipgloss.NewStyle().
            Border(lipgloss.RoundedBorder()).
            BorderForeground(lipgloss.Color("#A3A3A3")).
            Padding(1, 2)
        menuChoiceStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
        selectedChoiceStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#2563EB")).Bold(true)
        quitStyle = lipgloss.NewStyle().MarginTop(1).Foreground(lipgloss.Color("238"))
        docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)
        statusMessageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#059669")).Bold(true)
    }
}

// -- OVERLAYS: CONFETTI & ZZZ --
func (m *model) startEffectsForAction() {
    switch m.currentAction {
    case "Play":
        m.initConfetti()
    case "Sleep":
        m.initZzz()
    default:
        m.confetti = nil
        m.zzzs = nil
    }
}

func (m *model) initConfetti() {
    m.confetti = nil
    w := 24
    glyphs := []string{"*", "+", "x"}
    for i := 0; i < 22; i++ {
        g := glyphs[rand.Intn(len(glyphs))]
        p := confettiParticle{
            x:   rand.Intn(w),
            y:   rand.Intn(2), // top rows
            dx:  rand.Intn(3) - 1,
            dy:  1,
            life: 8 + rand.Intn(7),
            ch:  g,
        }
        m.confetti = append(m.confetti, p)
    }
}

func (m *model) updateConfetti() {
    w, h := 24, 7
    next := m.confetti[:0]
    for i := range m.confetti {
        p := m.confetti[i]
        p.x += p.dx
        p.y += p.dy
        p.life--
        if p.x < 0 || p.x >= w || p.y < 0 || p.y >= h || p.life <= 0 {
            continue
        }
        // small horizontal jitter
        if rand.Intn(3) == 0 {
            p.x += (rand.Intn(3) - 1)
            if p.x < 0 {
                p.x = 0
            }
            if p.x >= w {
                p.x = w - 1
            }
        }
        next = append(next, p)
    }
    m.confetti = next
    // Occasionally spawn new pieces while playing
    if len(m.confetti) < 18 {
        g := []string{"*", "+", "x"}
        m.confetti = append(m.confetti, confettiParticle{
            x:   rand.Intn(24),
            y:   0,
            dx:  rand.Intn(3) - 1,
            dy:  1,
            life: 10,
            ch:  g[rand.Intn(len(g))],
        })
    }
}

func (m *model) initZzz() {
    m.zzzs = nil
    // start near the head region of the buddy art
    startX := 12
    startY := 1
    m.zzzs = append(m.zzzs, zzzParticle{x: startX, y: startY, life: 14, text: "z"})
}

func (m *model) updateZzz() {
    w, _ := 24, 7
    next := m.zzzs[:0]
    for i := range m.zzzs {
        z := m.zzzs[i]
        // drift up-right slowly
        if rand.Intn(2) == 0 {
            z.x += 1
        }
        if rand.Intn(2) == 0 {
            z.y -= 1
        }
        if z.x >= w {
            z.x = w - 1
        }
        z.life--
        if z.y < 0 || z.life <= 0 {
            continue
        }
        next = append(next, z)
    }
    m.zzzs = next
    // spawn new z every few frames, up to a small count
    if len(m.zzzs) < 4 && m.frame%3 == 0 {
        startX := 12 + rand.Intn(3) - 1
        startY := 2
        texts := []string{"z", "zz", "zzz"}
        m.zzzs = append(m.zzzs, zzzParticle{x: startX, y: startY, life: 14, text: texts[rand.Intn(len(texts))]})
    }
}
