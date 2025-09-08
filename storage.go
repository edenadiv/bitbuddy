package main

import (
	"encoding/json"
	"os"
)

const saveFile = "bitbuddy.json"

// save takes a BitBuddy and saves its state to a JSON file.
func save(buddy *BitBuddy) error {
	data, err := json.MarshalIndent(buddy, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(saveFile, data, 0644)
}

// load reads the state from the JSON file and returns a BitBuddy.
// If the file doesn't exist, it creates a new BitBuddy.
func load() (*BitBuddy, error) {
    data, err := os.ReadFile(saveFile)
    if err != nil {
        if os.IsNotExist(err) {
            // If file doesn't exist, create a new BitBuddy
            return NewBitBuddy("BitBuddy"), nil
        }
        return nil, err
    }

    var buddy BitBuddy
    err = json.Unmarshal(data, &buddy)
    if err != nil {
        return nil, err
    }
    // Backward compatibility: default pet type if missing
    if buddy.PetType == "" {
        buddy.PetType = "Cat"
    }
    return &buddy, nil
}
