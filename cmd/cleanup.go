package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	keepDays  int
	keepCount int
	dryRun    bool
)

var cleanupCmd = &cobra.Command{
	Use:   "cleanup [mongodb|mysql|postgresql|all]",
	Short: "Clean up old backups based on retention policy",
	Long: `Remove old backups to save disk space.

You can specify retention by:
  --days N    Keep backups from last N days
  --keep N    Keep the N most recent backups
  --dry-run   Show what would be deleted without actually deleting

Examples:
  ./BackItUp cleanup mysql --days 30          # Keep last 30 days
  ./BackItUp cleanup postgresql --keep 5      # Keep 5 most recent
  ./BackItUp cleanup all --days 7 --dry-run   # Preview cleanup for all DBs`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]

		if keepDays == 0 && keepCount == 0 {
			fmt.Println("Error: You must specify either --days or --keep")
			fmt.Println("Example: ./BackItUp cleanup mysql --days 30")
			return
		}

		cleanupBackups(target)
	},
}

func init() {
	rootCmd.AddCommand(cleanupCmd)
	cleanupCmd.Flags().IntVarP(&keepDays, "days", "d", 0, "Keep backups from last N days")
	cleanupCmd.Flags().IntVarP(&keepCount, "keep", "k", 0, "Keep N most recent backups")
	cleanupCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be deleted without deleting")
}

func cleanupBackups(target string) {
	if dryRun {
		fmt.Println("🔍 DRY RUN MODE - No files will be deleted")
		fmt.Println()
	}

	var totalDeleted int
	var totalSize int64

	switch target {
	case "mongodb":
		deleted, size := cleanupDatabaseBackups("mongo", "MongoDB", true)
		totalDeleted += deleted
		totalSize += size

	case "mysql":
		deleted, size := cleanupDatabaseBackups("mysql", "MySQL", false)
		totalDeleted += deleted
		totalSize += size

	case "postgresql":
		deleted, size := cleanupDatabaseBackups("postgresql", "PostgreSQL", false)
		totalDeleted += deleted
		totalSize += size

	case "all":
		fmt.Println("🧹 Cleaning up backups for all databases...")
		fmt.Println()

		deleted, size := cleanupDatabaseBackups("mongo", "MongoDB", true)
		totalDeleted += deleted
		totalSize += size

		deleted, size = cleanupDatabaseBackups("mysql", "MySQL", false)
		totalDeleted += deleted
		totalSize += size

		deleted, size = cleanupDatabaseBackups("postgresql", "PostgreSQL", false)
		totalDeleted += deleted
		totalSize += size

	default:
		fmt.Printf("Unknown target: %s\n", target)
		fmt.Println("Supported targets: mongodb, mysql, postgresql, all")
		return
	}

	// Summary
	fmt.Println("\n" + strings.Repeat("═", 60))
	if dryRun {
		fmt.Printf("Would delete: %d backup(s), freeing %s\n", totalDeleted, formatSize(totalSize))
		fmt.Println("\nRun without --dry-run to actually delete these backups.")
	} else {
		if totalDeleted > 0 {
			fmt.Printf("✅ Deleted: %d backup(s), freed %s\n", totalDeleted, formatSize(totalSize))
		} else {
			fmt.Println("✅ No backups needed cleanup")
		}
	}
}

func cleanupDatabaseBackups(dbDir, dbName string, isDirectory bool) (int, int64) {
	backupPath := filepath.Join("BACKUP", dbDir)

	// Check if backup directory exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return 0, 0
	}

	var backups []BackupInfo

	if isDirectory {
		backups = listDirBackups(backupPath)
	} else {
		backups = listFileBackups(backupPath, ".sql")
	}

	if len(backups) == 0 {
		return 0, 0
	}

	fmt.Printf("📦 %s Backups\n", dbName)
	fmt.Println("──────────────────────────────────────────────────────────")

	// Determine which backups to delete
	var toDelete []BackupInfo

	if keepDays > 0 {
		cutoffDate := time.Now().AddDate(0, 0, -keepDays)
		for _, backup := range backups {
			if backup.ModTime.Before(cutoffDate) {
				toDelete = append(toDelete, backup)
			}
		}
	} else if keepCount > 0 {
		if len(backups) > keepCount {
			toDelete = backups[keepCount:]
		}
	}

	if len(toDelete) == 0 {
		fmt.Printf("   No old backups to clean (keeping ")
		if keepDays > 0 {
			fmt.Printf("last %d days)\n", keepDays)
		} else {
			fmt.Printf("%d most recent)\n", keepCount)
		}
		fmt.Println()
		return 0, 0
	}

	// Delete backups
	var deletedCount int
	var freedSpace int64

	for _, backup := range toDelete {
		age := time.Since(backup.ModTime)
		ageStr := formatAge(age)

		if dryRun {
			fmt.Printf("   Would delete: %-35s %10s  %s old\n",
				backup.Name,
				formatSize(backup.Size),
				ageStr,
			)
		} else {
			fmt.Printf("   Deleting: %-35s %10s  %s old\n",
				backup.Name,
				formatSize(backup.Size),
				ageStr,
			)

			var err error
			if backup.IsDir {
				err = os.RemoveAll(backup.Path)
			} else {
				err = os.Remove(backup.Path)
			}

			if err != nil {
				fmt.Printf("      Error: %v\n", err)
				continue
			}
		}

		deletedCount++
		freedSpace += backup.Size
	}

	fmt.Printf("   Total: %d backup(s), %s\n\n", deletedCount, formatSize(freedSpace))

	return deletedCount, freedSpace
}

func formatAge(d time.Duration) string {
	days := int(d.Hours() / 24)
	if days == 0 {
		hours := int(d.Hours())
		if hours == 0 {
			return "< 1 hour"
		}
		return fmt.Sprintf("%d hour(s)", hours)
	}
	if days < 30 {
		return fmt.Sprintf("%d day(s)", days)
	}
	months := days / 30
	if months < 12 {
		return fmt.Sprintf("%d month(s)", months)
	}
	years := days / 365
	return fmt.Sprintf("%d year(s)", years)
}
