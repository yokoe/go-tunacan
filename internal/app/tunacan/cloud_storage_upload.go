package tunacan

import (
	"context"
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/storage"
)

func uploadToCloudStorage(localFilepath string, bucketName string, remoteFilename string, contentType string) (*string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	bucket := client.Bucket(bucketName)
	object := bucket.Object(remoteFilename)
	writer := object.NewWriter(ctx)

	writer.ObjectAttrs.ContentType = contentType
	writer.ObjectAttrs.CacheControl = "no-cache"
	writer.ObjectAttrs.ACL = []storage.ACLRule{
		storage.ACLRule{
			Entity: storage.AllUsers,
			Role:   storage.RoleReader,
		},
	}

	f, err := os.Open(localFilepath)
	if _, err = io.Copy(writer, f); err != nil {
		return nil, err
	}
	defer f.Close()

	if err := writer.Close(); err != nil {
		return nil, err
	}

	attr := writer.Attrs()
	uploadedURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", attr.Bucket, attr.Name)

	return &uploadedURL, nil
}
