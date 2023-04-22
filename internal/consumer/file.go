package consumer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang-jwt/jwt"

	"github.com/intermediate-service-ta/boot"
	"github.com/intermediate-service-ta/helper"
	"github.com/intermediate-service-ta/internal/model"
)

func verifyJWTQueue(ctx context.Context, msg Message) error {
	if msg.Token != "" {
		secretKey, err := helper.GetJWTSecretFromContextQueue(ctx) // Get secret key if exist
		if err != nil {
			return err
		}

		token, err := jwt.Parse(msg.Token, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return "", errors.New("unauthorized")
			}
			return []byte(secretKey), nil
		})

		// parsing errors result
		if err != nil {
			return err
		}
		// if there's a token
		if token.Valid {
			return nil
		} else {
			return errors.New("invalid token")
		}
	} else {
		return errors.New("no token in the header")
	}
}

func (con *Consumer) exec(c context.Context, msg Message, log *log.Logger) error {
	err := verifyJWTQueue(c, msg)
	if err != nil {
		return err
	}

	comms := strings.Split(msg.Command, " ")
	switch comms[0] {
	case "upload":
		con.UploadFile(c, msg)
	case "rm":
		if len(comms) > 1 && comms[1] == "-r" {
		} else {
		}
	case "mkdir":
	case "touch":
	default:
		fmt.Println(comms[0], ": Command not found")
		return errors.New("command not found")
	}

	return nil
}

func (con *Consumer) UploadFile(c context.Context, msg Message) {
	file := model.File{
		Filename:     msg.AbsPath,
		OriginalName: msg.AbsPath,
		Client:       "gcs",
		Size:         int64(len(msg.Buffer)),
	}

	fl, err := con.fileRepo.Create(c, &file)
	if err != nil {
		fmt.Println(err)
		return
	}

	cli, ok := c.Value("vdfsClient").(boot.Client)
	if !ok {
		fmt.Println("Failed to get client")
		return
	}
	client := helper.ClientInitiation("gcs", cli)

	_, err = client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("bucket_vfs_1"),
		Key:    aws.String(fl.Filename),
		Body:   bytes.NewReader(msg.Buffer),
	})
	if err != nil {
		fmt.Println(err)
		return
	}

}
