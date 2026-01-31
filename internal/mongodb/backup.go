package mongodb

import (
	"fmt"
	"log"
	"os/exec"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Backup(client *mongo.Client, uri string) {
	// database := client.Database("sample_mfix")

	path := "BACKUP/mongo"

	cmd := exec.Command("mongodump", "--uri", uri, "--out", path)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("backup completed")
}
