package client

import (
	"context"
	"github.com/joho/godotenv"
	"os"
	"testing"
)

func TestSave(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Errorf("load dotenv error: %v", err)
	}

	fileUrl := os.Getenv("SRC_FILE_URL")
	endpoint := os.Getenv("DEST_S3_ENDPOINT")
	buck := os.Getenv("DEST_S3_BUCKET")
	key := os.Getenv("DEST_S3_key")
	saveClient := Client{
		Endpoint: endpoint,
	}

	ctx := context.Background()
	err = saveClient.Save(ctx, buck, key, fileUrl)
	if err != nil {
		t.Errorf("client error: %v", err)
	}
}
