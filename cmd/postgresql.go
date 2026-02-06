package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/tiyfiy/BackItUp/internal/config"
	"github.com/tiyfiy/BackItUp/internal/postgresql"
)

var postgresqlCmd = &cobra.Command{
	Use:   "postgresql",
	Short: "skibidi",
	Long:  "long skbidi",

	Run: backupPostgreSQL,
}

func init() {
	rootCmd.AddCommand(postgresqlCmd)

	postgresqlCmd.Flags().Bool("config", false, "Configure PostgreSQL settings")
	postgresqlCmd.Flags().String("host", "", "PostgreSQL host")
	postgresqlCmd.Flags().String("port", "", "PostgreSQL port")
	postgresqlCmd.Flags().String("user", "", "PostgreSQL user")
	postgresqlCmd.Flags().String("password", "", "PostgreSQL password")
	postgresqlCmd.Flags().String("database", "", "PostgreSQL database")
}

func backupPostgreSQL(cmd *cobra.Command, args []string) {
	configMode, _ := cmd.Flags().GetBool("config")
	host, _ := cmd.Flags().GetString("host")
	port, _ := cmd.Flags().GetString("port")
	user, _ := cmd.Flags().GetString("user")
	password, _ := cmd.Flags().GetString("password")
	database, _ := cmd.Flags().GetString("database")

	if configMode {
		if host != "" {
			config.SetPostgreSQLHost(host)
			fmt.Printf("PostgreSQL host saved to config\n")
			return
		} else if port != "" {
			config.SetPostgreSQLPort(port)
			fmt.Printf("PostgreSQL port saved to config\n")
			return
		} else if user != "" {
			config.SetPostgreSQLUser(user)
			fmt.Printf("PostgreSQL user saved to config\n")
			return
		} else if password != "" {
			config.SetPostgreSQLPassword(password)
			fmt.Printf("PostgreSQL password saved to config\n")
			return
		} else if database != "" {
			config.SetPostgreSQLDatabase(database)
			fmt.Printf("PostgreSQL database saved to config\n")
			return
		} else {
			log.Fatal("when using config you must provide a value")
		}
	}

	fmt.Println("Backing up postgresql...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := postgresql.Connection(cfg.PostgreSQL.Host, cfg.PostgreSQL.Port, cfg.PostgreSQL.User, cfg.PostgreSQL.Password, cfg.PostgreSQL.Database)
	if err != nil {
		fmt.Println("error from the connection")
		log.Fatal(err)
	}
	defer db.Close()

	postgresql.Backup(db, cfg.PostgreSQL.Host, cfg.PostgreSQL.Port, cfg.PostgreSQL.User, cfg.PostgreSQL.Password, cfg.PostgreSQL.Database)
}
