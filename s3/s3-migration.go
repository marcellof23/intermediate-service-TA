package s3

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Handler struct {
}

func NewS3Handler() *S3Handler {
	return &S3Handler{}
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
}

func (hdl *S3Handler) ListS3Objects(c *gin.Context) {
	s3Session, ok := c.MustGet("s3Session").(*session.Session)
	if !ok {
		fmt.Println("Failed to get google session")
		return
	}

	bucket := "testing-vdfs"

	// Create S3 service client
	svc := s3.New(s3Session)

	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
	if err != nil {
		exitErrorf("Unable to list items in bucket %q, %v", bucket, err)
	}

	for _, item := range resp.Contents {
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		fmt.Println("")
	}
}
