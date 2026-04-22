package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Entry struct {
	Name       string    `json:"name"`
	Image      string    `json:"image"`
	DeployedAt time.Time `json:"deployed_at"`
}

func historyPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".deployctl", "history.json")
}

func Load() ([]Entry, error) {
	path := historyPath()

	// si el archivo no existe, devuelve historial vacío
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []Entry{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading history: %w", err)
	}

	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("parsing history: %w", err)
	}

	return entries, nil
}

func Save(entries []Entry) error {
	path := historyPath()

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding history: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}

func Add(name, image string) error {
	entries, err := Load()
	if err != nil {
		return err
	}

	entries = append(entries, Entry{
		Name:       name,
		Image:      image,
		DeployedAt: time.Now(),
	})

	return Save(entries)
}

func Previous(name string) (string, error) {
	entries, err := Load()
	if err != nil {
		return "", err
	}

	// buscar las últimas dos entradas para este nombre
	var matches []Entry
	for _, e := range entries {
		if e.Name == name {
			matches = append(matches, e)
		}
	}

	if len(matches) < 2 {
		return "", fmt.Errorf("no previous deployment found for '%s'", name)
	}

	// devolver la penúltima
	return matches[len(matches)-2].Image, nil
}
