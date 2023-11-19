package main

import (
	"context"
	"fmt"
	"time"

	"github.com/marcosrosse/bucket-migration-tool/internal/database"
)

type Avatar struct {
	Id   int
	Path string
}

// Maybe here will have all the logic
func worker(workerId int, jobs chan Avatar) {
	for j := range jobs {
		fmt.Println("Worker: ", workerId, "jobs: ", j)
		time.Sleep(time.Second)
	}

}

func main() {

	db := database.ConnDB()

	var counter int
	db.QueryRow(context.Background(), "select count(*) from avatars").Scan(&counter)
	defer db.Close(context.Background())

	limit := 100

	// create the channel jobs with the avatar type
	jobs := make(chan Avatar, counter)

	// Start a go routine sending an id and a jobs to the worker function

	for w := 1; w <= 10; w++ {
		go worker(w, jobs)
	}

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
			jobs <- avatar // Send each line for the message channel
		}

		counter -= limit

		time.Sleep(3 * time.Second)

	}
}
