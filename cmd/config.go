package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/0v3rr1de0/mcsrvr/pkg/config"
)

var (
	defaultMemory   string
	defaultJavaArgs string
	rconPort        int
	rconPassword    string
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config [server-name] [config-type]",
	Short: "Configure server settings",
	Long: `Configure server settings such as startup script, server properties, or operator list.
Config types: start (startup script), properties (server.properties), ops (ops.json), rcon (RCON settings)

Example:
  mcsrvr config paper123 start
  mcsrvr config paper123 properties
  mcsrvr config paper123 ops
  mcsrvr config paper123 rcon --port 25575 --password mypassword
  mcsrvr config --default-memory 4G
  mcsrvr config --default-java-args "-XX:+UseG1GC -XX:+ParallelRefProcEnabled"`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if we're setting default values
		if cmd.Flags().Changed("default-memory") || cmd.Flags().Changed("default-java-args") {
			// Update default configuration
			if err := config.UpdateDefaults(defaultMemory, defaultJavaArgs); err != nil {
				fmt.Fprintf(os.Stderr, "Error: Failed to update default configuration: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Default configuration updated successfully")
			return
		}

		// Otherwise, we need a server name and config type
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "Error: Server name and config type are required\n")
			cmd.Help()
			os.Exit(1)
		}

		serverName := args[0]
		configType := args[1]

		// Get the server configuration
		serverConfig, err := config.GetServer(serverName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Handle different config types
		switch configType {
		case "start":
			// Configure startup script
			configureStartupScript(serverConfig)
		case "properties":
			// Configure server.properties
			configureServerProperties(serverConfig)
		case "ops":
			// Configure ops.json
			configureOps(serverConfig)
		case "rcon":
			// Configure RCON settings
			if err := configureRcon(serverConfig, rconPort, rconPassword); err != nil {
				fmt.Fprintf(os.Stderr, "Error: Failed to configure RCON: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("RCON configuration updated successfully")
		default:
			fmt.Fprintf(os.Stderr, "Error: Unknown config type: %s\n", configType)
			cmd.Help()
			os.Exit(1)
		}
	},
}

// configureStartupScript opens the startup script in the user's editor
func configureStartupScript(serverConfig config.ServerConfig) {
	// Determine the startup script path
	var scriptPath string
	if runtime.GOOS == "windows" {
		scriptPath = filepath.Join(serverConfig.Path, "start.bat")
	} else {
		scriptPath = filepath.Join(serverConfig.Path, "start.sh")
	}

	// Check if the startup script exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Startup script does not exist: %s\n", scriptPath)
		os.Exit(1)
	}

	// Open the startup script in the user's editor
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("notepad", scriptPath)
	} else {
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "nano"
		}
		cmd = exec.Command(editor, scriptPath)
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to open editor: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Startup script updated successfully")
}

// configureServerProperties opens the server.properties file in the user's editor
func configureServerProperties(serverConfig config.ServerConfig) {
	// Determine the server.properties path
	propertiesPath := filepath.Join(serverConfig.Path, "server.properties")

	// Check if the server.properties file exists
	if _, err := os.Stat(propertiesPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: server.properties file does not exist: %s\n", propertiesPath)
		os.Exit(1)
	}

	// Open the server.properties file in the user's editor
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("notepad", propertiesPath)
	} else {
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "nano"
		}
		cmd = exec.Command(editor, propertiesPath)
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to open editor: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("server.properties updated successfully")
}

// configureOps opens the ops.json file in the user's editor
func configureOps(serverConfig config.ServerConfig) {
	// Determine the ops.json path
	opsPath := filepath.Join(serverConfig.Path, "ops.json")

	// Check if the ops.json file exists
	if _, err := os.Stat(opsPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: ops.json file does not exist: %s\n", opsPath)
		os.Exit(1)
	}

	// Open the ops.json file in the user's editor
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("notepad", opsPath)
	} else {
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "nano"
		}
		cmd = exec.Command(editor, opsPath)
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to open editor: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("ops.json updated successfully")
}

// configureRcon configures RCON settings in server.properties
func configureRcon(serverConfig config.ServerConfig, port int, password string) error {
	// Determine the server.properties path
	propertiesPath := filepath.Join(serverConfig.Path, "server.properties")

	// Check if the server.properties file exists
	if _, err := os.Stat(propertiesPath); os.IsNotExist(err) {
		return fmt.Errorf("server.properties file does not exist: %s", propertiesPath)
	}

	// Read the server.properties file
	content, err := os.ReadFile(propertiesPath)
	if err != nil {
		return fmt.Errorf("failed to read server.properties: %w", err)
	}

	// Update the RCON settings
	// This is a simple implementation that doesn't handle all cases
	// In a real implementation, we would use a proper properties parser
	lines := make([]string, 0)
	for _, line := range strings.Split(string(content), "\n") {
		if strings.HasPrefix(line, "enable-rcon=") {
			line = "enable-rcon=true"
		} else if strings.HasPrefix(line, "rcon.port=") {
			line = fmt.Sprintf("rcon.port=%d", port)
		} else if strings.HasPrefix(line, "rcon.password=") {
			line = fmt.Sprintf("rcon.password=%s", password)
		}
		lines = append(lines, line)
	}

	// Write the updated server.properties file
	if err := os.WriteFile(propertiesPath, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		return fmt.Errorf("failed to write server.properties: %w", err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Define flags for the config command
	configCmd.Flags().StringVar(&defaultMemory, "default-memory", "", "Default memory allocation for new servers")
	configCmd.Flags().StringVar(&defaultJavaArgs, "default-java-args", "", "Default Java arguments for new servers")
	configCmd.Flags().IntVar(&rconPort, "port", 25575, "RCON port")
	configCmd.Flags().StringVar(&rconPassword, "password", "", "RCON password")
}
