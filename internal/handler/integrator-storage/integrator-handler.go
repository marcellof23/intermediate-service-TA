package integrator_storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
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

func (hdl *IntegratorHandler) GetFile(c *gin.Context) {
	clientType := c.PostForm("client")
	cli, ok := c.MustGet("vdfsClient").(boot.Client)
	if !ok {
		fmt.Println("Failed to get client")
		return
	}

	client := helper.ClientInitiation(clientType, cli)
	filename := "ehe"

	result, err := client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String("testing-vfs"),
		Key:    aws.String(filename),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    "Failed to get file",
			"uploader": err.Error(),
		})
		return
	}

	defer result.Body.Close()
	file, err := os.Create(filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    "Failed to get file",
			"uploader": err.Error(),
		})
		return
	}
	defer file.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    "Failed to get file",
			"uploader": err.Error(),
		})
		return
	}
	_, err = file.Write(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    "Failed to get file",
			"uploader": err.Error(),
		})
		return
	}
}

func (hdl *IntegratorHandler) UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "No values in photo",
			"error":   err.Error(),
		})
		return
	}

	clientType := c.PostForm("client")
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
		Key:    aws.String("./asdf/" + header.Filename),
		Body:   bytes.NewReader(fileBuffer),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to upload file",
			"error":   err.Error(),
		})
	}
}

func (hdl *IntegratorHandler) DeleteFile(c *gin.Context) {
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
		Bucket: aws.String("testing-vdfs"),
		Key:    aws.String("h/logger.go"),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete file",
			"error":   err.Error(),
		})
	}
}

func (hdl *IntegratorHandler) TestFile(c *gin.Context) {
}
