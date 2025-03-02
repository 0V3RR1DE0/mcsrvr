package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/0v3rr1de0/mcsrvr/pkg/server"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop [server-name]",
	Short: "Stop a Minecraft server",
	Long: `Stop a Minecraft server by name.

Example:
  mcsrvr stop paper123`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serverName := args[0]

		// Stop the server
		if err := server.StopServer(serverName); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to stop server: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
