package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tiyfiy/BackItUp/internal/mongodb"
)

var mongodbCmd = &cobra.Command{
	Use:   "mongodb",
	Short: "skibidi",
	Long:  "long skbidi",

	Run: backupMongodb,
}

func init() {
	rootCmd.AddCommand(mongodbCmd)
}

func backupMongodb(cmd *cobra.Command, args []string) {
	fmt.Println("mongodb...")
	mongodb.ConnectionMongodb()
}
