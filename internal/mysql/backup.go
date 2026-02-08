package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func Backup(db *sql.DB, host, port, user, password, database string) {
	path := "BACKUP/mysql"

	err := os.MkdirAll(path, 0755)
	if err != nil {
		log.Fatal(err)
	}

	// Add timestamp to filename to prevent overwrites
	now := time.Now()
	timestamp := fmt.Sprintf("%d-%02d-%02d_%02d-%02d-%02d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second())
	outfile := fmt.Sprintf("%s/%s_%s.sql", path, database, timestamp)

	cmd := exec.Command("mysqldump",
		"-h", host,
		"-P", port,
		"-u", user,
		fmt.Sprintf("-p%s", password),
		database,
	)

	output, err := os.Create(outfile)
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()

	cmd.Stdout = output
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("✅ Backup completed: %s\n", outfile)
}
