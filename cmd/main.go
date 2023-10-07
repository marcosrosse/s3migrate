package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/marcosrosse/copys3files/internal/database"
	"github.com/marcosrosse/copys3files/internal/s3"
)

var (
    srcBucket *string
    dstBucket *string
    srcPath *string
    dstPath *string
)


func init (){
    //TODO: Implement a verification if the flag is realy parsed
	srcBucket = flag.String("src-bucket", "legacy-s3", "The source bucket where the legacy files are stored.")
	dstBucket = flag.String("dst-bucket", "production-s3", "The destination bucket where the files will be placed.")
	srcPath = flag.String("src-path", "/image", "Source files path.")
	dstPath = flag.String("dst-path", "/avatar", "Destination files path.")
}

func main() {

    // Flag parse
    flag.Parse()

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
		objName := strings.TrimPrefix(path, "image/")
		legacyObj[id] = objName
	}

	// Range all key and values base in the id and path from DB
	for key, value := range legacyObj {
		// Concatenate path with the image name
		objSrcName := (*srcPath + value)
		objDstName := (*dstPath + value)

		// If object didn't exist, copy it to the bucket
		if obj, _ := s3.ObjExists(*dstBucket, objDstName); !obj {

			err = s3.CopyObjs(*srcBucket, *dstBucket, objSrcName, objDstName)

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
