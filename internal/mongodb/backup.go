package mongodb

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Backup(client *mongo.Client, uri string) {
	// Add timestamp to directory name
	now := time.Now()
	timestamp := fmt.Sprintf("%d-%02d-%02d_%02d-%02d-%02d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second())
	path := fmt.Sprintf("BACKUP/mongo/backup_%s", timestamp)

	cmd := exec.Command("mongodump", "--uri", uri, "--out", path)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("✅ Backup completed: %s\n", path)
}
