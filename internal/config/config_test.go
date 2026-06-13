package config

import (
	"os"
	"testing"
)

// crea un archivo deployctl.yaml temporal para los tests
func setupConfig(t *testing.T, content string) {
	t.Helper()
	err := os.WriteFile("deployctl.yaml", []byte(content), 0644)
	if err != nil {
		t.Fatalf("could not create test config: %v", err)
	}
}

// limpia el archivo después de cada test
func teardown(t *testing.T) {
	t.Helper()
	os.Remove("deployctl.yaml")
}

func TestLoad_ValidConfig(t *testing.T) {
	setupConfig(t, `
environments:
  dev:
    name: mi-nginx-dev
    image: nginx:latest
  prod:
    name: mi-nginx-prod
    image: nginx:1.25
`)
	defer teardown(t)

	cfg, err := Load()

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(cfg.Environments) != 2 {
		t.Fatalf("expected 2 environments, got %d", len(cfg.Environments))
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	teardown(t) // asegurarse que no existe

	_, err := Load()

	if err == nil {
		t.Fatal("expected error when file not found")
	}
}

func TestGetEnvironment_ValidEnv(t *testing.T) {
	setupConfig(t, `
environments:
  dev:
    name: mi-nginx-dev
    image: nginx:latest
`)
	defer teardown(t)

	cfg, _ := Load()
	env, err := cfg.GetEnvironment("dev")

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if env.Name != "mi-nginx-dev" {
		t.Errorf("expected name 'mi-nginx-dev', got '%s'", env.Name)
	}
	if env.Image != "nginx:latest" {
		t.Errorf("expected image 'nginx:latest', got '%s'", env.Image)
	}
}

func TestGetEnvironment_InvalidEnv(t *testing.T) {
	setupConfig(t, `
environments:
  dev:
    name: mi-nginx-dev
    image: nginx:latest
`)
	defer teardown(t)

	cfg, _ := Load()
	_, err := cfg.GetEnvironment("produccion")

	if err == nil {
		t.Fatal("expected error for non-existent environment")
	}
}
