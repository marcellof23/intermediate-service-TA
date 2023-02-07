package gcs

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

type GCSHandler struct {
}

func NewGCSHandler() *GCSHandler {
	return &GCSHandler{}
}

func (hdl *GCSHandler) ListGCSBuckets(c *gin.Context) {
	gcsSession, ok := c.MustGet("gcsSession").(*session.Session)
	if !ok {
		fmt.Println("Failed to get google session")
		return
	}

	client := s3.New(gcsSession)
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	result, err := client.ListBucketsWithContext(ctx, &s3.ListBucketsInput{})
	if err != nil {
		fmt.Printf("ListBucketsWithContext: %s", err.Error())
	}

	fmt.Println("Buckets:")
	for _, b := range result.Buckets {
		fmt.Printf("%s\n", aws.StringValue(b.Name))
	}

}
