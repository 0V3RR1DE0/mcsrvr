package init

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/0v3rr1de0/mcsrvr/pkg/config"
	"github.com/0v3rr1de0/mcsrvr/pkg/downloader"
	"github.com/0v3rr1de0/mcsrvr/pkg/server/rcon"
)

// CreateStartupScript creates a startup script for the server
func CreateStartupScript(serverPath, jarPath, serverName, memory, javaArgs string) (string, error) {
	var scriptPath string
	var scriptContent string

	// Determine the script extension based on the OS
	if isWindows() {
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
	if err := rcon.EnableRCON(serverPath); err != nil {
		return "", fmt.Errorf("failed to enable RCON: %w", err)
	}

	return scriptPath, nil
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

// isWindows returns true if the current OS is Windows
func isWindows() bool {
	return os.PathSeparator == '\\' && os.PathListSeparator == ';'
}
