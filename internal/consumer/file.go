package consumer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
		fmt.Println(err)
	}

	log.Println(msg.Command, msg.AbsPathSource, uname)

	comms := strings.Split(msg.Command, " ")
	switch comms[0] {
	case "upload":
		con.UploadFile(c, msg)
	case "cp":
		if len(comms) > 1 && comms[1] == "-r" {
			con.CopyDir(c, msg)
		} else {
			con.CopyFile(c, msg)
		}
	case "rm":
		if len(comms) > 1 && comms[1] == "-r" {
			con.RemoveDir(c, msg)
		} else {
			con.RemoveFile(c, msg)
		}
	case "mkdir":
		con.CreateFolder(c, msg)
	default:
		return errors.New("command not found")
	}

	return nil
}

type Effector func(context.Context, Message) error

func (con *Consumer) Retry(effector Effector, delay time.Duration) Effector {
	return func(ctx context.Context, msg Message) error {
		for {
			err := effector(ctx, msg)
			if err == nil {
				return nil
			}

			con.errorLog.Printf("Function call failed, retrying in %v err: %s", delay, err.Error())

			select {
			case <-time.After(delay):
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
	fullPath := filepath.Join(msg.AbsPathDest, msg.AbsPathSource)
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

	r := con.Retry(BackupFiletoDisk, 3e9)
	go r(c, msg)
}

func (con *Consumer) CreateFolder(c context.Context, msg Message) {
	err := os.MkdirAll(filepath.Join(boot.Backup, msg.AbsPathSource), os.ModePerm)
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

	err = os.Remove(filepath.Join(boot.Backup, msg.AbsPathSource))
	if err != nil {
		con.errorLog.Println(err)
		return
	}
}

func (con *Consumer) CopyFile(c context.Context, msg Message) {
	flSource, err := con.fileRepo.Get(c, msg.AbsPathSource)
	if err != nil {
		con.errorLog.Println(err)
		return
	}

	arrRes := helper.SortSlice(storage.TotalSizeClient)
	file := model.File{
		Filename:     msg.AbsPathSource,
		OriginalName: msg.AbsPathSource,
		Client:       arrRes[0],
		Size:         flSource.Size,
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

	_, err = client.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(bucketName),
		Key:        aws.String(fl.Filename),
		CopySource: aws.String(msg.AbsPathSource),
	})

	err = CopyFiletoDisk(c, msg.AbsPathSource, msg.AbsPathDest)
	if err != nil {
		con.errorLog.Println(err)
		return
	}
}

func (con *Consumer) CopyDir(c context.Context, msg Message) {
	//arrRes := helper.SortSlice(storage.TotalSizeClient)
	//fullPath := filepath.Join(msg.AbsPathDest, msg.AbsPathSource)
	//file := model.File{
	//	Filename:     fullPath,
	//	OriginalName: fullPath,
	//	Client:       arrRes[0],
	//	Size:         int64(len(msg.Buffer)),
	//}
	//storage.UpdateTotalSizeClient(arrRes[0], int64(len(msg.Buffer)))

}

func (con *Consumer) RemoveDir(c context.Context, msg Message) {
}
