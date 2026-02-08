// Package cmd implements the command-line interface for BackItUp.
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "BackItUp",
	Short: "Simple CLI tool for backing up databases",
	Long: `BackItUp is a simple yet powerful CLI tool for backing up your databases.

Currently supports:
  - MongoDB
  - MySQL
  - PostgreSQL

Use the database-specific subcommands to configure and run backups.
All backups are stored in the BACKUP/ directory organized by database type.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Root command flags can be added here if needed
}
