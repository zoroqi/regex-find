package app

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const MaxHistorySize = 500

// HistoryItem represents a single entry in the history.
type HistoryItem struct {
	Regex      string `json:"regex"`
	FirstMatch string `json:"firstMatch"`
	Timestamp  int64  `json:"ts"`
	Count      int    `json:"count"`
}

// History represents the collection of history items.
type History struct {
	Patterns []HistoryItem `json:"patterns"`
}

// LoadHistory loads the history from the given path.
// If the path is empty, it returns an empty History struct.
func LoadHistory(path string) (History, error) {
	if path == "" {
		return History{Patterns: []HistoryItem{}}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return History{Patterns: []HistoryItem{}}, nil
		}
		return History{}, err
	}

	var history History
	if err := json.Unmarshal(data, &history); err != nil {
		return History{}, err
	}

	return history, nil
}

// SaveHistory saves the history to the given path.
// If the path is empty, it does nothing.
func SaveHistory(path string, data History) error {
	if path == "" {
		return nil
	}

	// Ensure the directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, jsonData, 0644)
}
