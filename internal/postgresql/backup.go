package postgresql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func Backup(db *sql.DB, host, port, user, password, database string) {
	path := "BACKUP/postgresql"

	err := os.MkdirAll(path, 0755)
	if err != nil {
		log.Fatal(err)
	}

	outfile := fmt.Sprintf("%s/%s.sql", path, database)

	cmd := exec.Command("pg_dump",
		"-h", host,
		"-p", port,
		"-U", user,
		"-d", database,
		"-f", outfile,
	)

	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", password))

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("backup completed")
}
