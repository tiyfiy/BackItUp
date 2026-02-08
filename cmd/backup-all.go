package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/tiyfiy/BackItUp/internal/config"
	"github.com/tiyfiy/BackItUp/internal/mongodb"
	"github.com/tiyfiy/BackItUp/internal/mysql"
	"github.com/tiyfiy/BackItUp/internal/postgresql"
)

var backupAllCmd = &cobra.Command{
	Use:   "backup-all",
	Short: "Backup all configured databases at once",
	Long: `Run backups for all configured databases (MongoDB, MySQL, PostgreSQL) in sequence.

This command will:
- Check which databases are configured
- Run backups for each configured database
- Report success/failure for each
- Show total time taken`,
	Run: func(cmd *cobra.Command, args []string) {
		runBackupAll()
	},
}

func init() {
	rootCmd.AddCommand(backupAllCmd)
}

func runBackupAll() {
	startTime := time.Now()

	fmt.Println("🔄 Starting backup for all configured databases...")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	cfg, err := config.Load()
	if err != nil {
		fmt.Println("❌ Error loading configuration:", err)
		return
	}

	successCount := 0
	failCount := 0
	skippedCount := 0

	// Backup MongoDB
	if cfg.MongoDB.URI != "" && cfg.MongoDB.URI != "mongodb://localhost:27017" {
		fmt.Println("📦 Backing up MongoDB...")
		client, err := mongodb.Connection(cfg.MongoDB.URI)
		if err != nil {
			fmt.Printf("   ❌ Failed: %v\n\n", err)
			failCount++
		} else {
			mongodb.Backup(client, cfg.MongoDB.URI)
			successCount++
			fmt.Println()
		}
	} else {
		fmt.Println("⏭️  Skipping MongoDB (not configured)")
		skippedCount++
	}

	// Backup MySQL
	if cfg.MySQL.Database != "" {
		fmt.Println("📦 Backing up MySQL...")
		db, err := mysql.Connection(cfg.MySQL.Host, cfg.MySQL.Port, cfg.MySQL.User, cfg.MySQL.Password, cfg.MySQL.Database)
		if err != nil {
			fmt.Printf("   ❌ Failed: %v\n\n", err)
			failCount++
		} else {
			mysql.Backup(db, cfg.MySQL.Host, cfg.MySQL.Port, cfg.MySQL.User, cfg.MySQL.Password, cfg.MySQL.Database)
			db.Close()
			successCount++
			fmt.Println()
		}
	} else {
		fmt.Println("⏭️  Skipping MySQL (not configured)")
		skippedCount++
	}

	// Backup PostgreSQL
	if cfg.PostgreSQL.Database != "" {
		fmt.Println("📦 Backing up PostgreSQL...")
		db, err := postgresql.Connection(cfg.PostgreSQL.Host, cfg.PostgreSQL.Port, cfg.PostgreSQL.User, cfg.PostgreSQL.Password, cfg.PostgreSQL.Database)
		if err != nil {
			fmt.Printf("   ❌ Failed: %v\n\n", err)
			failCount++
		} else {
			postgresql.Backup(db, cfg.PostgreSQL.Host, cfg.PostgreSQL.Port, cfg.PostgreSQL.User, cfg.PostgreSQL.Password, cfg.PostgreSQL.Database)
			db.Close()
			successCount++
			fmt.Println()
		}
	} else {
		fmt.Println("⏭️  Skipping PostgreSQL (not configured)")
		skippedCount++
	}

	// Summary
	duration := time.Since(startTime)
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Printf("📊 Backup Summary:\n")
	fmt.Printf("   ✅ Successful: %d\n", successCount)
	if failCount > 0 {
		fmt.Printf("   ❌ Failed: %d\n", failCount)
	}
	if skippedCount > 0 {
		fmt.Printf("   ⏭️  Skipped: %d\n", skippedCount)
	}
	fmt.Printf("   ⏱️  Duration: %.2f seconds\n", duration.Seconds())
	fmt.Println()

	if successCount == 0 {
		fmt.Println("💡 Tip: Configure databases first using the --config flag")
		fmt.Println("   Example: ./BackItUp mysql --config --database mydb")
	}
}
