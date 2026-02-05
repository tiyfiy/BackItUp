package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/tiyfiy/BackItUp/internal/config"
	"github.com/tiyfiy/BackItUp/internal/mysql"
)

var mysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "skibidi",
	Long:  "long skbidi",

	Run: backupMySQL,
}

func init() {
	rootCmd.AddCommand(mysqlCmd)

	mysqlCmd.Flags().Bool("config", false, "Configure MySQL settings")
	mysqlCmd.Flags().String("host", "", "MySQL host")
	mysqlCmd.Flags().String("port", "", "MySQL port")
	mysqlCmd.Flags().String("user", "", "MySQL user")
	mysqlCmd.Flags().String("password", "", "MySQL password")
	mysqlCmd.Flags().String("database", "", "MySQL database")
}

func backupMySQL(cmd *cobra.Command, args []string) {
	configMode, _ := cmd.Flags().GetBool("config")
	host, _ := cmd.Flags().GetString("host")
	port, _ := cmd.Flags().GetString("port")
	user, _ := cmd.Flags().GetString("user")
	password, _ := cmd.Flags().GetString("password")
	database, _ := cmd.Flags().GetString("database")

	if configMode {
		if host != "" {
			config.SetMySQLHost(host)
			fmt.Printf("MySQL host saved to config\n")
			return
		} else if port != "" {
			config.SetMySQLPort(port)
			fmt.Printf("MySQL port saved to config\n")
			return
		} else if user != "" {
			config.SetMySQLUser(user)
			fmt.Printf("MySQL user saved to config\n")
			return
		} else if password != "" {
			config.SetMySQLPassword(password)
			fmt.Printf("MySQL password saved to config\n")
			return
		} else if database != "" {
			config.SetMySQLDatabase(database)
			fmt.Printf("MySQL database saved to config\n")
			return
		} else {
			log.Fatal("when using config you must provide a value")
		}
	}

	fmt.Println("Backing up mysql...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := mysql.Connection(cfg.MySQL.Host, cfg.MySQL.Port, cfg.MySQL.User, cfg.MySQL.Password, cfg.MySQL.Database)
	if err != nil {
		fmt.Println("error from the connection")
		log.Fatal(err)
	}
	defer db.Close()

	mysql.Backup(db, cfg.MySQL.Host, cfg.MySQL.Port, cfg.MySQL.User, cfg.MySQL.Password, cfg.MySQL.Database)
}
