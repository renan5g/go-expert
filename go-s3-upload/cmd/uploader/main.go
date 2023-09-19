package main

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	s3Client *s3.S3
	s3Bucket string
	wg       sync.WaitGroup
)

func init() {
	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials("minioadmin", "minioadmin", ""),
		Endpoint:    aws.String("http://localhost:9000"),
		Region:      aws.String("us-east-1"),
	}
	newSession, err := session.NewSession(s3Config)
	if err != nil {
		panic(err)
	}

	s3Client = s3.New(newSession)
	s3Bucket = "goexpert"
}

func main() {
	dir, err := os.Open("./tmp")
	if err != nil {
		panic(err)
	}
	defer dir.Close()

	upControl := make(chan struct{}, 100)
	errorFileUpload := make(chan string, 10)

	go func() {
		for {
			select {
			case filename := <-errorFileUpload:
				upControl <- struct{}{}
				wg.Add(1)
				go uploadFile(filename, upControl, errorFileUpload)
			}
		}
	}()

	for {
		files, err := dir.ReadDir(1)
		if err != nil {
			if err == io.EOF {
				break
			}

			fmt.Printf("Error reading directory: %s\n", err)
			continue
		}
		wg.Add(1)
		upControl <- struct{}{}
		go uploadFile(files[0].Name(), upControl, errorFileUpload)
	}

	wg.Wait()
}

func uploadFile(filename string, upControl <-chan struct{}, errorFileUpload chan<- string) {
	defer wg.Done()
	completeFilename := fmt.Sprintf("./tmp/%s", filename)

	fmt.Printf("Uploading file %s starting to bucket %s\n", completeFilename, s3Bucket)

	f, err := os.Open(completeFilename)
	if err != nil {
		fmt.Printf("Error opening file %s\n", completeFilename)
		<-upControl
		errorFileUpload <- filename
		return
	}
	defer f.Close()

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(filename),
		Body:   f,
	})
	if err != nil {
		fmt.Printf("Error uploading file %s, error: %v\n", completeFilename, err.Error())
		<-upControl
		errorFileUpload <- filename
		return
	}

	fmt.Printf("File %s uploaded successfully\n", completeFilename)
	<-upControl
}
