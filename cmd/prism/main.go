package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version information (set during build)
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "prism",
	Short: "PRISM - Phase Render & Inspection for Structural Mockups",
	Long: `PRISM is a CLI tool that generates visual PNG representations 
of Phase 1 structural mockups from your AI Design Agent process.

It takes JSON structure files created in Phase 1 and renders them as 
black-and-white wireframe images for easy review and approval.`,
	Version: version,
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().Bool("json", false, "Output in JSON format")
	rootCmd.PersistentFlags().StringP("project", "p", "./", "Project directory path")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Suppress non-essential output")
	rootCmd.PersistentFlags().String("config", "", "Config file path (default: ~/.prism)")

	// Add subcommands
	rootCmd.AddCommand(renderCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(compareCmd)
}
