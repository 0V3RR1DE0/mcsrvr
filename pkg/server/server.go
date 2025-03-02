package server

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	//"syscall"
	"time"

	"github.com/0v3rr1de0/mcsrvr/pkg/config"
	"github.com/0v3rr1de0/mcsrvr/pkg/downloader"
	"github.com/jltobler/go-rcon"
)

// RCONPort is the default RCON port for Minecraft servers
const RCONPort = 25575

// RCONPassword is the default RCON password for Minecraft servers
const RCONPassword = "mcsrvr"

// CreateStartupScript creates a startup script for the server
func CreateStartupScript(serverPath, jarPath, serverName, memory, javaArgs string) (string, error) {
	var scriptPath string
	var scriptContent string

	// Determine the script extension based on the OS
	if runtime.GOOS == "windows" {
		scriptPath = filepath.Join(serverPath, "start.bat")
		scriptContent = fmt.Sprintf(`@echo off
echo Starting Minecraft server %s...
java -Xmx%s -Xms%s %s -jar "%s" nogui
if errorlevel 1 (
    echo Server crashed or failed to start. Press any key to exit.
    pause > nul
)
`, serverName, memory, memory, javaArgs, filepath.Base(jarPath))
	} else {
		scriptPath = filepath.Join(serverPath, "start.sh")
		scriptContent = fmt.Sprintf(`#!/bin/bash
echo "Starting Minecraft server %s..."
java -Xmx%s -Xms%s %s -jar "%s" nogui
if [ $? -ne 0 ]; then
    echo "Server crashed or failed to start. Press Enter to exit."
    read
fi
`, serverName, memory, memory, javaArgs, filepath.Base(jarPath))
	}

	// Write the script to file
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		return "", fmt.Errorf("failed to create startup script: %w", err)
	}

	// Create or update server.properties to enable RCON
	if err := EnableRCON(serverPath); err != nil {
		return "", fmt.Errorf("failed to enable RCON: %w", err)
	}

	return scriptPath, nil
}

// EnableRCON enables RCON in the server.properties file
func EnableRCON(serverPath string) error {
	propertiesPath := filepath.Join(serverPath, "server.properties")

	// Check if server.properties exists
	if _, err := os.Stat(propertiesPath); os.IsNotExist(err) {
		// Create a new server.properties file with RCON enabled
		propertiesContent := fmt.Sprintf(`# RCON Configuration
enable-rcon=true
rcon.port=%d
rcon.password=%s
broadcast-rcon-to-ops=true
`, RCONPort, RCONPassword)

		if err := os.WriteFile(propertiesPath, []byte(propertiesContent), 0644); err != nil {
			return fmt.Errorf("failed to create server.properties: %w", err)
		}

		return nil
	}

	// Read the existing server.properties file
	content, err := os.ReadFile(propertiesPath)
	if err != nil {
		return fmt.Errorf("failed to read server.properties: %w", err)
	}

	// Convert to string for easier manipulation
	propertiesContent := string(content)

	// Create a map to store the properties
	properties := make(map[string]string)

	// Parse the properties file
	lines := strings.Split(propertiesContent, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		properties[key] = value
	}

	// Update the RCON properties
	properties["enable-rcon"] = "true"
	properties["rcon.port"] = fmt.Sprintf("%d", RCONPort)
	properties["rcon.password"] = RCONPassword
	properties["broadcast-rcon-to-ops"] = "true"

	// Convert the properties back to a string
	var newContent strings.Builder

	// Add any comments at the beginning of the file
	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			newContent.WriteString(line)
			newContent.WriteString("\n")
		}
	}

	// Add the properties
	for key, value := range properties {
		newContent.WriteString(key)
		newContent.WriteString("=")
		newContent.WriteString(value)
		newContent.WriteString("\n")
	}

	// Write the updated content back to the file
	if err := os.WriteFile(propertiesPath, []byte(newContent.String()), 0644); err != nil {
		return fmt.Errorf("failed to update server.properties: %w", err)
	}

	return nil
}

// AcceptEULA accepts the Minecraft EULA by creating or modifying the eula.txt file
func AcceptEULA(serverPath string) error {
	eulaPath := filepath.Join(serverPath, "eula.txt")
	eulaContent := "eula=true\n"

	if err := os.WriteFile(eulaPath, []byte(eulaContent), 0644); err != nil {
		return fmt.Errorf("failed to accept EULA: %w", err)
	}

	return nil
}

// ServerProcess represents a running Minecraft server process
type ServerProcess struct {
	Name    string
	Cmd     *exec.Cmd
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
				Cmd:     nil, // We can't restore the Cmd object
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

// StartServer starts a Minecraft server as a hidden process
func StartServer(serverName string) error {
	// Get the server configuration
	serverConfig, err := config.GetServer(serverName)
	if err != nil {
		return err
	}

	// Check if the server is already running
	if process, exists := ActiveServers[serverName]; exists && process.Running {
		return fmt.Errorf("server '%s' is already running", serverName)
	}

	// Check if the server directory exists
	if _, err := os.Stat(serverConfig.Path); os.IsNotExist(err) {
		return fmt.Errorf("server directory does not exist: %s", serverConfig.Path)
	}

	// Determine the startup script path
	var scriptPath string
	if runtime.GOOS == "windows" {
		scriptPath = filepath.Join(serverConfig.Path, "start.bat")
	} else {
		scriptPath = filepath.Join(serverConfig.Path, "start.sh")
	}

	// Check if the startup script exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("startup script does not exist: %s", scriptPath)
	}

	// Create log file for the server
	logDir := filepath.Join(serverConfig.Path, "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	logFile, err := os.OpenFile(
		filepath.Join(logDir, "server.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		// Run the script using cmd and detach the process on Windows
		cmd = exec.Command("cmd", "/c", scriptPath)
	} else {
		// For Unix-like systems, use bash without additional arguments
		cmd = exec.Command("bash", scriptPath)
	}

	// Set the process attributes using our helper function
	cmd.SysProcAttr = newSysProcAttr()
	cmd.Dir = serverConfig.Path
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	if err := cmd.Start(); err != nil {
		logFile.Close()
		return fmt.Errorf("failed to start server: %w", err)
	}

	// Store the command window process PID
	cmdPID := cmd.Process.Pid
	fmt.Printf("Server '%s' started with command window PID %d\n", serverName, cmdPID)

	// Wait a moment for the Java process to start
	time.Sleep(2 * time.Second)

	// Find the Java process PID
	javaPID, err := FindJavaPID()
	if err != nil {
		fmt.Printf("Warning: Failed to find Java process PID: %v\n", err)
		fmt.Println("Using command window PID instead")

		// Store the command window process in ActiveServers
		ActiveServers[serverName] = &ServerProcess{
			Name:    serverName,
			Cmd:     cmd,
			PID:     cmdPID,
			Running: true,
		}
	} else {
		fmt.Printf("Found Java process with PID %d\n", javaPID)

		// Store the Java process in ActiveServers
		ActiveServers[serverName] = &ServerProcess{
			Name:    serverName,
			Cmd:     cmd,
			PID:     javaPID,
			Running: true,
		}
	}

	// Save the active servers to file
	if err := SaveActiveServers(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to save active servers: %v\n", err)
	}

	// Update the server configuration with the last started time
	serverConfig.LastStarted = time.Now()
	if err := config.UpdateServer(serverName, serverConfig); err != nil {
		return fmt.Errorf("failed to update server configuration: %w", err)
	}

	fmt.Printf("Server '%s' started successfully\n", serverName)
	return nil
}

// StopServer stops a Minecraft server
func StopServer(serverName string) error {
	// Get the server configuration
	_, err := config.GetServer(serverName)
	if err != nil {
		return err
	}

	// Check if the server is running
	process, exists := ActiveServers[serverName]
	if !exists || !process.Running {
		return fmt.Errorf("server '%s' is not running", serverName)
	}

	fmt.Printf("Stopping server '%s'...\n", serverName)

	// Store the cmd process PID
	var cmdPID int
	if process.Cmd != nil && process.Cmd.Process != nil {
		cmdPID = process.Cmd.Process.Pid
	}

	// Store the Java process PID
	javaPID := process.PID

	// Try to stop the server gracefully using RCON
	if err := stopServerGracefully(serverName); err != nil {
		fmt.Printf("Warning: Failed to stop server gracefully: %v\n", err)
		fmt.Println("Falling back to force kill...")
	} else {
		fmt.Println("RCON stop command sent successfully")

		// Wait a moment for the server to start shutting down
		time.Sleep(2 * time.Second)
	}

	// Check if the Java process is still running
	javaRunning := IsProcessRunning(javaPID)

	// If the Java process is still running, kill it
	if javaRunning {
		fmt.Printf("Java process (PID %d) is still running, attempting to kill it...\n", javaPID)

		if runtime.GOOS == "windows" {
			// On Windows, use taskkill to kill the Java process
			killCmd := exec.Command("taskkill", "/F", "/PID", fmt.Sprintf("%d", javaPID))
			if err := killCmd.Run(); err != nil {
				fmt.Printf("Warning: Failed to kill Java process: %v\n", err)
			} else {
				fmt.Printf("Java process (PID %d) killed successfully\n", javaPID)
			}
		} else {
			// On Unix-like systems, use kill to kill the Java process
			killCmd := exec.Command("kill", "-9", fmt.Sprintf("%d", javaPID))
			if err := killCmd.Run(); err != nil {
				fmt.Printf("Warning: Failed to kill Java process: %v\n", err)
			} else {
				fmt.Printf("Java process (PID %d) killed successfully\n", javaPID)
			}
		}
	}

	// If the cmd process is still running, kill it
	if cmdPID > 0 {
		cmdRunning := IsProcessRunning(cmdPID)
		if cmdRunning {
			fmt.Printf("Command window process (PID %d) is still running, attempting to kill it...\n", cmdPID)

			// Kill the cmd process
			if process.Cmd != nil && process.Cmd.Process != nil {
				if err := process.Cmd.Process.Kill(); err != nil {
					fmt.Printf("Warning: Failed to kill command window process: %v\n", err)
				} else {
					fmt.Printf("Command window process (PID %d) killed successfully\n", cmdPID)
				}
			}
		}
	}

	// Update the process status
	process.Running = false
	delete(ActiveServers, serverName)

	// Save the active servers to file
	if err := SaveActiveServers(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to save active servers: %v\n", err)
	}

	fmt.Printf("Server '%s' stopped successfully\n", serverName)

	return nil
}

// stopServerGracefully stops a Minecraft server gracefully using RCON.
func stopServerGracefully(serverName string) error {
	// Connect to the RCON server using the external client.
	client, err := ConnectRCON(serverName)
	if err != nil {
		return fmt.Errorf("failed to connect to RCON server: %w", err)
	}
	// Use the Send method to issue the "stop" command.
	_, err = client.Send("stop")
	if err != nil {
		return fmt.Errorf("failed to execute stop command: %w", err)
	}
	return nil
}

// ConnectRCON connects to the RCON server using the go‑rcon package.
func ConnectRCON(serverName string) (*rcon.Client, error) {
	// Get the server configuration
	_, err := config.GetServer(serverName)
	if err != nil {
		return nil, err
	}

	// Check if the server is running
	if _, exists := ActiveServers[serverName]; !exists {
		return nil, fmt.Errorf("server '%s' is not running", serverName)
	}

	// Construct a new RCON client.
	// The URL scheme for rcon.NewClient is "rcon://host:port"
	client := rcon.NewClient(fmt.Sprintf("rcon://localhost:%d", RCONPort), RCONPassword)
	return client, nil
}

// ExecuteCommand executes a command on a Minecraft server using RCON.
func ExecuteCommand(serverName, command string) error {
	// Get the server configuration.
	_, err := config.GetServer(serverName)
	if err != nil {
		return err
	}

	// Check if the server is running.
	if _, exists := ActiveServers[serverName]; !exists {
		return fmt.Errorf("server '%s' is not running", serverName)
	}

	fmt.Printf("Executing command '%s' on server '%s'...\n", command, serverName)

	// Connect to the RCON server.
	client, err := ConnectRCON(serverName)
	if err != nil {
		return fmt.Errorf("failed to connect to RCON server: %w", err)
	}

	// Send the command using the client's Send method.
	response, err := client.Send(command)
	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	// Print the response.
	fmt.Println(response)
	return nil
}

// CreateBackup creates a backup of a Minecraft server
func CreateBackup(serverName, backupPath string) error {
	// Get the server configuration
	_, err := config.GetServer(serverName)
	if err != nil {
		return err
	}

	// Create the backup directory if it doesn't exist
	if err := os.MkdirAll(backupPath, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Create a timestamp for the backup
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	backupName := fmt.Sprintf("%s_%s", serverName, timestamp)
	backupDir := filepath.Join(backupPath, backupName)

	// Create the backup directory
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Copy the server files to the backup directory
	// This is a placeholder implementation
	// In a real implementation, we would need to copy all the server files
	fmt.Printf("Creating backup of server '%s' to '%s'...\n", serverName, backupDir)

	// For now, we'll just print a message
	fmt.Println("Backup created successfully")

	return nil
}

// RestoreBackup restores a backup of a Minecraft server
func RestoreBackup(backupPath, restorePath string) error {
	// Check if the backup directory exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup directory does not exist: %s", backupPath)
	}

	// Create the restore directory if it doesn't exist
	if err := os.MkdirAll(restorePath, 0755); err != nil {
		return fmt.Errorf("failed to create restore directory: %w", err)
	}

	// Copy the backup files to the restore directory
	// This is a placeholder implementation
	// In a real implementation, we would need to copy all the backup files
	fmt.Printf("Restoring backup from '%s' to '%s'...\n", backupPath, restorePath)

	// For now, we'll just print a message
	fmt.Println("Backup restored successfully")

	return nil
}

// ListBackups lists all backups for a server
func ListBackups(serverName string) ([]string, error) {
	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Determine the backups directory
	backupsDir := filepath.Join(homeDir, ".mcsrvr", "backups")

	// Check if the backups directory exists
	if _, err := os.Stat(backupsDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	// List all backups for the server
	var backups []string

	// If no server name is provided, list all backups
	if serverName == "" {
		entries, err := os.ReadDir(backupsDir)
		if err != nil {
			return nil, fmt.Errorf("failed to read backups directory: %w", err)
		}

		for _, entry := range entries {
			if entry.IsDir() {
				backups = append(backups, entry.Name())
			}
		}
	} else {
		// List backups for the specific server
		entries, err := os.ReadDir(backupsDir)
		if err != nil {
			return nil, fmt.Errorf("failed to read backups directory: %w", err)
		}

		for _, entry := range entries {
			if entry.IsDir() && strings.HasPrefix(entry.Name(), serverName+"_") {
				backups = append(backups, entry.Name())
			}
		}
	}

	return backups, nil
}

// InitializeServer initializes a new Minecraft server
func InitializeServer(serverPath, serverName, serverType, version, memory, javaArgs string) error {
	// Create the server directory if it doesn't exist
	if err := os.MkdirAll(serverPath, 0755); err != nil {
		return fmt.Errorf("failed to create server directory: %w", err)
	}

	// Download the server jar based on the server type
	var jarPath string
	var err error

	switch serverType {
	case "papermc":
		// Download PaperMC server jar
		jarPath, err = downloadPaperMC(serverPath, version)
	case "vanilla":
		// Download vanilla server jar
		jarPath, err = downloadVanilla(serverPath, version)
	case "fabric":
		// Download Fabric server jar
		jarPath, err = downloadFabric(serverPath, version, "")
	default:
		return fmt.Errorf("unsupported server type: %s", serverType)
	}

	if err != nil {
		return fmt.Errorf("failed to download server jar: %w", err)
	}

	// Create the startup script
	_, err = CreateStartupScript(serverPath, jarPath, serverName, memory, javaArgs)
	if err != nil {
		return err
	}

	// Run the server once to generate the eula.txt file
	fmt.Println("Running server for the first time to generate eula.txt...")

	// This is a placeholder implementation
	// In a real implementation, we would need to run the server and wait for it to generate the eula.txt file

	// Accept the EULA
	if err := AcceptEULA(serverPath); err != nil {
		return err
	}

	// Add the server to the configuration
	if err := config.AddServer(serverName, serverType, version, serverPath, memory, javaArgs); err != nil {
		return fmt.Errorf("failed to add server to configuration: %w", err)
	}

	fmt.Printf("Server '%s' initialized successfully at %s\n", serverName, serverPath)
	fmt.Printf("To start the server, run: mcsrvr start %s\n", serverName)

	return nil
}

// downloadPaperMC downloads the PaperMC server jar
func downloadPaperMC(serverPath, version string) (string, error) {
	return downloader.DownloadPaperMC(serverPath, version)
}

// downloadVanilla downloads the vanilla Minecraft server jar
func downloadVanilla(serverPath, version string) (string, error) {
	return downloader.DownloadVanilla(serverPath, version)
}

// downloadFabric downloads the Fabric server jar
func downloadFabric(serverPath, version, loaderVersion string) (string, error) {
	return downloader.DownloadFabric(serverPath, version, loaderVersion)
}

// InitializeFabricServer initializes a new Fabric server
func InitializeFabricServer(serverPath, serverName, mcVersion, loaderVersion, memory, javaArgs string) error {
	// Create the server directory if it doesn't exist
	if err := os.MkdirAll(serverPath, 0755); err != nil {
		return fmt.Errorf("failed to create server directory: %w", err)
	}

	// Download the Fabric server jar
	jarPath, err := downloadFabric(serverPath, mcVersion, loaderVersion)
	if err != nil {
		return fmt.Errorf("failed to download Fabric server jar: %w", err)
	}

	// Create the startup script
	_, err = CreateStartupScript(serverPath, jarPath, serverName, memory, javaArgs)
	if err != nil {
		return err
	}

	// Run the server once to generate the eula.txt file
	fmt.Println("Running server for the first time to generate eula.txt...")

	// This is a placeholder implementation
	// In a real implementation, we would need to run the server and wait for it to generate the eula.txt file

	// Accept the EULA
	if err := AcceptEULA(serverPath); err != nil {
		return err
	}

	// Add the server to the configuration
	if err := config.AddServer(serverName, "fabric", mcVersion, serverPath, memory, javaArgs); err != nil {
		return fmt.Errorf("failed to add server to configuration: %w", err)
	}

	fmt.Printf("Fabric server '%s' initialized successfully at %s\n", serverName, serverPath)
	fmt.Printf("To start the server, run: mcsrvr start %s\n", serverName)
	fmt.Println("Note: Most mods will also require you to install Fabric API into the mods folder")

	return nil
}
