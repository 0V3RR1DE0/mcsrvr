package process

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/0v3rr1de0/mcsrvr/pkg/config"
)

// ServerProcess represents a running Minecraft server process
type ServerProcess struct {
	Name    string
	PID     int
	Running bool
}

// ServerProcessInfo represents the serializable information about a running server process
type ServerProcessInfo struct {
	Name    string `json:"name"`
	PID     int    `json:"pid"`
	Running bool   `json:"running"`
	Path    string `json:"path"`
}

// ActiveServers keeps track of running server processes
var ActiveServers = make(map[string]*ServerProcess)

// getActiveServersFilePath returns the path to the file where active servers are stored
func getActiveServersFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".mcsrvr")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return filepath.Join(configDir, "active_servers.json"), nil
}

// SaveActiveServers saves the active servers to a file
func SaveActiveServers() error {
	filePath, err := getActiveServersFilePath()
	if err != nil {
		return err
	}

	// Convert ActiveServers to a serializable format
	activeServersInfo := make(map[string]ServerProcessInfo)
	for name, process := range ActiveServers {
		// Get the server configuration to save the path
		serverConfig, err := config.GetServer(name)
		var path string
		if err == nil {
			path = serverConfig.Path
		}

		activeServersInfo[name] = ServerProcessInfo{
			Name:    process.Name,
			PID:     process.PID,
			Running: process.Running,
			Path:    path,
		}
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(activeServersInfo, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal active servers: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write active servers file: %w", err)
	}

	return nil
}

// LoadActiveServers loads the active servers from a file
func LoadActiveServers() error {
	filePath, err := getActiveServersFilePath()
	if err != nil {
		return err
	}

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// No active servers file, that's okay
		return nil
	}

	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read active servers file: %w", err)
	}

	// Unmarshal from JSON
	var activeServersInfo map[string]ServerProcessInfo
	if err := json.Unmarshal(data, &activeServersInfo); err != nil {
		return fmt.Errorf("failed to unmarshal active servers: %w", err)
	}

	// Convert to ActiveServers format
	for name, info := range activeServersInfo {
		// Check if the process is still running
		running := IsProcessRunning(info.PID)

		if running {
			ActiveServers[name] = &ServerProcess{
				Name:    info.Name,
				PID:     info.PID,
				Running: true,
			}
		}
	}

	return nil
}

// IsProcessRunning checks if a process with the given PID is running
func IsProcessRunning(pid int) bool {
	if pid <= 0 {
		return false
	}

	if runtime.GOOS == "windows" {
		// On Windows, use tasklist to check if the process is running
		// For Minecraft servers, we need to check for java.exe
		cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %d", pid), "/NH")
		output, err := cmd.Output()
		if err == nil && strings.Contains(string(output), fmt.Sprintf("%d", pid)) {
			return true
		}
	} else {
		// On Unix-like systems, use kill -0 to check if the process is running
		cmd := exec.Command("kill", "-0", fmt.Sprintf("%d", pid))
		if err := cmd.Run(); err == nil {
			return true
		}
	}

	return false
}

// FindJavaPID finds the PID of the Java process for a Minecraft server
func FindJavaPID() (int, error) {
	if runtime.GOOS == "windows" {
		// On Windows, use tasklist to find java.exe processes
		cmd := exec.Command("tasklist", "/FI", "IMAGENAME eq java.exe", "/NH", "/FO", "CSV")
		output, err := cmd.Output()
		if err != nil {
			return 0, fmt.Errorf("failed to execute tasklist: %w", err)
		}

		// Parse the output to find the Java PID
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) == "" {
				continue
			}

			// Parse the CSV line
			parts := strings.Split(line, ",")
			if len(parts) < 2 {
				continue
			}

			// Extract the PID
			pidStr := strings.Trim(parts[1], "\"")
			pid, err := strconv.Atoi(pidStr)
			if err != nil {
				continue
			}

			// Return the first Java PID we find
			// In a real implementation, we would need to be more specific
			// to find the correct Java process for the server
			return pid, nil
		}

		return 0, fmt.Errorf("no Java process found")
	} else {
		// On Unix-like systems, use ps to find java processes
		cmd := exec.Command("ps", "-ef", "|", "grep", "java", "|", "grep", "-v", "grep")
		output, err := cmd.Output()
		if err != nil {
			return 0, fmt.Errorf("failed to execute ps: %w", err)
		}

		// Parse the output to find the Java PID
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) == "" {
				continue
			}

			// Parse the line
			parts := strings.Fields(line)
			if len(parts) < 2 {
				continue
			}

			// Extract the PID
			pid, err := strconv.Atoi(parts[1])
			if err != nil {
				continue
			}

			// Return the first Java PID we find
			// In a real implementation, we would need to be more specific
			// to find the correct Java process for the server
			return pid, nil
		}

		return 0, fmt.Errorf("no Java process found")
	}
}

// RefreshServerStatus updates the status of all active servers
func RefreshServerStatus() {
	for name, process := range ActiveServers {
		process.Running = IsProcessRunning(process.PID)
		if !process.Running {
			delete(ActiveServers, name)
		}
	}

	// Save the updated active servers
	if err := SaveActiveServers(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to save active servers: %v\n", err)
	}
}

// init loads the active servers when the package is initialized
func init() {
	if err := LoadActiveServers(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to load active servers: %v\n", err)
	}

	// Refresh the status of all active servers
	RefreshServerStatus()
}
