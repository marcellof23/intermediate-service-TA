package gcs

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	iam "google.golang.org/api/iam/v1"
)

func ListGCSBuckets(w io.Writer, googleAccessKeyID string, googleAccessKeySecret string) error {
	// googleAccessKeyID := "Your Google Access Key ID"
	// googleAccessKeySecret := "Your Google Access Key Secret"

	// Create a new client and do the following:
	// 1. Change the endpoint URL to use the Google Cloud Storage XML API endpoint.
	// 2. Use Cloud Storage HMAC Credentials.
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("auto"),
		Endpoint:    aws.String("https://storage.googleapis.com"),
		Credentials: credentials.NewStaticCredentials(googleAccessKeyID, googleAccessKeySecret, ""),
	}))

	client := s3.New(sess)
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	result, err := client.ListBucketsWithContext(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return fmt.Errorf("ListBucketsWithContext: %v", err)
	}

	fmt.Fprintf(w, "Buckets:")
	for _, b := range result.Buckets {
		fmt.Fprintf(w, "%s\n", aws.StringValue(b.Name))
	}

	return nil
}

func ListKeys(w io.Writer, serviceAccountEmail string) ([]*iam.ServiceAccountKey, error) {
	ctx := context.Background()
	service, err := iam.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("iam.NewService: %v", err)
	}

	resource := "projects/-/serviceAccounts/" + serviceAccountEmail
	response, err := service.Projects.ServiceAccounts.Keys.List(resource).Do()
	if err != nil {
		return nil, fmt.Errorf("Projects.ServiceAccounts.Keys.List: %v", err)
	}
	for _, key := range response.Keys {
		fmt.Fprintf(w, "Listing key: %v", key.Name)
	}
	return response.Keys, nil
}
