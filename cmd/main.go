package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/marcosrosse/s3migrate/internal/database"
	"github.com/marcosrosse/s3migrate/internal/s3"
)

var (
	srcBucket = os.Getenv("S3_SRC_BUCKET")
	dstBucket = os.Getenv("S3_DST_BUCKET")
	srcPath   = os.Getenv("S3_SRC_PATH_OBJ")
	dstPath   = os.Getenv("S3_DST_PATH_OBJ")
)

func main() {
	fmt.Println("Starting Job")

	// To the connectin with the DB
	db, err := database.ConnDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	// Query all the values with specified path
	rows, err := db.Query("SELECT id, path from avatars WHERE path LIKE 'image/%'")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	// Read the rows of the select query and remove the path image
	var id int
	var path string
	var legacyObj = make(map[int]string)
	for rows.Next() {
		err = rows.Scan(&id, &path)
		if err != nil {
			log.Fatal(err)
		}
		objName := strings.TrimPrefix(path, srcPath)
		legacyObj[id] = objName
	}

	// Range all key and values base in the id and path from DB
	for key, value := range legacyObj {
		// Concatenate path with the image name
		objSrcName := (srcPath + value)
		objDstName := (dstPath + value)

		// If object didn't exist, copy it to the bucket
		if obj, _ := s3.ObjExists(dstBucket, objDstName); !obj {

			err = s3.CopyObjs(srcBucket, dstBucket, objSrcName, objDstName)

			if err != nil {
				fmt.Println("Failed to copy the object", err)
			} else {
				// Update each success upload in the Postgres with the new path
				// TODO: Implement a bulk update instead line by line
				sqlStatement := `UPDATE avatars SET path = $1 WHERE id = $2`
				_, err = db.Exec(sqlStatement, objDstName, key)
				if err != nil {
					log.Panic("Error to insert in the DB", err)

				}
			}

		} else {
			log.Printf("Object: %#value already exists\n", objDstName)
		}
	}

}
