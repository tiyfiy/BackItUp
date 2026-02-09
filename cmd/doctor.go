package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Analyze backup health and get recommendations",
	Long: `Doctor analyzes your backup patterns, detects anomalies, and provides
smart recommendations to optimize your backup strategy.

Features:
  - Backup trend visualization with ASCII charts
  - Growth rate analysis
  - Anomaly detection
  - Health score calculation
  - Disk space warnings
  - Personalized recommendations`,
	Run: func(cmd *cobra.Command, args []string) {
		runDoctorAnalysis()
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}

type BackupStats struct {
	TotalBackups    int
	TotalSize       int64
	OldestBackup    time.Time
	NewestBackup    time.Time
	AverageSize     int64
	GrowthRate      float64
	Anomalies       []string
	SizeHistory     []int64
	TimeHistory     []time.Time
	DiskUsagePercent float64
}

func runDoctorAnalysis() {
	fmt.Println("\n🏥 BackItUp Doctor - Backup Health Analysis")
	fmt.Println("═══════════════════════════════════════════════════════════════")

	backupDir := "BACKUP"
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		fmt.Println("\n❌ No backups found. Run some backups first!")
		return
	}

	// Analyze each database type
	mongoStats := analyzeBackups(filepath.Join(backupDir, "mongo"), "MongoDB")
	mysqlStats := analyzeBackups(filepath.Join(backupDir, "mysql"), "MySQL")
	pgStats := analyzeBackups(filepath.Join(backupDir, "postgresql"), "PostgreSQL")

	// Print individual database analyses
	if mongoStats.TotalBackups > 0 {
		printDatabaseAnalysis("MongoDB", mongoStats)
	}
	if mysqlStats.TotalBackups > 0 {
		printDatabaseAnalysis("MySQL", mysqlStats)
	}
	if pgStats.TotalBackups > 0 {
		printDatabaseAnalysis("PostgreSQL", pgStats)
	}

	// Calculate overall health score
	allStats := []*BackupStats{&mongoStats, &mysqlStats, &pgStats}
	healthScore := calculateHealthScore(allStats)

	// Print overall summary
	printOverallSummary(allStats, healthScore)

	// Print recommendations
	printRecommendations(allStats, healthScore)
}

func analyzeBackups(path, dbType string) BackupStats {
	stats := BackupStats{
		SizeHistory: make([]int64, 0),
		TimeHistory: make([]time.Time, 0),
		Anomalies:   make([]string, 0),
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return stats
	}

	backups := make([]BackupInfo, 0)

	// Collect all backup info
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		var size int64
		if entry.IsDir() {
			size = getDirSize(filepath.Join(path, entry.Name()))
		} else {
			size = info.Size()
		}

		backups = append(backups, BackupInfo{
			Name:    entry.Name(),
			Size:    size,
			ModTime: info.ModTime(),
		})
	}

	if len(backups) == 0 {
		return stats
	}

	// Sort by time
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].ModTime.Before(backups[j].ModTime)
	})

	// Calculate statistics
	stats.TotalBackups = len(backups)
	stats.OldestBackup = backups[0].ModTime
	stats.NewestBackup = backups[len(backups)-1].ModTime

	var totalSize int64
	for _, backup := range backups {
		totalSize += backup.Size
		stats.SizeHistory = append(stats.SizeHistory, backup.Size)
		stats.TimeHistory = append(stats.TimeHistory, backup.ModTime)
	}

	stats.TotalSize = totalSize
	stats.AverageSize = totalSize / int64(len(backups))

	// Calculate growth rate
	if len(backups) >= 2 {
		firstSize := float64(backups[0].Size)
		lastSize := float64(backups[len(backups)-1].Size)
		if firstSize > 0 {
			stats.GrowthRate = ((lastSize - firstSize) / firstSize) * 100
		}
	}

	// Detect anomalies
	stats.Anomalies = detectAnomalies(backups, stats.AverageSize)

	return stats
}

func detectAnomalies(backups []BackupInfo, avgSize int64) []string {
	anomalies := make([]string, 0)

	// Check for size anomalies (backups >2x or <0.5x average)
	for _, backup := range backups {
		if avgSize > 0 {
			ratio := float64(backup.Size) / float64(avgSize)
			if ratio > 2.0 {
				anomalies = append(anomalies, fmt.Sprintf("📈 %s is %.1fx larger than average",
					backup.Name, ratio))
			} else if ratio < 0.5 && backup.Size > 0 {
				anomalies = append(anomalies, fmt.Sprintf("📉 %s is %.1fx smaller than average",
					backup.Name, 1/ratio))
			}
		}
	}

	// Check for gaps in backup schedule
	if len(backups) >= 2 {
		for i := 1; i < len(backups); i++ {
			gap := backups[i].ModTime.Sub(backups[i-1].ModTime)
			if gap > 7*24*time.Hour {
				anomalies = append(anomalies, fmt.Sprintf("⏰ %d-day gap between backups (%s to %s)",
					int(gap.Hours()/24),
					backups[i-1].ModTime.Format("Jan 2"),
					backups[i].ModTime.Format("Jan 2")))
			}
		}
	}

	return anomalies
}

func printDatabaseAnalysis(dbName string, stats BackupStats) {
	fmt.Printf("\n\n📊 %s Analysis\n", dbName)
	fmt.Println("───────────────────────────────────────────────────────────────")

	fmt.Printf("  Total Backups:  %d\n", stats.TotalBackups)
	fmt.Printf("  Total Size:     %s\n", formatSize(stats.TotalSize))
	fmt.Printf("  Average Size:   %s\n", formatSize(stats.AverageSize))
	fmt.Printf("  Date Range:     %s → %s\n",
		stats.OldestBackup.Format("2006-01-02"),
		stats.NewestBackup.Format("2006-01-02"))

	if stats.GrowthRate != 0 {
		growthEmoji := "📈"
		if stats.GrowthRate < 0 {
			growthEmoji = "📉"
		}
		fmt.Printf("  Growth Rate:    %s %.1f%%\n", growthEmoji, stats.GrowthRate)
	}

	// Print size trend chart
	if len(stats.SizeHistory) > 1 {
		fmt.Println("\n  Size Trend:")
		printMiniChart(stats.SizeHistory)
	}

	// Print anomalies
	if len(stats.Anomalies) > 0 {
		fmt.Println("\n  ⚠️  Anomalies Detected:")
		for _, anomaly := range stats.Anomalies {
			fmt.Printf("     • %s\n", anomaly)
		}
	}
}

func printMiniChart(sizes []int64) {
	if len(sizes) == 0 {
		return
	}

	// Find min and max for normalization
	minSize := sizes[0]
	maxSize := sizes[0]
	for _, size := range sizes {
		if size < minSize {
			minSize = size
		}
		if size > maxSize {
			maxSize = size
		}
	}

	chartWidth := 50
	chartHeight := 8

	// Normalize and create chart
	chart := make([][]string, chartHeight)
	for i := range chart {
		chart[i] = make([]string, chartWidth)
		for j := range chart[i] {
			chart[i][j] = " "
		}
	}

	// Plot points
	for i, size := range sizes {
		x := (i * chartWidth) / len(sizes)
		if x >= chartWidth {
			x = chartWidth - 1
		}

		var y int
		if maxSize > minSize {
			normalized := float64(size-minSize) / float64(maxSize-minSize)
			y = chartHeight - 1 - int(normalized*float64(chartHeight-1))
		} else {
			y = chartHeight / 2
		}

		if y >= 0 && y < chartHeight && x >= 0 && x < chartWidth {
			chart[y][x] = "●"
		}

		// Draw line to next point
		if i < len(sizes)-1 {
			nextSize := sizes[i+1]
			nextX := ((i + 1) * chartWidth) / len(sizes)
			if nextX >= chartWidth {
				nextX = chartWidth - 1
			}

			var nextY int
			if maxSize > minSize {
				normalized := float64(nextSize-minSize) / float64(maxSize-minSize)
				nextY = chartHeight - 1 - int(normalized*float64(chartHeight-1))
			} else {
				nextY = chartHeight / 2
			}

			// Simple line drawing
			for dx := x; dx <= nextX && dx < chartWidth; dx++ {
				progress := float64(dx-x) / float64(nextX-x+1)
				lineY := y + int(progress*float64(nextY-y))
				if lineY >= 0 && lineY < chartHeight {
					if chart[lineY][dx] == " " {
						chart[lineY][dx] = "─"
					}
				}
			}
		}
	}

	// Print chart
	fmt.Printf("     %s\n", formatSize(maxSize))
	for i, row := range chart {
		if i == chartHeight-1 {
			fmt.Printf("     %s %s\n", strings.Join(row, ""), formatSize(minSize))
		} else {
			fmt.Printf("     %s\n", strings.Join(row, ""))
		}
	}
	fmt.Printf("     %-50s\n", "oldest → newest")
}

func calculateHealthScore(allStats []*BackupStats) int {
	score := 100

	for _, stats := range allStats {
		if stats.TotalBackups == 0 {
			continue
		}

		// Deduct points for no recent backups
		daysSinceLastBackup := time.Since(stats.NewestBackup).Hours() / 24
		if daysSinceLastBackup > 7 {
			score -= 15
		} else if daysSinceLastBackup > 3 {
			score -= 5
		}

		// Deduct points for anomalies
		score -= len(stats.Anomalies) * 3

		// Deduct points for too few backups
		if stats.TotalBackups < 3 {
			score -= 10
		}

		// Deduct points for rapid growth (potential issue)
		if stats.GrowthRate > 100 {
			score -= 10
		}
	}

	if score < 0 {
		score = 0
	}

	return score
}

func printOverallSummary(allStats []*BackupStats, healthScore int) {
	fmt.Println("\n\n💊 Overall Health Score")
	fmt.Println("═══════════════════════════════════════════════════════════════")

	// Health score bar
	scoreBar := strings.Repeat("█", healthScore/5) + strings.Repeat("░", (100-healthScore)/5)
	scoreEmoji := "🟢"
	scoreLabel := "Excellent"

	if healthScore < 70 {
		scoreEmoji = "🟡"
		scoreLabel = "Good"
	}
	if healthScore < 50 {
		scoreEmoji = "🟠"
		scoreLabel = "Fair"
	}
	if healthScore < 30 {
		scoreEmoji = "🔴"
		scoreLabel = "Needs Attention"
	}

	fmt.Printf("\n  %s  %s\n", scoreEmoji, scoreLabel)
	fmt.Printf("  %s %d/100\n\n", scoreBar, healthScore)

	// Total statistics
	var totalBackups int
	var totalSize int64
	for _, stats := range allStats {
		totalBackups += stats.TotalBackups
		totalSize += stats.TotalSize
	}

	fmt.Printf("  📦 Total Backups: %d\n", totalBackups)
	fmt.Printf("  💾 Total Storage: %s\n", totalSize)
}

func printRecommendations(allStats []*BackupStats, healthScore int) {
	fmt.Println("\n\n💡 Recommendations")
	fmt.Println("═══════════════════════════════════════════════════════════════")

	recommendations := make([]string, 0)

	for _, stats := range allStats {
		if stats.TotalBackups == 0 {
			continue
		}

		// Check backup frequency
		daysSinceLastBackup := time.Since(stats.NewestBackup).Hours() / 24
		if daysSinceLastBackup > 7 {
			recommendations = append(recommendations,
				"⏰ It's been more than 7 days since your last backup. Consider running 'backup-all'")
		}

		// Check for rapid growth
		if stats.GrowthRate > 50 {
			recommendations = append(recommendations,
				"📈 Your database is growing rapidly (%.1f%%). Consider more frequent backups")
		}

		// Check disk space
		if stats.TotalSize > 10*1024*1024*1024 { // > 10GB
			recommendations = append(recommendations,
				"🗄️  Large backup storage detected. Consider using 'cleanup' to remove old backups")
		}

		// Check backup count
		if stats.TotalBackups < 3 {
			recommendations = append(recommendations,
				"📊 You have fewer than 3 backups. Build a backup history for better protection")
		}

		// Check for anomalies
		if len(stats.Anomalies) > 0 {
			recommendations = append(recommendations,
				"⚠️  Anomalies detected in backup sizes. Review your backup patterns")
		}
	}

	// Schedule recommendation
	hasSchedule := false // We'd check crontab in real implementation
	if !hasSchedule {
		recommendations = append(recommendations,
			"🤖 Set up automated backups with 'schedule' command for peace of mind")
	}

	if len(recommendations) == 0 {
		fmt.Println("\n  ✅ Everything looks great! Your backups are healthy.")
		fmt.Println("  ✅ Keep up the good work!")
	} else {
		for i, rec := range recommendations {
			fmt.Printf("\n  %d. %s\n", i+1, rec)
		}
	}

	fmt.Println("\n")
}
