package helper

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/intermediate-service-ta/boot"
)

var ErrDBNotFound = errors.New("failed to get database from context")
var ErrJWTKeyNotFound = errors.New("failed to get jwt key from context")

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

func GetDatabaseFromContext(c *gin.Context) (*gorm.DB, error) {
	tmp, exists := c.Get("db")
	if !exists {
		return nil, ErrDBNotFound
	}
	db, ok := tmp.(*gorm.DB)
	if !ok {
		return nil, ErrDBNotFound
	}
	return db, nil
}

func GetJWTSecretFromContext(c *gin.Context) (string, error) {
	tmp, exists := c.Get("jwt-secret")
	if !exists {
		return "", ErrJWTKeyNotFound
	}
	jwtkey, ok := tmp.(string)
	if !ok {
		return "", ErrJWTKeyNotFound
	}
	return jwtkey, nil
}

func GetHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}
