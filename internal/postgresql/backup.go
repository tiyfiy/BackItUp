package postgresql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func Backup(db *sql.DB, host, port, user, password, database string) {
	path := "BACKUP/postgresql"

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

	fmt.Printf("✅ Backup completed: %s\n", outfile)
}
