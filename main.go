package main

import (
    "fmt"
    "math/rand"
    "os"
    "time"

    tea "github.com/charmbracelet/bubbletea"
)

func main() {
    rand.Seed(time.Now().UnixNano())
    buddy, err := load()
    if err != nil {
        fmt.Println("Error loading saved data:", err)
        os.Exit(1)
    }

	p := tea.NewProgram(initialModel(buddy))
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
