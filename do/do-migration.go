package do

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func ListDOBuckets(w io.Writer) error {
	key := "DO00RPCBRKL48QL6FK8W"
	secret := "8r7PcHaYzOZEXAoZiQvabSZqTem5Z8FS5jAjKRiR3GU"

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:         aws.String("https://sgp1.digitaloceanspaces.com"),
		Region:           aws.String("us-east-1"),
		S3ForcePathStyle: aws.Bool(false), // // Configures to use subdomain/virtual calling format. Depending on your version, alternatively use o.UsePathStyle = false
	}

	newSession := session.New(s3Config)
	s3Client := s3.New(newSession)

	spaces, err := s3Client.ListBuckets(nil)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	for _, b := range spaces.Buckets {
		fmt.Println(aws.StringValue(b.Name))
	}

	return nil
}
