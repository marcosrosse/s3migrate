package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
)

var (
	host     = os.Getenv("POSTGRES_HOST")
	port     = os.Getenv("POSTGRES_PORT")
	username = os.Getenv("POSTGRES_USERNAME")
	password = os.Getenv("POSTGRES_PASSWORD")
	dbname   = os.Getenv("POSTGRES_DBNAME")
	DB       *pgx.Conn
)

// Connect to the Postgres Database
func ConnDB() *pgx.Conn {
	retries := 10
	psqlConn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		username, password, host, port, dbname)

	var db *pgx.Conn
	var err error

	db, err = pgx.Connect(context.Background(), psqlConn)
	for err != nil {
		log.Printf("%s\n", err)
		if retries > 1 {
			retries--
			time.Sleep(5 * time.Second)
			db, err = pgx.Connect(context.Background(), psqlConn)
			continue
		}
		panic(err)
	}
	DB = db
	return db
}
