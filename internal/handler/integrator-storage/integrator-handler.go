package integrator_storage

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"

	"github.com/intermediate-service-ta/boot"
	"github.com/intermediate-service-ta/helper"
	"github.com/intermediate-service-ta/internal/repository"
)

type IntegratorHandler struct {
	fileRepo    repository.FileRepository
	requestRepo repository.RequestRepository
	Mu          sync.Mutex
}

func NewIntegratorHandler(filerepo repository.FileRepository) *IntegratorHandler {
	return &IntegratorHandler{fileRepo: filerepo}
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

	filename := c.Query("filename")

	file, err := hdl.fileRepo.Get(c, filename)
	if err != nil {
		fmt.Println("Failed to get client")
		return
	}

	cli, ok := c.MustGet("vdfsClient").(boot.Client)
	if !ok {
		fmt.Println("Failed to get client")
		return
	}

	client := helper.ClientInitiation(file.Client, cli)

	result, err := client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String("testing-vdfs"),
		Key:    aws.String(filename),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get file",
			"error":   err.Error(),
		})
		return
	}

	body, err := io.ReadAll(result.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to parse body",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successful get data",
		"data":    body,
	})
	return

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

func Tar(source, target string) error {
	filename := filepath.Base(source)
	target = filepath.Join(target, fmt.Sprintf("%s.tar", filename))
	tarfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer tarfile.Close()

	tarball := tar.NewWriter(tarfile)
	defer tarball.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	return filepath.Walk(source,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return err
			}

			if baseDir != "" {
				header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
			}

			if err := tarball.WriteHeader(header); err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(tarball, file)
			return err
		})
}

func (hdl *IntegratorHandler) GetFolder(c *gin.Context) {
	if err := Tar("backup", "."); err != nil {
		log.Fatal(err)
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.Header("Content-Type", "application/x-tar")
	c.Header("Content-Disposition", "attachment; filename=backup.tar")
	tarFile, err := os.Open("backup.tar")
	defer tarFile.Close()
	defer os.RemoveAll("backup.tar")

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	_, err = io.Copy(c.Writer, tarFile)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
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

func (hdl *IntegratorHandler) GetFolder2(c *gin.Context) {
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

func GetPath(client string, conf boot.Rclone) string {
	switch client {
	case "gcs":
		return conf.GcsConfig
	case "dos":
		return conf.DosConfig
	case "s3":
		return conf.S3Config
	}
	return ""
}

func (hdl *IntegratorHandler) MigrateObjects(c *gin.Context) {
	rcloneConf, ok := c.MustGet("rclone").(boot.Rclone)
	if !ok {
		fmt.Println("Failed to get rclone config")
		return
	}

	clientSource := c.PostForm("clientSource")
	if clientSource == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Parameter in clientSource not found",
		})
		return
	}

	clientDest := c.PostForm("clientDest")
	if clientDest == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Parameter in clientDest not found",
		})
		return
	}

	srcMigration := GetPath(clientSource, rcloneConf)
	dstMigration := GetPath(clientDest, rcloneConf)
	cmd := exec.Command("rclone", "moveto", srcMigration+"/", dstMigration+"/")

	// Run the rclone command and print any errors
	_, err := cmd.CombinedOutput()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = hdl.fileRepo.MigrateProvider(c, clientSource, clientDest)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Successfully migrate from %s to %s", clientSource, clientDest),
	})
}

func (hdl *IntegratorHandler) SyncOperations(c *gin.Context) {
	// Get the JSON body and decode into credentials
	RequestID := c.PostForm("RequestID")
	if RequestID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Parameter in RequestID not found",
		})
		return
	}

	TotalCommand := c.PostForm("TotalCommand")
	if TotalCommand == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Parameter in TotalCommand not found",
		})
		return
	}
	TotalCommandInt, err := strconv.Atoi(TotalCommand)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error",
		})
		return
	}

	if val, ok := boot.RequestCommand[RequestID]; ok {
		boot.RequestCommand[RequestID] = boot.CountRequest{TotalCommand: TotalCommandInt, TotalExecuted: val.TotalExecuted}
	} else {
		boot.RequestCommand[RequestID] = boot.CountRequest{TotalCommand: TotalCommandInt, TotalExecuted: 0}
	}

	var isSuccess = false
	cr := cron.New()
	var wg sync.WaitGroup
	wg.Add(1)
	cr.AddFunc("@every 0.5s", func() {
		hdl.Mu.Lock()
		if val, ok := boot.RequestCommand[RequestID]; ok {
			if val.TotalCommand == val.TotalExecuted && val.TotalExecuted != 0 {
				isSuccess = true
				wg.Done()
				cr.Stop()
			}
		}
		hdl.Mu.Unlock()
	})
	cr.Start()
	wg.Wait()

	if isSuccess {
		c.AbortWithStatusJSON(200, gin.H{
			"success": true,
		})
		return
	} else {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   "err",
		})
		return
	}
}
