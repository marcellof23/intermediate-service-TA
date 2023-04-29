package integrator_storage

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

func zipSource(source, target string) error {
	// 1. Create a ZIP file and zip.Writer
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	// 2. Go through all the files of the source
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 3. Create a local file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// set compression
		header.Method = zip.Deflate

		// 4. Set relative path of a file as the header name
		header.Name, err = filepath.Rel(filepath.Dir(source), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += "/"
		}

		// 5. Create writer for the file header and save content of the file
		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		return err
	})
}

func (hdl *IntegratorHandler) GetFolder(c *gin.Context) {
	if err := zipSource("backup", "backup.zip"); err != nil {
		log.Fatal(err)
	}
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", "attachment; filename=backup.zip")
	zipFile, err := os.Open("backup.zip")
	defer zipFile.Close()
	defer os.RemoveAll("backup.zip")

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	_, err = io.Copy(c.Writer, zipFile)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
