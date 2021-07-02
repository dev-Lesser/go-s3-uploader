package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joho/godotenv"
)

var AWS_SERVER_PUBLIC_KEY string
var AWS_SERVER_SECRET_KEY string

const (
	S3_REGION = "" // bucket region
	S3_BUCKET = "" // bucket name
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
		os.Exit(1)
	}
	AWS_SERVER_PUBLIC_KEY = os.Getenv("AWS_PUBLIC_KEY") // s3 upload pub key
	AWS_SERVER_SECRET_KEY = os.Getenv("AWS_SECRET_KEY") // s3 upload secret key
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(S3_REGION),
			Credentials: credentials.NewStaticCredentials(
				AWS_SERVER_PUBLIC_KEY,
				AWS_SERVER_SECRET_KEY,
				"",
			),
		})
	if err != nil {
		panic(err)
	}

	targetDir := "YOUR TARGET DIR NAME" // read files in target dir
	files, err := ioutil.ReadDir(targetDir)
	if err != nil {
		exitErrorf("Unable to read data %q", targetDir)
	}
	uploader := s3manager.NewUploader(sess) // create uploader
	for _, file := range files {
		fmt.Println(file.Name(), S3_BUCKET)
		go func(uploader *s3manager.Uploader, filename string) {
			file, err := os.Open(targetDir + filename)
			if err != nil { // file err check
				exitErrorf("Unable to open file %q, %v", err)
			}
			defer file.Close()

			_, err = uploader.Upload(&s3manager.UploadInput{ // s3 upload
				Bucket: aws.String(S3_BUCKET),
				Key:    aws.String(filename), // create (folder&file or file) path in s3 bucket
				Body:   file,
			})
			if err != nil {
				exitErrorf("Unable to upload %q to %q, %v", filename, S3_BUCKET, err)
			}

		}(uploader, file.Name())
	}

}
func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
