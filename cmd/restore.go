package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tiyfiy/BackItUp/internal/config"
)

var (
	restoreLatest bool
	restoreFile   string
)

var restoreCmd = &cobra.Command{
	Use:   "restore [mongodb|mysql|postgresql]",
	Short: "Restore a database from backup",
	Long: `Restore a database from a previously created backup.

By default, shows available backups and prompts for selection.
Use --latest to automatically restore the most recent backup.
Use --file to specify a specific backup file/directory.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dbType := args[0]
		restoreDatabase(dbType)
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
	restoreCmd.Flags().BoolVarP(&restoreLatest, "latest", "l", false, "Restore the latest backup")
	restoreCmd.Flags().StringVarP(&restoreFile, "file", "f", "", "Restore from specific backup file/directory")
}

func restoreDatabase(dbType string) {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}

	var backupPath string
	var backups []BackupInfo

	switch dbType {
	case "mongodb":
		backups = listDirBackups(filepath.Join("BACKUP", "mongo"))
		if len(backups) == 0 {
			fmt.Println("No MongoDB backups found.")
			return
		}

		if restoreFile != "" {
			backupPath = restoreFile
		} else if restoreLatest {
			backupPath = backups[0].Path
		} else {
			backupPath = selectBackup(backups, "MongoDB")
		}

		if backupPath == "" {
			fmt.Println("No backup selected. Restore cancelled.")
			return
		}

		if !confirmRestore("MongoDB") {
			fmt.Println("Restore cancelled.")
			return
		}

		restoreMongoDB(cfg.MongoDB.URI, backupPath)

	case "mysql":
		backups = listFileBackups(filepath.Join("BACKUP", "mysql"), ".sql")
		if len(backups) == 0 {
			fmt.Println("No MySQL backups found.")
			return
		}

		if restoreFile != "" {
			backupPath = restoreFile
		} else if restoreLatest {
			backupPath = backups[0].Path
		} else {
			backupPath = selectBackup(backups, "MySQL")
		}

		if backupPath == "" {
			fmt.Println("No backup selected. Restore cancelled.")
			return
		}

		if !confirmRestore("MySQL") {
			fmt.Println("Restore cancelled.")
			return
		}

		restoreMySQL(cfg.MySQL, backupPath)

	case "postgresql":
		backups = listFileBackups(filepath.Join("BACKUP", "postgresql"), ".sql")
		if len(backups) == 0 {
			fmt.Println("No PostgreSQL backups found.")
			return
		}

		if restoreFile != "" {
			backupPath = restoreFile
		} else if restoreLatest {
			backupPath = backups[0].Path
		} else {
			backupPath = selectBackup(backups, "PostgreSQL")
		}

		if backupPath == "" {
			fmt.Println("No backup selected. Restore cancelled.")
			return
		}

		if !confirmRestore("PostgreSQL") {
			fmt.Println("Restore cancelled.")
			return
		}

		restorePostgreSQL(cfg.PostgreSQL, backupPath)

	default:
		fmt.Printf("Unknown database type: %s\n", dbType)
		fmt.Println("Supported types: mongodb, mysql, postgresql")
		return
	}
}

func selectBackup(backups []BackupInfo, dbType string) string {
	fmt.Printf("\n📦 Available %s Backups:\n", dbType)
	fmt.Println("══════════════════════════════════════════════════════════════")

	for i, backup := range backups {
		fmt.Printf("  [%d] %-40s %10s  %s\n",
			i+1,
			backup.Name,
			formatSize(backup.Size),
			backup.ModTime.Format("2006-01-02 15:04:05"),
		)
	}

	fmt.Println()
	fmt.Print("Select backup number to restore (or 0 to cancel): ")

	var selection int
	_, err := fmt.Scanln(&selection)
	if err != nil || selection < 1 || selection > len(backups) {
		return ""
	}

	return backups[selection-1].Path
}

func confirmRestore(dbType string) bool {
	fmt.Println()
	fmt.Printf("⚠️  WARNING: This will restore the %s database.\n", dbType)
	fmt.Println("   This operation may overwrite existing data!")
	fmt.Println()
	fmt.Print("Are you sure you want to continue? (yes/no): ")

	var response string
	fmt.Scanln(&response)

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "yes" || response == "y"
}

func restoreMongoDB(uri, backupPath string) {
	fmt.Println("\n🔄 Restoring MongoDB from backup...")
	fmt.Printf("   Source: %s\n", backupPath)
	fmt.Println()

	// Check if mongorestore is available
	if _, err := exec.LookPath("mongorestore"); err != nil {
		log.Fatal("mongorestore command not found. Please install MongoDB tools.")
	}

	cmd := exec.Command("mongorestore", "--uri", uri, "--drop", backupPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatal("Restore failed:", err)
	}

	fmt.Println("\n✅ MongoDB restore completed successfully!")
}

func restoreMySQL(mysqlCfg config.MySQLConfig, backupPath string) {
	fmt.Println("\n🔄 Restoring MySQL from backup...")
	fmt.Printf("   Source: %s\n", backupPath)
	fmt.Printf("   Database: %s\n", mysqlCfg.Database)
	fmt.Println()

	if mysqlCfg.Database == "" {
		log.Fatal("MySQL database not configured. Run: ./BackItUp mysql --config --database yourdb")
	}

	// Check if mysql is available
	if _, err := exec.LookPath("mysql"); err != nil {
		log.Fatal("mysql command not found. Please install MySQL client.")
	}

	// Read the backup file
	backupData, err := os.ReadFile(backupPath)
	if err != nil {
		log.Fatal("Failed to read backup file:", err)
	}

	// Execute restore
	cmd := exec.Command("mysql",
		"-h", mysqlCfg.Host,
		"-P", mysqlCfg.Port,
		"-u", mysqlCfg.User,
		fmt.Sprintf("-p%s", mysqlCfg.Password),
		mysqlCfg.Database,
	)

	cmd.Stdin = strings.NewReader(string(backupData))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Fatal("Restore failed:", err)
	}

	fmt.Println("\n✅ MySQL restore completed successfully!")
}

func restorePostgreSQL(pgCfg config.PostgreSQLConfig, backupPath string) {
	fmt.Println("\n🔄 Restoring PostgreSQL from backup...")
	fmt.Printf("   Source: %s\n", backupPath)
	fmt.Printf("   Database: %s\n", pgCfg.Database)
	fmt.Println()

	if pgCfg.Database == "" {
		log.Fatal("PostgreSQL database not configured. Run: ./BackItUp postgresql --config --database yourdb")
	}

	// Check if psql is available
	if _, err := exec.LookPath("psql"); err != nil {
		log.Fatal("psql command not found. Please install PostgreSQL client.")
	}

	cmd := exec.Command("psql",
		"-h", pgCfg.Host,
		"-p", pgCfg.Port,
		"-U", pgCfg.User,
		"-d", pgCfg.Database,
		"-f", backupPath,
	)

	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", pgCfg.Password))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatal("Restore failed:", err)
	}

	fmt.Println("\n✅ PostgreSQL restore completed successfully!")
}
