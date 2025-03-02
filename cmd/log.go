package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/0v3rr1de0/mcsrvr/pkg/config"
)

var (
	follow bool
	lines  int
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log [server-name]",
	Short: "View server logs",
	Long: `View the logs of a Minecraft server.
You can follow the logs in real-time or view a specific number of lines.

Example:
  mcsrvr log paper123
  mcsrvr log paper123 --follow
  mcsrvr log paper123 --lines 50`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serverName := args[0]

		// Get the server configuration
		serverConfig, err := config.GetServer(serverName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Determine the log file path
		logPath := filepath.Join(serverConfig.Path, "logs", "latest.log")

		// Check if the log file exists
		if _, err := os.Stat(logPath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: Log file does not exist: %s\n", logPath)
			os.Exit(1)
		}

		// Open the log file
		file, err := os.Open(logPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to open log file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		if follow {
			// Set up signal handling to avoid stopping the MC server
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

			go func() {
				<-sigChan
				fmt.Println("\nLog monitoring stopped.")
				os.Exit(0)
			}()

			// Follow logs in real-time
			fmt.Printf("Following logs for server '%s'...\n", serverName)
			fmt.Println("Press Ctrl+C to stop following logs")

			reader := bufio.NewReader(file)
			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					time.Sleep(500 * time.Millisecond)
					continue
				}
				fmt.Print(line)
			}
		} else {
			// Print the last 'lines' lines of the log file
			if lines <= 0 {
				lines = 10 // Default to 10 lines
			}

			scanner := bufio.NewScanner(file)
			var logLines []string

			for scanner.Scan() {
				logLines = append(logLines, scanner.Text())
				if len(logLines) > lines {
					logLines = logLines[1:]
				}
			}

			for _, line := range logLines {
				fmt.Println(line)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(logCmd)

	// Define flags for the log command
	logCmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow the log file in real-time")
	logCmd.Flags().IntVarP(&lines, "lines", "n", 10, "Number of lines to show")
}
