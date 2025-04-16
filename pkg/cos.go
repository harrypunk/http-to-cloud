package feature

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// url must have bucket-appid: https://examplebucket-1250000000.cos.ap-guangzhou.myqcloud.com
// file less than 5G
type CosPutClient struct {
	Endpoint string
	TCId     string
	TCKey    string
}

func (cl *CosPutClient) Save(ctx context.Context, bucketName, objectKey, fileURL string) error {
	u, _ := url.Parse(cl.Endpoint)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cl.TCId,
			SecretKey: cl.TCKey,
		},
	})

	resp, err := http.Get(fileURL)
	if err != nil {
		return fmt.Errorf("http get error: %v", err)
	}
	defer resp.Body.Close()

	res, err := client.Object.Put(
		ctx, objectKey, &logHttp{src: resp.Body}, nil,
	)
	if err != nil {
		return fmt.Errorf("object put error: %v", err)
	}
	fmt.Printf("%+v\n", *res)
	return nil
}
