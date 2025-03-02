package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/0v3rr1de0/mcsrvr/pkg/config"
)

var (
	forceDelete bool
)

// delCmd represents the del command
var delCmd = &cobra.Command{
	Use:   "del [server-name]",
	Short: "Delete a Minecraft server",
	Long: `Delete a Minecraft server by name.
This will remove the server from the configuration, but will not delete the server files.

Example:
  mcsrvr del paper123
  mcsrvr del paper123 -y`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serverName := args[0]

		// Get the server configuration
		serverConfig, err := config.GetServer(serverName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Confirm deletion if not forced
		if !forceDelete {
			fmt.Printf("Are you sure you want to delete server '%s'? This will remove the server from the configuration, but will not delete the server files. (y/N): ", serverName)
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
				os.Exit(1)
			}

			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				fmt.Println("Deletion cancelled.")
				return
			}
		}

		// Delete the server
		if err := config.DeleteServer(serverName); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to delete server: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Server '%s' deleted from configuration. Server files at '%s' were not deleted.\n", serverName, serverConfig.Path)
	},
}

func init() {
	rootCmd.AddCommand(delCmd)

	// Define flags for the del command
	delCmd.Flags().BoolVarP(&forceDelete, "yes", "y", false, "Skip confirmation prompt")
}
