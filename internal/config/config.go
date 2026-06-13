package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Environment representa la configuración de un entorno
type Environment struct {
	Name  string `yaml:"name"`
	Image string `yaml:"image"`
}

// Config representa el archivo deployctl.yaml completo
type Config struct {
	Environments map[string]Environment `yaml:"environments"`
}

// Load lee el archivo deployctl.yaml del directorio actual
func Load() (*Config, error) {
	data, err := os.ReadFile("deployctl.yaml")
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("deployctl.yaml not found in current directory")
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	return &cfg, nil
}

// GetEnvironment devuelve la configuración de un entorno específico
func (c *Config) GetEnvironment(env string) (Environment, error) {
	e, ok := c.Environments[env]
	if !ok {
		return Environment{}, fmt.Errorf("environment '%s' not found in deployctl.yaml", env)
	}
	return e, nil
}
