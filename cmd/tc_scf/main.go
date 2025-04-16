package main

import (
	"context"
	"fmt"
	"os"

	"github.com/harrypunk/http-to-cloud/pkg"
	"github.com/tencentyun/scf-go-lib/cloudfunction"
)

type Info struct {
	Url         string `json:"url"`
	CosEndpoint string `json:"cos-endpoint"`
	CosKey      string `json:"cos-key"`
}

func putToCos(ctx context.Context, event Info) (string, error) {
	fmt.Println(fmt.Sprintf("Received event: %v", event))

	client := feature.CosPutClient{
		Endpoint: event.CosEndpoint,
		TCId:     os.Getenv("TC_ID"),
		TCKey:    os.Getenv("TC_KEY"),
	}

	err := client.Save(ctx, "", event.CosKey, event.Url)
	if err != nil {
		return "cos save error", err
	}

	return "cos put ok", nil
}

func main() {
	cloudfunction.Start(putToCos)
}
