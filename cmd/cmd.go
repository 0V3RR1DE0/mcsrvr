package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/0v3rr1de0/mcsrvr/pkg/server"
)

// cmdCmd represents the cmd command
var cmdCmd = &cobra.Command{
	Use:   "cmd [server-name] [command]",
	Short: "Execute a command on a Minecraft server",
	Long: `Execute a command on a Minecraft server by name.
The command should be provided without the leading slash (/).

Example:
  mcsrvr cmd paper123 "say Hello, world!"
  mcsrvr cmd paper123 "op username"`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		serverName := args[0]
		command := args[1]

		// Execute the command
		if err := server.ExecuteCommand(serverName, command); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to execute command: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(cmdCmd)
}
