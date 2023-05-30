package consumer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang-jwt/jwt"

	"github.com/intermediate-service-ta/boot"
	"github.com/intermediate-service-ta/helper"
	"github.com/intermediate-service-ta/internal/model"
	"github.com/intermediate-service-ta/internal/storage"
)

func (con *Consumer) exec(c context.Context, msg Message, log *log.Logger) error {
	uname, err := helper.GetUsernameFromContext(c)
	if err != nil {
		log.Println(err)
	}

	log.Println(msg.Command, msg.AbsPathSource, uname, msg.Uid, msg.Gid, msg.FileMode, len(msg.Buffer), msg.Order)

	comms := strings.Split(msg.Command, " ")
	switch comms[0] {
	case "upload":
		con.UploadFile(c, msg)
	case "chmod":
		con.ChangeFileMode(c, msg)
	case "cp":
		con.CopyFile(c, msg)
	case "rm":
		if len(comms) > 1 && comms[1] == "-r" {
			con.RemoveDir(c, msg)
		} else {
			con.RemoveFile(c, msg)
		}
	case "mkdir":
		con.CreateFolder(c, msg)
	case "write":
		con.WriteFile(c, msg)
	default:
		return errors.New("command not found")
	}

	return nil
}

type Effector func(context.Context, Message) error

func (con *Consumer) Retry(effector Effector, delay time.Duration) Effector {
	return func(ctx context.Context, msg Message) error {
		for r := 0; ; r++ {
			err := effector(ctx, msg)
			if err == nil || r >= 40 {
				if r >= 40 {
					con.errorLog.Printf("Retrying function already done 40 times, error still persist err: %s", err.Error())
				}
				return nil
			}

			select {
			case <-time.After(delay):
				if r >= 20 {
					delay = delay * time.Duration(1.2*float64(time.Second))
					fmt.Println(delay)
				}
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

func (con *Consumer) AuthQueue(ctx context.Context, msg Message, log *log.Logger) error {
	if msg.Token != "" {
		secretKey, err := helper.GetJWTSecretFromContextQueue(ctx) // Get secret key if exist
		if err != nil {
			return err
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(msg.Token, claims, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return "", errors.New("unauthorized")
			}
			return []byte(secretKey), nil
		})
		ctx = context.WithValue(ctx, "username", claims["username"])

		// parsing errors result
		if err != nil {
			return err
		}
		// if there's a token
		if token.Valid {
			err = con.exec(ctx, msg, log)
			if err != nil {
				fmt.Println(err)
			}
			return nil
		} else {
			return errors.New("invalid token")
		}
	} else {
		return errors.New("no token in the header")
	}
}

func (con *Consumer) UploadFile(c context.Context, msg Message) {
	arrRes := helper.SortSlice(storage.TotalSizeClient)

	fullPath := helper.JoinPath(msg.AbsPathDest, msg.AbsPathSource)
	file := model.File{
		Filename:     fullPath,
		OriginalName: fullPath,
		Client:       arrRes[0],
		Size:         int64(len(msg.Buffer)),
	}
	storage.UpdateTotalSizeClient(arrRes[0], int64(len(msg.Buffer)))

	fl, err := con.fileRepo.Create(c, &file)
	if err != nil {
		con.errorLog.Println(err)
		return
	}

	cli, err := helper.GetVDFSClientFromContext(c)
	if err != nil {
		con.errorLog.Println(err)
		return
	}
	client := helper.ClientInitiation(arrRes[0], cli)

	bucketName, err := helper.GetBucketNameFromContext(c)
	if err != nil {
		con.errorLog.Println(err)
		return
	}

	_, err = client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fl.Filename),
		Body:   bytes.NewReader(msg.Buffer),
	})
	if err != nil {
		con.errorLog.Println(err)
		return
	}

	r := con.Retry(BackupFiletoDisk, 1e9)
	go r(c, msg)
}

func (con *Consumer) WriteFile(c context.Context, msg Message) {
	arrRes := helper.SortSlice(storage.TotalSizeClient)

	fullPath := helper.JoinPath(msg.AbsPathDest, msg.AbsPathSource)

	fl, err := con.fileRepo.Get(c, fullPath)
	if err != nil {
		con.errorLog.Println(err)
		return
	}

	chunkFile := model.ChunkFile{
		Filename: fl.Filename,
		FileID:   fl.ID,
		Order:    msg.Order,
		Client:   arrRes[0],
		Size:     int64(len(msg.Buffer)),
	}

	_, err = con.chunkFileRepo.Create(c, &chunkFile)
	if err != nil {
		con.errorLog.Println(err)
		return
	}

	storage.UpdateTotalSizeClient(arrRes[0], int64(len(msg.Buffer)))

	cli, err := helper.GetVDFSClientFromContext(c)
	if err != nil {
		con.errorLog.Println(err)
		return
	}
	client := helper.ClientInitiation(arrRes[0], cli)

	bucketName, err := helper.GetBucketNameFromContext(c)
	if err != nil {
		con.errorLog.Println(err)
		return
	}

	_, err = client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fmt.Sprintf("%s-partition-%d", fl.Filename, msg.Order)),
		Body:   bytes.NewReader(msg.Buffer),
	})
	if err != nil {
		con.errorLog.Println(err)
		return
	}

	r := con.Retry(WriteFileOnDisk, 1e9)
	go r(c, msg)
}

func (con *Consumer) ChangeFileMode(c context.Context, msg Message) {
	fullPath := helper.JoinPath(boot.Backup, msg.AbsPathSource)

	err := os.Chmod(fullPath, os.FileMode(msg.FileMode))
	if err != nil {
		con.errorLog.Println(err)
		return
	}
}

func (con *Consumer) CreateFolder(c context.Context, msg Message) {
	fullPath := helper.JoinPath(boot.Backup, msg.AbsPathSource)
	err := os.MkdirAll(fullPath, os.FileMode(msg.FileMode))
	os.Chown(fullPath, msg.Uid, msg.Gid)

	if err != nil {
		con.errorLog.Println(err)
		return
	}
}

func (con *Consumer) RemoveFile(c context.Context, msg Message) {
	file, err := con.fileRepo.Delete(c, msg.AbsPathSource)
	if err != nil {
		con.errorLog.Println(err)
		return
	}

	cli, err := helper.GetVDFSClientFromContext(c)
	if err != nil {
		con.errorLog.Println(err)
		return
	}

	client := helper.ClientInitiation(file.Client, cli)

	bucketName, err := helper.GetBucketNameFromContext(c)
	if err != nil {
		con.errorLog.Println(err)
		return
	}

	_, err = client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(file.Filename),
	})

	r := con.Retry(RemoveFileFromDisk, 1e9)
	go r(c, msg)
}

func (con *Consumer) CopyFile(c context.Context, msg Message) {
	flSource, err := con.fileRepo.Get(c, msg.AbsPathSource)
	if err != nil {
		con.errorLog.Println(err)
		return
	}

	file := model.File{
		Filename:     msg.AbsPathDest,
		OriginalName: msg.AbsPathDest,
		Client:       flSource.Client,
		Size:         flSource.Size,
	}
	storage.UpdateTotalSizeClient(flSource.Client, int64(len(msg.Buffer)))

	fl, err := con.fileRepo.Create(c, &file)
	if err != nil {
		con.errorLog.Println(err)
		return
	}

	cli, err := helper.GetVDFSClientFromContext(c)
	if err != nil {
		con.errorLog.Println(err)
		return
	}
	client := helper.ClientInitiation(flSource.Client, cli)

	bucketName, err := helper.GetBucketNameFromContext(c)
	if err != nil {
		con.errorLog.Println(err)
		return
	}

	_, err = client.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(bucketName),
		Key:        aws.String(fl.Filename),
		CopySource: aws.String(fmt.Sprintf("%v/%v", bucketName, msg.AbsPathSource)),
	})
	if err != nil {
		con.errorLog.Println(err)
		return
	}

	r := con.Retry(CopyFiletoDisk, 1e9)
	go r(c, msg)

}

func (con *Consumer) RemoveDir(c context.Context, msg Message) {
	r := con.Retry(RemoveFolderFromDisk, 1e9)
	go r(c, msg)
}
