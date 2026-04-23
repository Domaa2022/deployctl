package history

import (
	"os"
	"testing"
)

// antes de cada test limpiamos el historial
func setupTest(t *testing.T) {
	t.Helper()
	path := historyPath()
	os.Remove(path)
}

func TestLoad_EmptyWhenNoFile(t *testing.T) {
	setupTest(t)

	entries, err := Load()

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected empty history, got %d entries", len(entries))
	}
}

func TestAdd_SavesEntry(t *testing.T) {
	setupTest(t)

	err := Add("mi-nginx", "nginx:latest")

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	entries, _ := Load()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Name != "mi-nginx" {
		t.Errorf("expected name 'mi-nginx', got '%s'", entries[0].Name)
	}
	if entries[0].Image != "nginx:latest" {
		t.Errorf("expected image 'nginx:latest', got '%s'", entries[0].Image)
	}
}

func TestPrevious_ReturnsPreviousImage(t *testing.T) {
	setupTest(t)

	Add("mi-nginx", "nginx:latest")
	Add("mi-nginx", "nginx:1.25")

	previous, err := Previous("mi-nginx")

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if previous != "nginx:latest" {
		t.Errorf("expected 'nginx:latest', got '%s'", previous)
	}
}

func TestPrevious_ErrorWhenNotEnoughHistory(t *testing.T) {
	setupTest(t)

	Add("mi-nginx", "nginx:latest")

	_, err := Previous("mi-nginx")

	if err == nil {
		t.Fatal("expected error when only one deploy in history")
	}
}
