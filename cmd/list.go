package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"github.com/0v3rr1de0/mcsrvr/pkg/config"
	"github.com/0v3rr1de0/mcsrvr/pkg/server/process"
)

var (
	onlineOnly  bool
	offlineOnly bool
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Minecraft servers",
	Long: `List all Minecraft servers managed by mcsrvr.
You can filter the list to show only online or offline servers.

Example:
  mcsrvr list
  mcsrvr list --online
  mcsrvr list --offline`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get all servers
		servers, err := config.ListServers()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to list servers: %v\n", err)
			os.Exit(1)
		}

		if len(servers) == 0 {
			fmt.Println("No servers found. Use 'mcsrvr init' to create a new server.")
			return
		}

		// Create a tabwriter for formatted output
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tTYPE\tVERSION\tPATH\tSTATUS\tPID\tLAST STARTED")

		// Print each server
		for _, srv := range servers {
			// Determine server status
			status := "Offline"
			if proc, exists := process.ActiveServers[srv.Name]; exists && proc.Running {
				status = "Online"
			}
			
			// Skip if filtering by status
			if onlineOnly && status != "Online" {
				continue
			}
			if offlineOnly && status != "Offline" {
				continue
			}

			// Format the last started time
			lastStarted := "Never"
			if !srv.LastStarted.IsZero() {
				lastStarted = srv.LastStarted.Format(time.RFC1123)
			}

			// Get PID if server is running
			pid := "-"
			if proc, exists := process.ActiveServers[srv.Name]; exists && proc.Running {
				pid = fmt.Sprintf("%d", proc.PID)
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
				srv.Name, srv.Type, srv.Version, srv.Path, status, pid, lastStarted)
		}

		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Define flags for the list command
	listCmd.Flags().BoolVar(&onlineOnly, "online", false, "Show only online servers")
	listCmd.Flags().BoolVar(&offlineOnly, "offline", false, "Show only offline servers")
}
