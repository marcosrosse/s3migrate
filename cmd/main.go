package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/marcosrosse/s3migrate/internal/database"
)

type Avatar struct {
	Id   int
	Path string
}

var (
	client                *s3.S3
	AWS_ENDPOINT_URL      = os.Getenv("AWS_ENDPOINT_URL")
	AWS_ACCESS_KEY_ID     = os.Getenv("AWS_ACCESS_KEY_ID")
	AWS_SECRET_ACCESS_KEY = os.Getenv("AWS_SECRET_ACCESS_KEY")
	AWS_REGION            = os.Getenv("AWS_REGION")
)

func init() {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(
			AWS_ACCESS_KEY_ID,
			AWS_SECRET_ACCESS_KEY, ""),
		Region:   aws.String(AWS_REGION),
		Endpoint: aws.String(AWS_ENDPOINT_URL),
	})
	if err != nil {
		panic(err)
	}
	client = s3.New(sess)
}

func worker(workerId int, job chan Avatar) {
	for j := range job {
		fmt.Println("Worker: ", workerId, "job: ", j)

		time.Sleep(time.Second)
	}

}

func ListBuckets(client *s3.S3) (*s3.ListBucketsOutput, error) {
	res, err := client.ListBuckets(nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func ListObjects(client *s3.S3, bucket, object string) (*s3.ListObjectsV2Output, error) {
	res, err := client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(object),
	})
	if err != nil {
		fmt.Println(err)
	}
	return res, nil

}

func main() {

	buckets, err := ListBuckets(client)
	if err != nil {
		fmt.Printf("Couldn't list buckets: %v", err)
		return
	}

	listObjects, _ := ListObjects(client, "legacy-s3", "avatar")
	fmt.Println(listObjects)

	for _, bucket := range buckets.Buckets {
		fmt.Printf("Found bucket: %s, created at: %s\n", *bucket.Name, *bucket.CreationDate)
	}
	db := database.ConnDB()

	var counter int
	db.QueryRow(context.Background(), "select count(*) from avatars").Scan(&counter)
	defer db.Close(context.Background())

	limit := 100

	// create the channel job with the avatar type
	job := make(chan Avatar, counter)

	// Start a go routine sending an id and a job to the worker function

	for w := 1; w <= 10; w++ {
		go worker(w, job)
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
			job <- avatar // Send each line for the message channel
		}

		counter -= limit

		time.Sleep(3 * time.Second)

	}
}
