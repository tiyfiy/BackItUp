package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var PostgreSQLCmd = &cobra.Command{
	Use:   "postgresql",
	Short: "command for using postreSQL database",
	Long:  "this is the command for using PostreSQL database",
	Run:   backupPostgreSQL,
}

func init() {
	rootCmd.AddCommand(PostgreSQLCmd)
}

func backupPostgreSQL(cmd *cobra.Command, args []string) {
	fmt.Println("this is POstresql command")
}
