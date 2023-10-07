package s3

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Set configuration to access minIO or S3
func SetS3() (*minio.Client, error) {
	endpoint := os.Getenv("S3_ENDPOINT")
	accessKeyID := os.Getenv("S3_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("S3_SECRET_ACCESS_KEY")
	s3Ssl := os.Getenv("S3_USE_SSL")
	useSsl, _ := strconv.ParseBool(s3Ssl)

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSsl,
	})

	if err != nil {
		log.Fatalln(err)
	}

	return minioClient, err
}

// Check if Object already exists in the bucket
func ObjExists(bucket, object string) (bool, error) {
	s3, err := SetS3()
	if err != nil {
		return false, err
	}
	_, err = s3.StatObject(context.Background(), bucket, object, minio.StatObjectOptions{})
	if err != nil {
		return false, err
	}
	return true, err
}

// Copy objects from one bucket to another in the S3 or minIO
func CopyObjs(srcBucket, dstBucket, objSrcName, objDstName string) error {
	s3, err := SetS3()
	if err != nil {
		return err
	}
	// Preparing source object
	srcOpts := minio.CopySrcOptions{
		Bucket: srcBucket,
		Object: objSrcName,
	}

	// Preparing destination object
	dstOpts := minio.CopyDestOptions{
		Bucket: dstBucket,
		Object: objDstName,
	}

	// Copy object from source to destination bucket
	_, err = s3.CopyObject(context.Background(), dstOpts, srcOpts)
	if err != nil {
		return err
	}
	return nil

}
