package main

import "time"

// BitBuddy represents the state of our digital pet.
type BitBuddy struct {
	Name      string
	Hunger    int // Goes up over time, decreases when fed
	Happiness int // Goes down over time, increases when played with
	Energy    int // Goes down when playing, increases when sleeping
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewBitBuddy creates a new BitBuddy with default stats.
func NewBitBuddy(name string) *BitBuddy {
	return &BitBuddy{
		Name:      name,
		Hunger:    50,
		Happiness: 50,
		Energy:    50,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
