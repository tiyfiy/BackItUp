package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func Backup(db *sql.DB, host, port, user, password, database string) {
	path := "BACKUP/mysql"

	err := os.MkdirAll(path, 0755)
	if err != nil {
		log.Fatal(err)
	}

	outfile := fmt.Sprintf("%s/%s.sql", path, database)

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

	fmt.Println("backup completed")
}
