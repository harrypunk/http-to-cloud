package feature

import (
	"context"
)

type Downloader interface {
	Save(ctx context.Context, bucket, objectKey, url string) error
}
