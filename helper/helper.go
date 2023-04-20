package helper

import (
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/intermediate-service-ta/boot"
)

func ClientInitiation(clientType string, cli boot.Client) *s3.S3 {
	var client *s3.S3
	switch clientType {
	case "gcs":
		client = cli.ClientMap["gcs"]
	case "dos":
		client = cli.ClientMap["dos"]
	case "s3":
		client = cli.ClientMap["s3"]
	default:
		client = cli.ClientMap["s3"]
	}

	return client
}
