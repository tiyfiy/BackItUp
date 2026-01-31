package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var MySQLCmd = &cobra.Command{
	Use:   "random",
	Short: "short description",
	Long:  "Long very long desctiption of ranodm command idk",

	Run: generateAnswer,
}

func init() {
	rootCmd.AddCommand(MySQLCmd)

	// MySQLCmd.Flags()
}

func generateAnswer(cmd *cobra.Command, args []string) {
	fmt.Println("skbidi toilet")
}
