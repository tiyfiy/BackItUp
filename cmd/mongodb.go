package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/tiyfiy/BackItUp/internal/config"
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
	fmt.Println("Backing up mongodb...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	client, err := mongodb.Connection(cfg.MongoDB.URI)
	if err != nil {
		fmt.Println("error from the connection")
	}

	mongodb.Backup(client, cfg.MongoDB.URI)
}
