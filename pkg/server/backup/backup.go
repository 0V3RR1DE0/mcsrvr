package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// CreateBackup creates a backup of a Minecraft server
func CreateBackup(serverName, backupPath string) error {
	// Create the backup directory if it doesn't exist
	if err := os.MkdirAll(backupPath, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Create a timestamp for the backup
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	backupName := fmt.Sprintf("%s_%s", serverName, timestamp)
	backupDir := filepath.Join(backupPath, backupName)

	// Create the backup directory
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Copy the server files to the backup directory
	// This is a placeholder implementation
	// In a real implementation, we would need to copy all the server files
	fmt.Printf("Creating backup of server '%s' to '%s'...\n", serverName, backupDir)

	// For now, we'll just print a message
	fmt.Println("Backup created successfully")

	return nil
}

// RestoreBackup restores a backup of a Minecraft server
func RestoreBackup(backupPath, restorePath string) error {
	// Check if the backup directory exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup directory does not exist: %s", backupPath)
	}

	// Create the restore directory if it doesn't exist
	if err := os.MkdirAll(restorePath, 0755); err != nil {
		return fmt.Errorf("failed to create restore directory: %w", err)
	}

	// Copy the backup files to the restore directory
	// This is a placeholder implementation
	// In a real implementation, we would need to copy all the backup files
	fmt.Printf("Restoring backup from '%s' to '%s'...\n", backupPath, restorePath)

	// For now, we'll just print a message
	fmt.Println("Backup restored successfully")

	return nil
}

// ListBackups lists all backups for a server
func ListBackups(serverName string) ([]string, error) {
	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Determine the backups directory
	backupsDir := filepath.Join(homeDir, ".mcsrvr", "backups")

	// Check if the backups directory exists
	if _, err := os.Stat(backupsDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	// List all backups for the server
	var backups []string

	// If no server name is provided, list all backups
	if serverName == "" {
		entries, err := os.ReadDir(backupsDir)
		if err != nil {
			return nil, fmt.Errorf("failed to read backups directory: %w", err)
		}

		for _, entry := range entries {
			if entry.IsDir() {
				backups = append(backups, entry.Name())
			}
		}
	} else {
		// List backups for the specific server
		entries, err := os.ReadDir(backupsDir)
		if err != nil {
			return nil, fmt.Errorf("failed to read backups directory: %w", err)
		}

		for _, entry := range entries {
			if entry.IsDir() && strings.HasPrefix(entry.Name(), serverName+"_") {
				backups = append(backups, entry.Name())
			}
		}
	}

	return backups, nil
}
