package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/harrypunk/save-to-s3/pkg/client" // Replace with your module name
)

type Event struct {
	BucketName string `json:"bucketName"`
	ObjectKey  string `json:"objectKey"`
	FileURL    string `json:"fileURL"`
	Endpoint   string `json:"endpoint"` // Optional endpoint
}

func HandleRequest(ctx context.Context, event Event) (string, error) {
	if event.BucketName == "" || event.ObjectKey == "" || event.FileURL == "" {
		return "", fmt.Errorf("missing required input parameters")
	}

	saveClient := 
	err := s3upload.UploadFileToS3(ctx, event.BucketName, event.ObjectKey, event.FileURL, event.Endpoint)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("File uploaded to s3://%s/%s", event.BucketName, event.ObjectKey), nil
}

func main() {
	lambda.Start(HandleRequest)
}
