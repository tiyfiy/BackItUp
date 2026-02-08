package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tiyfiy/BackItUp/internal/config"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current configuration status",
	Long:  `Display the current configuration for all database connections and backup settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		showStatus()
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func showStatus() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		return
	}

	// Check if config file exists
	configExists := true
	if _, err := os.Stat("config.yaml"); os.IsNotExist(err) {
		configExists = false
		fmt.Println("⚠️  No config.yaml found - showing defaults")
		fmt.Println()
	}

	fmt.Println("📊 BackItUp Configuration Status")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// MongoDB Status
	fmt.Println("🍃 MongoDB")
	fmt.Println("───────────────────────────────────────────────────────────────")
	if cfg.MongoDB.URI != "" && cfg.MongoDB.URI != "mongodb://localhost:27017" {
		fmt.Printf("  URI:        %s\n", maskURI(cfg.MongoDB.URI))
		fmt.Printf("  Status:     ✅ Configured\n")
	} else {
		fmt.Printf("  URI:        %s (default)\n", cfg.MongoDB.URI)
		fmt.Printf("  Status:     ⚠️  Using defaults\n")
	}
	fmt.Println()

	// MySQL Status
	fmt.Println("🐬 MySQL")
	fmt.Println("───────────────────────────────────────────────────────────────")
	fmt.Printf("  Host:       %s\n", cfg.MySQL.Host)
	fmt.Printf("  Port:       %s\n", cfg.MySQL.Port)
	fmt.Printf("  User:       %s\n", cfg.MySQL.User)
	fmt.Printf("  Password:   %s\n", maskPassword(cfg.MySQL.Password))
	fmt.Printf("  Database:   %s\n", getValueOrDefault(cfg.MySQL.Database, "not set"))
	if cfg.MySQL.Database != "" {
		fmt.Printf("  Status:     ✅ Configured\n")
	} else {
		fmt.Printf("  Status:     ⚠️  Database not set\n")
	}
	fmt.Println()

	// PostgreSQL Status
	fmt.Println("🐘 PostgreSQL")
	fmt.Println("───────────────────────────────────────────────────────────────")
	fmt.Printf("  Host:       %s\n", cfg.PostgreSQL.Host)
	fmt.Printf("  Port:       %s\n", cfg.PostgreSQL.Port)
	fmt.Printf("  User:       %s\n", cfg.PostgreSQL.User)
	fmt.Printf("  Password:   %s\n", maskPassword(cfg.PostgreSQL.Password))
	fmt.Printf("  Database:   %s\n", getValueOrDefault(cfg.PostgreSQL.Database, "not set"))
	if cfg.PostgreSQL.Database != "" {
		fmt.Printf("  Status:     ✅ Configured\n")
	} else {
		fmt.Printf("  Status:     ⚠️  Database not set\n")
	}
	fmt.Println()

	// General Settings
	fmt.Println("⚙️  General Settings")
	fmt.Println("───────────────────────────────────────────────────────────────")
	fmt.Printf("  Backup Dir: %s\n", cfg.BackupDir)
	fmt.Printf("  Config:     %s\n", getConfigLocation(configExists))
	fmt.Println()

	// Quick tips
	if !configExists {
		fmt.Println("💡 Tip: Run a database config command to create config.yaml")
		fmt.Println("   Example: ./BackItUp mysql --config --database mydb")
	}
}

func maskPassword(password string) string {
	if password == "" {
		return "not set"
	}
	return "********"
}

func maskURI(uri string) string {
	// Simple masking - replace password in MongoDB URI
	// Format: mongodb://user:password@host or mongodb+srv://user:password@host
	if len(uri) > 20 {
		// Just show beginning and end for security
		return uri[:15] + "..." + uri[len(uri)-20:]
	}
	return "***masked***"
}

func getValueOrDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func getConfigLocation(exists bool) string {
	if exists {
		return "config.yaml ✅"
	}
	return "not found ⚠️"
}
