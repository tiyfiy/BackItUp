package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "short description",
	Long:  "Long very long desctiption of ranodm command idk",

	Run: generateAnswer,
}

func generateAnswer(cmd *cobra.Command, args []string) {
	fmt.Println("skbidi toilet")
}
