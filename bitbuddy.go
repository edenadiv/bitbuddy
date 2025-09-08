package main

import "time"

const (
	maxStat = 100
	minStat = 0
)

// BitBuddy represents the state of our digital pet.
type BitBuddy struct {
    Name      string
    PetType   string
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
        PetType:   "Cat",
        Hunger:    50,
        Happiness: 50,
        Energy:    50,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
}

// Feed decreases hunger and slightly increases happiness.
func (b *BitBuddy) Feed() {
	b.Hunger -= 20
	if b.Hunger < minStat {
		b.Hunger = minStat
	}
	b.Happiness += 5
	if b.Happiness > maxStat {
		b.Happiness = maxStat
	}
	b.UpdatedAt = time.Now()
}

// Play increases happiness but uses energy.
func (b *BitBuddy) Play() {
	b.Happiness += 20
	if b.Happiness > maxStat {
		b.Happiness = maxStat
	}
	b.Energy -= 15
	if b.Energy < minStat {
		b.Energy = minStat
	}
	b.UpdatedAt = time.Now()
}

// Sleep restores energy.
func (b *BitBuddy) Sleep() {
	b.Energy += 40
	if b.Energy > maxStat {
		b.Energy = maxStat
	}
	b.UpdatedAt = time.Now()
}

// UpdateStats is called on a timer to degrade stats over time.
func (b *BitBuddy) UpdateStats() {
	b.Hunger += 2
	if b.Hunger > maxStat {
		b.Hunger = maxStat
	}
	b.Happiness -= 2
	if b.Happiness < minStat {
		b.Happiness = minStat
	}
	b.UpdatedAt = time.Now()
}
