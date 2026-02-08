package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/spf13/cobra"
)

type BackupInfo struct {
	Name     string
	Path     string
	Size     int64
	ModTime  time.Time
	IsDir    bool
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available backups",
	Long:  `Display all available backups organized by database type, including size and modification time.`,
	Run: func(cmd *cobra.Command, args []string) {
		listBackups()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func listBackups() {
	backupDir := "BACKUP"

	// Check if BACKUP directory exists
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		fmt.Println("No backups found. The BACKUP directory doesn't exist yet.")
		fmt.Println("Run a backup command first to create backups.")
		return
	}

	hasBackups := false

	// List MongoDB backups
	mongoBackups := listDirBackups(filepath.Join(backupDir, "mongo"))
	if len(mongoBackups) > 0 {
		hasBackups = true
		fmt.Println("\n📦 MongoDB Backups:")
		fmt.Println("══════════════════════════════════════════════════════════════")
		printBackupList(mongoBackups)
	}

	// List MySQL backups
	mysqlBackups := listFileBackups(filepath.Join(backupDir, "mysql"), ".sql")
	if len(mysqlBackups) > 0 {
		hasBackups = true
		fmt.Println("\n📦 MySQL Backups:")
		fmt.Println("══════════════════════════════════════════════════════════════")
		printBackupList(mysqlBackups)
	}

	// List PostgreSQL backups
	pgBackups := listFileBackups(filepath.Join(backupDir, "postgresql"), ".sql")
	if len(pgBackups) > 0 {
		hasBackups = true
		fmt.Println("\n📦 PostgreSQL Backups:")
		fmt.Println("══════════════════════════════════════════════════════════════")
		printBackupList(pgBackups)
	}

	if !hasBackups {
		fmt.Println("\nNo backups found.")
		fmt.Println("Run a backup command to create your first backup.")
	} else {
		fmt.Println()
	}
}

func listDirBackups(path string) []BackupInfo {
	var backups []BackupInfo

	entries, err := os.ReadDir(path)
	if err != nil {
		return backups
	}

	for _, entry := range entries {
		if entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				continue
			}

			size := getDirSize(filepath.Join(path, entry.Name()))
			backups = append(backups, BackupInfo{
				Name:    entry.Name(),
				Path:    filepath.Join(path, entry.Name()),
				Size:    size,
				ModTime: info.ModTime(),
				IsDir:   true,
			})
		}
	}

	// Sort by modification time (newest first)
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].ModTime.After(backups[j].ModTime)
	})

	return backups
}

func listFileBackups(path, extension string) []BackupInfo {
	var backups []BackupInfo

	entries, err := os.ReadDir(path)
	if err != nil {
		return backups
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == extension {
			info, err := entry.Info()
			if err != nil {
				continue
			}

			backups = append(backups, BackupInfo{
				Name:    entry.Name(),
				Path:    filepath.Join(path, entry.Name()),
				Size:    info.Size(),
				ModTime: info.ModTime(),
				IsDir:   false,
			})
		}
	}

	// Sort by modification time (newest first)
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].ModTime.After(backups[j].ModTime)
	})

	return backups
}

func getDirSize(path string) int64 {
	var size int64
	filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size
}

func printBackupList(backups []BackupInfo) {
	for _, backup := range backups {
		typeStr := "file"
		if backup.IsDir {
			typeStr = "dir "
		}

		fmt.Printf("  [%s] %-40s %10s  %s\n",
			typeStr,
			backup.Name,
			formatSize(backup.Size),
			backup.ModTime.Format("2006-01-02 15:04:05"),
		)
	}
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
