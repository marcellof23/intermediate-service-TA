package helper

import (
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/intermediate-service-ta/boot"
)

func ClientInitiation(clientType string, cli boot.Client) *s3.S3 {
	var client *s3.S3
	switch clientType {
	case "gcs":
		client = cli.GCSClient
	case "dos":
		client = cli.DOSClient
	case "s3":
		client = cli.S3Client
	default:
		client = cli.S3Client
	}

	return client
}
