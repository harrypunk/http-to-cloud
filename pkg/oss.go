package feature

import (
	"context"
	"fmt"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"net/http"
)

type OssClient struct {
	// Endpoint string
	Region string
}

func (cl *OssClient) Save(ctx context.Context, bucketName, objectKey, fileURL string) error {

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(cl.Region)

	resp, err := http.Get(fileURL)
	if err != nil {
		return fmt.Errorf("http get error: %v", err)
	}
	defer resp.Body.Close()

	client := oss.NewClient(cfg)
	uploader := client.NewUploader()

	result, err := uploader.UploadFrom(context.TODO(),
		&oss.PutObjectRequest{
			Bucket: oss.Ptr(bucketName),
			Key:    oss.Ptr(objectKey)},
		resp.Body)
	if err != nil {
		return fmt.Errorf("failed with uploadfrom : %v", err)
	}

	fmt.Printf("upload file result:%#v\n", result)
	return nil
}
