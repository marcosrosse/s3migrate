package database

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
)

// Connect to the Postgres Database
func ConnDB(connectionUrl string, retries int) *pgx.Conn {

	var db *pgx.Conn
	var err error

	db, err = pgx.Connect(context.Background(), connectionUrl)
	for err != nil {
		log.Printf("%s\n", err)
		if retries > 1 {
			retries--
			time.Sleep(5 * time.Second)
			db, err = pgx.Connect(context.Background(), connectionUrl)
			continue
		}
		panic(err)
	}

	return db
}
