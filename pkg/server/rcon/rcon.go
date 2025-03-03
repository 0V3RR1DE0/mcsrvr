package rcon

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jltobler/go-rcon"
)

// RCONPort is the default RCON port for Minecraft servers
const RCONPort = 25575

// RCONPassword is the default RCON password for Minecraft servers
const RCONPassword = "mcsrvr"

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

// ConnectRCON connects to the RCON server using the go‑rcon package.
func ConnectRCON(serverName string) (*rcon.Client, error) {
	// Construct a new RCON client.
	// The URL scheme for rcon.NewClient is "rcon://host:port"
	client := rcon.NewClient(fmt.Sprintf("rcon://localhost:%d", RCONPort), RCONPassword)
	return client, nil
}

// ExecuteCommand executes a command on a Minecraft server using RCON.
func ExecuteCommand(serverName, command string) error {
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

// StopServerGracefully stops a Minecraft server gracefully using RCON.
func StopServerGracefully(serverName string) error {
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
