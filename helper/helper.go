package helper

import (
	"context"
	"errors"
	"log"
	"path/filepath"
	"sort"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/intermediate-service-ta/boot"
)

var (
	ErrDBNotFound       = errors.New("failed to get database from context")
	ErrJWTKeyNotFound   = errors.New("failed to get jwt key from context")
	ErrUsernameNotFound = errors.New("failed to get username from context")
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

func GetVDFSClientFromContext(c context.Context) (boot.Client, error) {
	tmp := c.Value("vdfsClient")
	db, ok := tmp.(boot.Client)
	if !ok {
		return boot.Client{}, ErrDBNotFound
	}
	return db, nil
}

func GetBucketNameFromContext(c context.Context) (string, error) {
	tmp := c.Value("bucketName")
	bucket, ok := tmp.(string)
	if !ok {
		return "", ErrDBNotFound
	}
	return bucket, nil
}

func GetUsernameFromContext(c context.Context) (string, error) {
	tmp := c.Value("username")
	uname, ok := tmp.(string)
	if !ok {
		return "", ErrUsernameNotFound
	}
	return uname, nil
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

func GetClientsFromContext(c *gin.Context) ([]string, error) {
	tmp, exists := c.Get("clients")
	if !exists {
		return []string{}, ErrJWTKeyNotFound
	}
	clients, ok := tmp.([]string)
	if !ok {
		return []string{}, ErrJWTKeyNotFound
	}
	return clients, nil
}

func GetConfigFromContext(c context.Context) (boot.Config, error) {
	tmp := c.Value("config")
	conf, ok := tmp.(boot.Config)
	if !ok {
		return boot.Config{}, ErrJWTKeyNotFound
	}
	return conf, nil
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

func GetLogger(ctx context.Context) *zap.Logger {
	if ctx.Value("pubsub-logger") != nil {
		return ctx.Value("pubsub-logger").(*zap.Logger)
	}

	return zap.NewNop()
}

func JoinPath(path ...string) string {
	joinedPath := filepath.Join(path...)
	unixPath := filepath.ToSlash(joinedPath)
	return unixPath
}

func SortSlice(m map[string]int64) []string {
	keys := make([]string, 0, len(m))

	for key := range m {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return m[keys[i]] < m[keys[j]]
	})

	return keys
}
