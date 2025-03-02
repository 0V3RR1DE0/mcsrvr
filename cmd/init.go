package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/0v3rr1de0/mcsrvr/pkg/config"
	"github.com/0v3rr1de0/mcsrvr/pkg/server"
)

var (
	serverName         string
	serverType         string
	serverVersion      string
	serverMemory       string
	serverJavaArgs     string
	fabricLoaderVersion string
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Initialize a new Minecraft server",
	Long: `Initialize a new Minecraft server at the specified path.
If no path is provided, the current directory will be used.

Example:
  mcsrvr init . -n paper123 --type papermc -v 1.21.4
  mcsrvr init D:/serverfolder -n vanilla123 --type vanilla -v 1.21.4
  mcsrvr init D:/serverfolder -n fabric123 --type fabric -v 1.21.4 --fabric-loader 0.16.10`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Determine the server path
		serverPath := "."
		if len(args) > 0 {
			serverPath = args[0]
		}

		// Validate server type
		validTypes := map[string]bool{
			"papermc": true,
			"vanilla": true,
			"fabric":  true,
			// Add other server types as they are implemented
		}

		if !validTypes[serverType] {
			fmt.Fprintf(os.Stderr, "Error: Invalid server type '%s'. Supported types: papermc, vanilla, fabric\n", serverType)
			os.Exit(1)
		}

		// Create the server directory if it doesn't exist
		serverPath, err := filepath.Abs(serverPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to resolve absolute path: %v\n", err)
			os.Exit(1)
		}

		if err := os.MkdirAll(serverPath, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to create server directory: %v\n", err)
			os.Exit(1)
		}

		// Get default configuration if needed
		if serverMemory == "" || serverJavaArgs == "" {
			defaults, err := config.GetDefaults()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Failed to get default configuration: %v\n", err)
				os.Exit(1)
			}

			if serverMemory == "" {
				serverMemory = defaults.Memory
			}
			if serverJavaArgs == "" {
				serverJavaArgs = defaults.JavaArgs
			}
		}

		fmt.Printf("Initializing %s server '%s' at %s with version %s\n", serverType, serverName, serverPath, serverVersion)
		fmt.Printf("Memory: %s, Java Args: %s\n", serverMemory, serverJavaArgs)

		// Initialize the server
		var initErr error
		if serverType == "fabric" {
			// For Fabric servers, pass the loader version
			initErr = server.InitializeFabricServer(serverPath, serverName, serverVersion, fabricLoaderVersion, serverMemory, serverJavaArgs)
		} else {
			// For other server types
			initErr = server.InitializeServer(serverPath, serverName, serverType, serverVersion, serverMemory, serverJavaArgs)
		}
		
		if initErr != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to initialize server: %v\n", initErr)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Define flags for the init command
	initCmd.Flags().StringVarP(&serverName, "name", "n", "", "Name of the server (required)")
	initCmd.Flags().StringVar(&serverType, "type", "", "Server type (papermc, vanilla, fabric, etc.) (required)")
	initCmd.Flags().StringVarP(&serverVersion, "version", "v", "latest", "Server version")
	initCmd.Flags().StringVarP(&serverMemory, "memory", "m", "2G", "Memory allocation for the server (e.g., 2G, 4G)")
	initCmd.Flags().StringVar(&serverJavaArgs, "java-args", "", "Additional Java arguments")
	initCmd.Flags().StringVar(&fabricLoaderVersion, "fabric-loader", "0.16.10", "Fabric loader version (only for fabric server type)")

	// Mark required flags
	initCmd.MarkFlagRequired("name")
	initCmd.MarkFlagRequired("type")
	
	// Initialize the configuration
	config.Initialize()
}
