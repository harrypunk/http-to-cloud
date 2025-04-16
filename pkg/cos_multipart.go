package feature

import (
	"context"
	"fmt"
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
	err string
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

	httpBuf := HttpBuf{
		Size: cl.BufSize
	}

	chunkCh := make(chan []byte)
	go httpBuf.Get(ctx, fileURL, chunkCh)
	for chunk := range ch {
		
	}

	return nil
}
