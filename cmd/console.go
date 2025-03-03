package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/0v3rr1de0/mcsrvr/pkg/config"
	"github.com/0v3rr1de0/mcsrvr/pkg/server"
	"github.com/0v3rr1de0/mcsrvr/pkg/server/process"
)

// consoleCmd represents the console command
var consoleCmd = &cobra.Command{
	Use:   "console [server-name]",
	Short: "Access the console of a Minecraft server",
	Long: `Access the console of a Minecraft server by name.
You can exit the console by typing 'exit'.

Example:
  mcsrvr console paper123`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serverName := args[0]

		// Get the server configuration.
		serverConfig, err := config.GetServer(serverName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Check if the server is running.
		if _, exists := process.ActiveServers[serverName]; !exists {
			fmt.Fprintf(os.Stderr, "Error: Server '%s' is not running. Start it first with 'mcsrvr start %s'\n", serverName, serverName)
			os.Exit(1)
		}

		// Determine the appropriate key combination based on OS
		exitKey := "Ctrl+C"
		if runtime.GOOS == "darwin" {
			exitKey = "Cmd+C"
		}

		fmt.Printf("Connecting to console of server '%s'...\n", serverName)
		fmt.Printf("Type 'exit' or press %s to exit the console.\n", exitKey)

		// Channel to signal when to exit the console.
		exitChan := make(chan bool)

		// Start a goroutine to follow the server log file.
		go func() {
			// Determine the log file path.
			logPath := filepath.Join(serverConfig.Path, "logs", "latest.log")
			
			// Wait for the log file to be created if it doesn't exist.
			for {
				if _, err := os.Stat(logPath); !os.IsNotExist(err) {
					break
				}
				time.Sleep(500 * time.Millisecond)
			}

			// Open the log file.
			file, err := os.Open(logPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Failed to open log file: %v\n", err)
				return
			}
			defer file.Close()

			// Seek to the end of the file to only show new logs.
			fileInfo, err := file.Stat()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Failed to get file info: %v\n", err)
				return
			}
			file.Seek(fileInfo.Size(), 0)
			reader := bufio.NewReader(file)
			
			// Continuously follow the file.
			for {
				select {
				case <-exitChan:
					return
				default:
					line, err := reader.ReadString('\n')
					if err != nil {
						time.Sleep(100 * time.Millisecond)
						continue
					}
					fmt.Print(line)
					// Print the console prompt after each log line
					fmt.Print("> ")
				}
			}
		}()

		// Start a goroutine to handle user input.
		go func() {
			// Print initial console prompt
			fmt.Print("> ")
			
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				input := scanner.Text()
				
				// If the user types "exit", signal exit and do not execute any command.
				if strings.ToLower(input) == "exit" {
					exitChan <- true
					return
				}
				
				// Execute the command on the server.
				if err := server.ExecuteCommand(serverName, input); err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				}
				
				// Print the console prompt again after command execution
				fmt.Print("> ")
			}
			if err := scanner.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
				exitChan <- true
			}
		}()

		// Wait for the exit signal from the user.
		<-exitChan
		fmt.Println("Exiting console...")
	},
}

func init() {
	rootCmd.AddCommand(consoleCmd)
}
