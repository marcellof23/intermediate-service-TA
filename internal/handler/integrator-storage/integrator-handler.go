package integrator_storage

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"

	"github.com/intermediate-service-ta/boot"
	"github.com/intermediate-service-ta/helper"
)

type IntegratorHandler struct {
}

func NewIntegratorHandler() *IntegratorHandler {
	return &IntegratorHandler{}
}

func (hdl *IntegratorHandler) ListBuckets(c *gin.Context) {
	client, ok := c.MustGet("vdfsClient").(boot.Client)
	if !ok {
		fmt.Println("Failed to get session")
		return
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	fmt.Println("Buckets: ")
	for key, cli := range client.ClientMap {
		res, err := cli.ListBucketsWithContext(ctx, &s3.ListBucketsInput{})
		if err != nil {
			fmt.Printf("ListBucketsWithContext %s: %s", key, err.Error())
		}

		for _, b := range res.Buckets {
			fmt.Printf("%s\n", aws.StringValue(b.Name))
		}
	}
}

func (hdl *IntegratorHandler) UploadObject(c *gin.Context) {
	file, header, err := c.Request.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":        "No values in photo",
			"uploadObject": err.Error(),
		})
		return
	}

	clientType := c.PostForm("client")
	if clientType == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":        "Parameter in client not found",
			"uploadObject": err.Error(),
		})
		return
	}

	cli, ok := c.MustGet("vdfsClient").(boot.Client)
	if !ok {
		fmt.Println("Failed to get client")
		return
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	fileSize := header.Size
	fileBuffer := make([]byte, fileSize)
	file.Read(fileBuffer)

	client := helper.ClientInitiation(clientType, cli)

	_, err = client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("testing-vfs"),
		Key:    aws.String(header.Filename),
		Body:   bytes.NewReader(fileBuffer),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    "Failed to upload file",
			"uploader": err.Error(),
		})

	}
}

func (hdl *IntegratorHandler) DeleteObject(c *gin.Context) {
	clientType := c.PostForm("client")
	if clientType == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Parameter in client not found",
		})
		return
	}

	cli, ok := c.MustGet("vdfsClient").(boot.Client)
	if !ok {
		fmt.Println("Failed to get client")
		return
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	client := helper.ClientInitiation(clientType, cli)

	_, err := client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String("bucket_vfs_12"),
		Key:    aws.String("Screenshot from 2022-07-20 16-37-42.png"),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "Failed to upload file",
			"delete obj": err.Error(),
		})
	}
}
