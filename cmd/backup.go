package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/0v3rr1de0/mcsrvr/pkg/server"
)

var (
	backupPath string
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup [server-name]",
	Short: "Create a backup of a Minecraft server",
	Long: `Create a backup of a Minecraft server by name.
If no backup path is provided, the backup will be created in the default backup directory.

Example:
  mcsrvr backup paper123
  mcsrvr backup paper123 --path D:/backups`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serverName := args[0]

		// If no backup path is provided, use the default
		if backupPath == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Failed to get user home directory: %v\n", err)
				os.Exit(1)
			}
			backupPath = filepath.Join(homeDir, ".mcsrvr", "backups")
		}

		// Create the backup
		if err := server.CreateBackup(serverName, backupPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to create backup: %v\n", err)
			os.Exit(1)
		}
	},
}

// backupsCmd represents the backups command
var backupsCmd = &cobra.Command{
	Use:   "backups [server-name]",
	Short: "List backups for a Minecraft server",
	Long: `List backups for a Minecraft server by name.
If no server name is provided, all backups will be listed.

Example:
  mcsrvr backups
  mcsrvr backups paper123`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var serverName string
		if len(args) > 0 {
			serverName = args[0]
		}

		// List the backups
		backups, err := server.ListBackups(serverName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to list backups: %v\n", err)
			os.Exit(1)
		}

		if len(backups) == 0 {
			if serverName == "" {
				fmt.Println("No backups found.")
			} else {
				fmt.Printf("No backups found for server '%s'.\n", serverName)
			}
			return
		}

		// Print the backups
		fmt.Println("Backups:")
		for _, backup := range backups {
			fmt.Println(backup)
		}
	},
}

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore [backup-path] [restore-path]",
	Short: "Restore a backup of a Minecraft server",
	Long: `Restore a backup of a Minecraft server to a new location.

Example:
  mcsrvr restore ~/.mcsrvr/backups/paper123_2025-03-02_00-00-00 D:/servers/restored_paper123`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		backupPath := args[0]
		restorePath := args[1]

		// Restore the backup
		if err := server.RestoreBackup(backupPath, restorePath); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to restore backup: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(backupsCmd)
	rootCmd.AddCommand(restoreCmd)

	// Define flags for the backup command
	backupCmd.Flags().StringVar(&backupPath, "path", "", "Path to store the backup")
}
