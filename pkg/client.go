package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// UploadFileToS3 downloads a file from a URL and uploads it to an S3 bucket with an optional endpoint.
func UploadFileToS3(ctx context.Context, bucketName, objectKey, fileURL, endpoint string) error {
	var cfg aws.Config
	var err error

	if endpoint != "" {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if service == s3.ServiceID {
				return aws.Endpoint{
					URL:           endpoint,
					SigningRegion: region,
				}, nil
			}
			// returning EndpointNotFoundError will allow the service to fall back to its default resolution
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		})

		cfg, err = config.LoadDefaultConfig(ctx, config.WithEndpointResolverWithOptions(customResolver))

	} else {
		cfg, err = config.LoadDefaultConfig(ctx)
	}

	if err != nil {
		return fmt.Errorf("unable to load SDK config, %v", err)
	}

	client := s3.NewFromConfig(cfg)

	resp, err := http.Get(fileURL)
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file, status code: %d", resp.StatusCode)
	}

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   resp.Body,
	})

	if err != nil {
		return fmt.Errorf("failed to upload file to S3: %v", err)
	}

	fmt.Printf("Successfully uploaded %s to s3://%s/%s\n", fileURL, bucketName, objectKey)
	return nil
}
