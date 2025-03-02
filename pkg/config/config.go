package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ServerConfig represents the configuration for a Minecraft server
type ServerConfig struct {
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Version     string    `json:"version"`
	Path        string    `json:"path"`
	Memory      string    `json:"memory"`
	JavaArgs    string    `json:"javaArgs,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	LastStarted time.Time `json:"lastStarted,omitempty"`
}

// Config represents the global configuration for the mcsrvr tool
type Config struct {
	Servers map[string]ServerConfig `json:"servers"`
}

// configDir is the directory where the configuration file is stored
var configDir string

// configFile is the path to the configuration file
var configFile string

// Initialize initializes the configuration
func Initialize() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDir = filepath.Join(homeDir, ".mcsrvr")
	configFile = filepath.Join(configDir, "config.json")

	// Create the config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create the config file if it doesn't exist
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		config := Config{
			Servers: make(map[string]ServerConfig),
		}
		return saveConfig(config)
	}

	return nil
}

// LoadConfig loads the configuration from the config file
func LoadConfig() (Config, error) {
	var config Config

	data, err := os.ReadFile(configFile)
	if err != nil {
		return config, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("failed to parse config file: %w", err)
	}

	return config, nil
}

// saveConfig saves the configuration to the config file
func saveConfig(config Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// AddServer adds a server to the configuration
func AddServer(name, serverType, version, path, memory, javaArgs string) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	// Check if a server with the same name already exists
	if _, exists := config.Servers[name]; exists {
		return fmt.Errorf("server with name '%s' already exists", name)
	}

	// Add the server to the configuration
	config.Servers[name] = ServerConfig{
		Name:      name,
		Type:      serverType,
		Version:   version,
		Path:      path,
		Memory:    memory,
		JavaArgs:  javaArgs,
		CreatedAt: time.Now(),
	}

	return saveConfig(config)
}

// GetServer gets a server from the configuration
func GetServer(name string) (ServerConfig, error) {
	config, err := LoadConfig()
	if err != nil {
		return ServerConfig{}, err
	}

	server, exists := config.Servers[name]
	if !exists {
		return ServerConfig{}, fmt.Errorf("server with name '%s' does not exist", name)
	}

	return server, nil
}

// ListServers lists all servers in the configuration
func ListServers() ([]ServerConfig, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	servers := make([]ServerConfig, 0, len(config.Servers))
	for _, server := range config.Servers {
		servers = append(servers, server)
	}

	return servers, nil
}

// UpdateServer updates a server in the configuration
func UpdateServer(name string, updatedServer ServerConfig) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	if _, exists := config.Servers[name]; !exists {
		return fmt.Errorf("server with name '%s' does not exist", name)
	}

	config.Servers[name] = updatedServer
	return saveConfig(config)
}

// DeleteServer deletes a server from the configuration
func DeleteServer(name string) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	if _, exists := config.Servers[name]; !exists {
		return fmt.Errorf("server with name '%s' does not exist", name)
	}

	delete(config.Servers, name)
	return saveConfig(config)
}

// DefaultConfig represents the default configuration for new servers
type DefaultConfig struct {
	Memory   string `json:"memory"`
	JavaArgs string `json:"javaArgs,omitempty"`
}

// UpdateDefaults updates the default configuration for new servers
func UpdateDefaults(memory, javaArgs string) error {
	// Create the defaults file path
	defaultsFile := filepath.Join(configDir, "defaults.json")

	// Load existing defaults if they exist
	var defaults DefaultConfig
	if _, err := os.Stat(defaultsFile); !os.IsNotExist(err) {
		data, err := os.ReadFile(defaultsFile)
		if err != nil {
			return fmt.Errorf("failed to read defaults file: %w", err)
		}

		if err := json.Unmarshal(data, &defaults); err != nil {
			return fmt.Errorf("failed to parse defaults file: %w", err)
		}
	}

	// Update the defaults
	if memory != "" {
		defaults.Memory = memory
	}
	if javaArgs != "" {
		defaults.JavaArgs = javaArgs
	}

	// Save the defaults
	data, err := json.MarshalIndent(defaults, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal defaults: %w", err)
	}

	if err := os.WriteFile(defaultsFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write defaults file: %w", err)
	}

	return nil
}

// GetDefaults gets the default configuration for new servers
func GetDefaults() (DefaultConfig, error) {
	// Create the defaults file path
	defaultsFile := filepath.Join(configDir, "defaults.json")

	// Load existing defaults if they exist
	var defaults DefaultConfig
	if _, err := os.Stat(defaultsFile); !os.IsNotExist(err) {
		data, err := os.ReadFile(defaultsFile)
		if err != nil {
			return defaults, fmt.Errorf("failed to read defaults file: %w", err)
		}

		if err := json.Unmarshal(data, &defaults); err != nil {
			return defaults, fmt.Errorf("failed to parse defaults file: %w", err)
		}
	} else {
		// Set default values
		defaults.Memory = "2G"
		defaults.JavaArgs = "-XX:+UseG1GC -XX:+ParallelRefProcEnabled"
	}

	return defaults, nil
}
