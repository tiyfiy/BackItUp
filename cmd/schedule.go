package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Generate cron schedule examples for automated backups",
	Long: `Display cron job examples to help you schedule automated backups.

Shows ready-to-use cron expressions for common backup schedules like:
- Daily backups
- Weekly backups
- Monthly backups
- Custom intervals

You can copy these directly into your crontab.`,
	Run: func(cmd *cobra.Command, args []string) {
		showScheduleExamples()
	},
}

func init() {
	rootCmd.AddCommand(scheduleCmd)
}

func showScheduleExamples() {
	// Get absolute path to the binary
	execPath, err := os.Executable()
	if err != nil {
		execPath = "./BackItUp"
	} else {
		execPath, _ = filepath.Abs(execPath)
	}

	workingDir, err := os.Getwd()
	if err != nil {
		workingDir = "/path/to/backitup"
	}

	fmt.Println("⏰ Automated Backup Scheduling Guide")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	fmt.Println("📋 Common Cron Schedules")
	fmt.Println("───────────────────────────────────────────────────────────────")
	fmt.Println()

	fmt.Println("1️⃣  Daily at 2:00 AM (all databases):")
	fmt.Printf("   0 2 * * * cd %s && %s backup-all\n", workingDir, execPath)
	fmt.Println()

	fmt.Println("2️⃣  Every 6 hours:")
	fmt.Printf("   0 */6 * * * cd %s && %s backup-all\n", workingDir, execPath)
	fmt.Println()

	fmt.Println("3️⃣  Daily at 3:00 AM (MySQL only):")
	fmt.Printf("   0 3 * * * cd %s && %s mysql\n", workingDir, execPath)
	fmt.Println()

	fmt.Println("4️⃣  Weekly on Sunday at 1:00 AM:")
	fmt.Printf("   0 1 * * 0 cd %s && %s backup-all\n", workingDir, execPath)
	fmt.Println()

	fmt.Println("5️⃣  Monthly on 1st day at 2:00 AM:")
	fmt.Printf("   0 2 1 * * cd %s && %s backup-all\n", workingDir, execPath)
	fmt.Println()

	fmt.Println("6️⃣  Weekdays at 11:00 PM:")
	fmt.Printf("   0 23 * * 1-5 cd %s && %s backup-all\n", workingDir, execPath)
	fmt.Println()

	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("📚 Cron Format Reference:")
	fmt.Println("   ┌─────── minute (0-59)")
	fmt.Println("   │ ┌────── hour (0-23)")
	fmt.Println("   │ │ ┌───── day of month (1-31)")
	fmt.Println("   │ │ │ ┌──── month (1-12)")
	fmt.Println("   │ │ │ │ ┌─── day of week (0-7, Sun=0 or 7)")
	fmt.Println("   │ │ │ │ │")
	fmt.Println("   * * * * *")
	fmt.Println()

	fmt.Println("💡 How to Set Up:")
	fmt.Println("───────────────────────────────────────────────────────────────")
	fmt.Println("1. Edit your crontab:")
	fmt.Println("   crontab -e")
	fmt.Println()
	fmt.Println("2. Add one of the examples above")
	fmt.Println()
	fmt.Println("3. Save and exit")
	fmt.Println()
	fmt.Println("4. Verify with:")
	fmt.Println("   crontab -l")
	fmt.Println()

	fmt.Println("🔧 Pro Tips:")
	fmt.Println("───────────────────────────────────────────────────────────────")
	fmt.Println("• Combine with cleanup for space management:")
	fmt.Printf("  0 2 * * * cd %s && %s backup-all && %s cleanup all --days 30\n", workingDir, execPath, execPath)
	fmt.Println()
	fmt.Println("• Redirect output to log file:")
	fmt.Printf("  0 2 * * * cd %s && %s backup-all >> backup.log 2>&1\n", workingDir, execPath)
	fmt.Println()
	fmt.Println("• Test your cron job first by running the command manually!")
	fmt.Println()
}
