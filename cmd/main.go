package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/marcosrosse/bucket-migration-tool/internal/database"
)

type Avatar struct {
	Id   int
	Path string
}

// Maybe here will have all the logic
func worker(workerId int, msg chan Avatar) {
	for res := range msg {
		fmt.Println("Worker: ", workerId, " Msg: ", res)
		time.Sleep(time.Second)
	}

}

func main() {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	username := os.Getenv("POSTGRES_USERNAME")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DBNAME")

	psqlConn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		username, password, host, port, dbname)
	db := database.ConnDB(psqlConn, 10)

	var counter int
	// Comented cause there is to much rows in the table
	// db.QueryRow(context.Background(), "select count(*) from avatars").Scan(&counter)
	// fmt.Println("This is the total of rows", counter)

	counter = 40
	limit := 10

	// create the channel msg with the avatar type
	msg := make(chan Avatar)

	// Start a go routine sending an id and a msg to the worker function
	go worker(1, msg)
	go worker(2, msg)

	for counter > 0 {
		page := counter / limit
		offset := limit * (page - 1)

		SQL := `SELECT "id","path" FROM "avatars" ORDER BY "id" LIMIT $1 OFFSET $2`

		rows, _ := db.Query(context.Background(), SQL, limit, offset)
		defer rows.Close()

		var id int
		var path string
		for rows.Next() {
			rows.Scan(&id, &path)
			// Populate the struct
			avatar := Avatar{
				Id:   id,
				Path: path,
			}
			msg <- avatar // Send each line for the message channel
		}

		counter -= limit

		time.Sleep(3 * time.Second)

	}
}
