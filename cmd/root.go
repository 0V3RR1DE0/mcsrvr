package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mcsrvr",
	Short: "A CLI tool for managing Minecraft servers",
	Long: `mcsrvr is a CLI tool for easily setting up and managing Minecraft servers.
It supports various server types like vanilla, PaperMC, Spigot, Bukkit, Velocity, 
Forge, Fabric, BungeeCord, and Cuberite.

You can initialize, start, stop, backup, and manage your Minecraft servers with simple commands.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mcsrvr.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
