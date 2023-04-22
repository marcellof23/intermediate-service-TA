package helper

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/intermediate-service-ta/boot"
)

var (
	ErrDBNotFound     = errors.New("failed to get database from context")
	ErrJWTKeyNotFound = errors.New("failed to get jwt key from context")
)

func ClientInitiation(clientType string, cli boot.Client) *s3.S3 {
	var client *s3.S3
	switch clientType {
	case "gcs":
		client = cli.ClientMap["gcs"]
	case "dos":
		client = cli.ClientMap["dos"]
	case "s3":
		client = cli.ClientMap["s3"]
	default:
		client = cli.ClientMap["s3"]
	}
	return client
}

func GetDatabaseFromContext(c context.Context) (*gorm.DB, error) {
	tmp := c.Value("db")
	db, ok := tmp.(*gorm.DB)
	if !ok {
		return nil, ErrDBNotFound
	}
	return db, nil
}

func GetJWTSecretFromContext(c *gin.Context) (string, error) {
	tmp, exists := c.Get("jwtSecret")
	if !exists {
		return "", ErrJWTKeyNotFound
	}
	jwtKey, ok := tmp.(string)
	if !ok {
		return "", ErrJWTKeyNotFound
	}
	return jwtKey, nil
}

func GetJWTSecretFromContextQueue(c context.Context) (string, error) {
	tmp := c.Value("jwtSecret")
	jwtKey, ok := tmp.(string)
	if !ok {
		return "", ErrJWTKeyNotFound
	}
	return jwtKey, nil
}

func GetHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}
