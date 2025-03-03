package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/0v3rr1de0/mcsrvr/pkg/config"
	"github.com/0v3rr1de0/mcsrvr/pkg/server"
	"github.com/0v3rr1de0/mcsrvr/pkg/server/process"
)

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "restart [server-name]",
	Short: "Restart a Minecraft server",
	Long: `Restart a Minecraft server by name.

Example:
  mcsrvr restart paper123`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serverName := args[0]

		// Get the server configuration
		_, err := config.GetServer(serverName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Check if the server is running
		proc, exists := process.ActiveServers[serverName]
		if !exists || !proc.Running {
			fmt.Fprintf(os.Stderr, "Error: Server '%s' is not running\n", serverName)
			os.Exit(1)
		}

		// Stop the server
		fmt.Printf("Stopping server '%s'...\n", serverName)
		if err := server.StopServer(serverName); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to stop server: %v\n", err)
			os.Exit(1)
		}

		// Wait a moment for the server to fully stop
		fmt.Println("Waiting for server to stop...")
		time.Sleep(5 * time.Second)

		// Verify that the server is fully stopped
		if proc, exists := process.ActiveServers[serverName]; exists && proc.Running {
			fmt.Fprintf(os.Stderr, "Error: Server '%s' is still running after stop command\n", serverName)
			os.Exit(1)
		}

		// Start the server
		fmt.Printf("Starting server '%s'...\n", serverName)
		if err := server.StartServer(serverName); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to start server: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Server '%s' restarted successfully\n", serverName)
	},
}

func init() {
	rootCmd.AddCommand(restartCmd)
}
