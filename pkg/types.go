package feature

import (
	"context"
	"io"
	"log"
)

type Downloader interface {
	Save(ctx context.Context, bucket, objectKey, url string) error
}

type logHttp struct {
	src io.Reader
}

func (lh *logHttp) Read(b []byte) (n int, err error) {
	n, err = lh.src.Read(b)
	log.Printf("log http length %v\n", n)
	return n, err
}
