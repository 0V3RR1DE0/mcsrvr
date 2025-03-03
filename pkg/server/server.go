package server

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/0v3rr1de0/mcsrvr/pkg/config"
	"github.com/0v3rr1de0/mcsrvr/pkg/server/backup"
	serverInit "github.com/0v3rr1de0/mcsrvr/pkg/server/init"
	"github.com/0v3rr1de0/mcsrvr/pkg/server/process"
	"github.com/0v3rr1de0/mcsrvr/pkg/server/rcon"
)

// StartServer starts a Minecraft server as a detached process
func StartServer(serverName string) error {
	// Get the server configuration
	serverConfig, err := config.GetServer(serverName)
	if err != nil {
		return err
	}

	// Check if the server is already running
	if proc, exists := process.ActiveServers[serverName]; exists && proc.Running {
		return fmt.Errorf("server '%s' is already running", serverName)
	}

	// Check if the server directory exists
	if _, err := os.Stat(serverConfig.Path); os.IsNotExist(err) {
		return fmt.Errorf("server directory does not exist: %s", serverConfig.Path)
	}

	// Determine the startup script path
	var scriptPath string
	if runtime.GOOS == "windows" {
		scriptPath = filepath.Join(serverConfig.Path, "start.bat")
	} else {
		scriptPath = filepath.Join(serverConfig.Path, "start.sh")
	}

	// Check if the startup script exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("startup script does not exist: %s", scriptPath)
	}

	// Create log file for the server
	logDir := filepath.Join(serverConfig.Path, "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	logFile, err := os.OpenFile(
		filepath.Join(logDir, "server.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}
	defer logFile.Close()

	// Create the command to start the server
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		// On Windows, use cmd.exe to run the batch file
		cmd = exec.Command("cmd", "/c", scriptPath)
	} else {
		// For Unix-like systems, use bash to run the shell script
		cmd = exec.Command("bash", scriptPath)
	}

	// Set the working directory and output files
	cmd.Dir = serverConfig.Path
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	// Set the process attributes using our helper function
	cmd.SysProcAttr = process.NewSysProcAttr()

	// Start the server process
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	// Store the command window process PID
	cmdPID := cmd.Process.Pid
	fmt.Printf("Server '%s' started with command window PID %d\n", serverName, cmdPID)

	// Release the process to allow the CLI to exit without killing the server
	if err := cmd.Process.Release(); err != nil {
		fmt.Printf("Warning: Failed to release process: %v\n", err)
	}

	// For Unix, we need to find the Java process PID
	if runtime.GOOS != "windows" {
		// Wait a moment for the Java process to start
		time.Sleep(2 * time.Second)

		// Find the Java process PID
		javaPID, err := process.FindJavaPID()
		if err != nil {
			fmt.Printf("Warning: Failed to find Java process PID: %v\n", err)
			fmt.Println("Using command window PID instead")
		} else {
			fmt.Printf("Found Java process with PID %d\n", javaPID)
			cmdPID = javaPID
		}
	}

	// Store the process in ActiveServers
	process.ActiveServers[serverName] = &process.ServerProcess{
		Name:    serverName,
		PID:     cmdPID,
		Running: true,
	}

	// Save the active servers to file
	if err := process.SaveActiveServers(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to save active servers: %v\n", err)
	}

	// Update the server configuration with the last started time
	serverConfig.LastStarted = time.Now()
	if err := config.UpdateServer(serverName, serverConfig); err != nil {
		return fmt.Errorf("failed to update server configuration: %w", err)
	}

	// We've already released the process earlier, so no need to do it again

	fmt.Printf("Server '%s' started successfully\n", serverName)
	return nil
}

// StopServer stops a Minecraft server
func StopServer(serverName string) error {
	// Get the server configuration
	_, err := config.GetServer(serverName)
	if err != nil {
		return err
	}

	// Check if the server is running
	proc, exists := process.ActiveServers[serverName]
	if !exists || !proc.Running {
		return fmt.Errorf("server '%s' is not running", serverName)
	}

	fmt.Printf("Stopping server '%s'...\n", serverName)

	// Store the Java process PID
	javaPID := proc.PID

	// Try to stop the server gracefully using RCON
	if err := rcon.StopServerGracefully(serverName); err != nil {
		fmt.Printf("Warning: Failed to stop server gracefully: %v\n", err)
		fmt.Println("Falling back to force kill...")
	} else {
		fmt.Println("RCON stop command sent successfully")

		// Wait a moment for the server to start shutting down
		time.Sleep(2 * time.Second)
	}

	// Check if the Java process is still running
	javaRunning := process.IsProcessRunning(javaPID)

	// If the Java process is still running, kill it
	if javaRunning {
		fmt.Printf("Java process (PID %d) is still running, attempting to kill it...\n", javaPID)

		if runtime.GOOS == "windows" {
			// On Windows, use taskkill to kill the Java process
			killCmd := exec.Command("taskkill", "/F", "/PID", fmt.Sprintf("%d", javaPID))
			if err := killCmd.Run(); err != nil {
				fmt.Printf("Warning: Failed to kill Java process: %v\n", err)
			} else {
				fmt.Printf("Java process (PID %d) killed successfully\n", javaPID)
			}
		} else {
			// On Unix-like systems, use kill to kill the Java process
			killCmd := exec.Command("kill", "-9", fmt.Sprintf("%d", javaPID))
			if err := killCmd.Run(); err != nil {
				fmt.Printf("Warning: Failed to kill Java process: %v\n", err)
			} else {
				fmt.Printf("Java process (PID %d) killed successfully\n", javaPID)
			}
		}
	}

	// Update the process status
	proc.Running = false
	delete(process.ActiveServers, serverName)

	// Save the active servers to file
	if err := process.SaveActiveServers(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to save active servers: %v\n", err)
	}

	fmt.Printf("Server '%s' stopped successfully\n", serverName)

	return nil
}

// ExecuteCommand executes a command on a Minecraft server using RCON.
func ExecuteCommand(serverName, command string) error {
	// Get the server configuration.
	_, err := config.GetServer(serverName)
	if err != nil {
		return err
	}

	// Check if the server is running.
	if _, exists := process.ActiveServers[serverName]; !exists {
		return fmt.Errorf("server '%s' is not running", serverName)
	}

	return rcon.ExecuteCommand(serverName, command)
}

// CreateBackup creates a backup of a Minecraft server
func CreateBackup(serverName, backupPath string) error {
	// Get the server configuration
	_, err := config.GetServer(serverName)
	if err != nil {
		return err
	}

	return backup.CreateBackup(serverName, backupPath)
}

// RestoreBackup restores a backup of a Minecraft server
func RestoreBackup(backupPath, restorePath string) error {
	return backup.RestoreBackup(backupPath, restorePath)
}

// ListBackups lists all backups for a server
func ListBackups(serverName string) ([]string, error) {
	return backup.ListBackups(serverName)
}

// InitializeServer initializes a new Minecraft server
func InitializeServer(serverPath, serverName, serverType, version, memory, javaArgs string) error {
	return serverInit.InitializeServer(serverPath, serverName, serverType, version, memory, javaArgs)
}

// InitializeFabricServer initializes a new Fabric server
func InitializeFabricServer(serverPath, serverName, mcVersion, loaderVersion, memory, javaArgs string) error {
	return serverInit.InitializeFabricServer(serverPath, serverName, mcVersion, loaderVersion, memory, javaArgs)
}

// AcceptEULA accepts the Minecraft EULA by creating or modifying the eula.txt file
func AcceptEULA(serverPath string) error {
	return serverInit.AcceptEULA(serverPath)
}

// RefreshServerStatus updates the status of all active servers
func RefreshServerStatus() {
	process.RefreshServerStatus()
}

// init loads the active servers when the package is initialized
func init() {
	// The initialization is now handled in the process package
}
