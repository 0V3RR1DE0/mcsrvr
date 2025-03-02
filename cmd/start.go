package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/0v3rr1de0/mcsrvr/pkg/server"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start [server-name]",
	Short: "Start a Minecraft server",
	Long: `Start a Minecraft server by name.

Example:
  mcsrvr start paper123`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serverName := args[0]

		// Start the server
		if err := server.StartServer(serverName); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to start server: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
