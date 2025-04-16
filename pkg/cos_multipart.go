package feature

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/tencentyun/cos-go-sdk-v5"
)

type CosMultiClient struct {
	Endpoint string
	TCId     string
	TCKey    string
	BufSize  int
}

type uploadResult struct {
	Part int
	Err  error
}

func (cl *CosMultiClient) Save(ctx context.Context, bucketName, objectKey, fileURL string) error {
	u, _ := url.Parse(cl.Endpoint)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cl.TCId,
			SecretKey: cl.TCKey,
		},
	})

	httpBuf := HttpBuf{Size: cl.BufSize}

	result, _, err := client.Object.InitiateMultipartUpload(ctx, objectKey, nil)
	if err != nil {
		return fmt.Errorf("init multipart failed: %v", err)
	}

	uploadId := result.UploadID
	log.Printf("upload id: %s", uploadId)

	chunkCh := make(chan []byte)
	go httpBuf.Get(ctx, fileURL, chunkCh)

	resultCh := make(chan uploadResult)
	defer close(resultCh)

	totalParts := 0
	partNum := 1
	for chunk := range chunkCh {
		go uploadPart(ctx, client, objectKey, uploadId, partNum, chunk, resultCh)
		totalParts += 1
		partNum += 1
	}

	for range totalParts {
		rs := <-resultCh
		log.Printf("part result: %v", rs)
		if rs.Err != nil {
			return rs.Err
		}
	}

	resp, _, err := client.Object.CompleteMultipartUpload(ctx, objectKey, uploadId, &cos.CompleteMultipartUploadOptions{})
	if err != nil {
		return err
	}

	log.Printf("complete multipart %v\n", resp)

	return nil
}

func uploadPart(ctx context.Context, client *cos.Client, objectKey, uploadId string, part int, chunk []byte, resultChan chan<- uploadResult) {
	resp, err := client.Object.UploadPart(ctx, objectKey, uploadId, part, bytes.NewReader(chunk), nil)
	if err != nil {
		resultChan <- uploadResult{Part: part, Err: err}
	}
	partETag := resp.Header.Get("ETag")
	log.Printf("part %d etag: %s", part, partETag)
	resultChan <- uploadResult{Part: part, Err: nil}
}
