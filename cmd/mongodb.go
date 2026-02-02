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

	mongodbCmd.Flags().Bool("config", false, "Configure MongoDB settings")
	mongodbCmd.Flags().String("uri", "", "MongoDB connection URI")
	mongodbCmd.Flags().String("path", "", "Path where the backups should be saved")
}

func backupMongodb(cmd *cobra.Command, args []string) {
	configMode, _ := cmd.Flags().GetBool("config")
	uri, _ := cmd.Flags().GetString("uri")
	path, _ := cmd.Flags().GetString("path")

	if configMode {
		if uri != "" {
			config.SetMongodbURI(uri)

			fmt.Printf("MongoDB URI saved to config\n")
			return
		} else if path != "" {
			config.SetMongodbPath(path)

			fmt.Printf("MongoDB backup path saved to config\n")
			return
		} else {
			log.Fatal("when using config us must provide URI")
		}
	}

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
