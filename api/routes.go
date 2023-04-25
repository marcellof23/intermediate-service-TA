package api

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/intermediate-service-ta/boot"
	"github.com/intermediate-service-ta/internal/consumer"
	integratehandler "github.com/intermediate-service-ta/internal/handler/integrator-storage"
	userhandler "github.com/intermediate-service-ta/internal/handler/user"
	repository "github.com/intermediate-service-ta/internal/storage"
)

func newClient(endpoint, region, accessKeyID, secretKey string) *session.Session {
	cl := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	awsConfig := &aws.Config{
		Region:           aws.String(region),
		S3ForcePathStyle: aws.Bool(true),
		HTTPClient:       cl,
		Credentials:      credentials.NewStaticCredentials(accessKeyID, secretKey, ""),
	}

	if endpoint != "" {
		awsConfig.Endpoint = aws.String(endpoint)
	}
	s := session.Must(session.NewSession(awsConfig))

	return s
}

func initSession(dep *boot.Dependencies) boot.Sess {
	gcsSession := newClient(dep.Config().GCSEndpoint, dep.Config().GCSRegion, dep.Config().GoogleAccessKeyID, dep.Config().GoogleAccessKeySecret)
	dosSession := newClient(dep.Config().DOSEndpoint, dep.Config().DOSRegion, dep.Config().DigitalOceanAccessKeyID, dep.Config().DigitalOceanAccessKeySecret)
	s3Session := newClient(dep.Config().S3Endpoint, dep.Config().S3Region, dep.Config().AmazonAccessKeyID, dep.Config().AmazonAccessKeySecret)

	sessionMap := make(map[string]*session.Session)
	sessionMap["gcs"] = gcsSession
	sessionMap["dos"] = dosSession
	sessionMap["s3"] = s3Session

	return boot.Sess{
		SessionMap: sessionMap,
	}
}

func initClient(dep *boot.Dependencies, sess boot.Sess) boot.Client {
	var clientMap = make(map[string]*s3.S3)
	for _, v := range boot.Clients {
		clientMap[v] = s3.New(sess.SessionMap[v])
		_, err := clientMap[v].CreateBucket(&s3.CreateBucketInput{
			Bucket: aws.String(dep.Config().BucketName),
		})
		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				if awsErr.Code() == s3.ErrCodeBucketAlreadyOwnedByYou || awsErr.Code() == s3.ErrCodeBucketAlreadyExists {
					continue
				} else {
					fmt.Println(err.Error())
				}
			} else {
				fmt.Println(err.Error())
			}
		}
	}

	return boot.Client{
		ClientMap: clientMap,
	}
}

func InitRoutes(dep *boot.Dependencies) *gin.Engine {
	// init repos
	userRepo := repository.NewUserRepo()
	fileRepo := repository.NewFileRepo()

	// init Handler
	integrateHdl := integratehandler.NewIntegratorHandler()
	userHdl := userhandler.NewUserHandler(userRepo)

	// init blank engine
	r := gin.New()

	// init session
	sess := initSession(dep)
	client := initClient(dep, sess)

	// attach session to context
	r.Use(func(c *gin.Context) {
		c.Set("vdfsClient", client)
		c.Set("db", dep.DB())
		c.Set("jwtSecret", dep.Config().JWTSecretKey)
	})

	// attach value to context
	ctx := context.Background()
	ctx = context.WithValue(ctx, "db", dep.DB())
	ctx = context.WithValue(ctx, "vdfsClient", client)
	ctx = context.WithValue(ctx, "jwtSecret", dep.Config().JWTSecretKey)
	ctx = context.WithValue(ctx, "bucketName", dep.Config().BucketName)

	// init total size client
	var errs error
	repository.TotalSizeClient, errs = fileRepo.GetTotalSizeClient(ctx)
	if errs != nil {
		panic(errors.New("failed to get total client size"))
	}

	// init consumer

	consumer := consumer.NewConsumer(fileRepo)
	go consumer.ConsumeCommand(ctx, dep)

	// setup cors config
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length", "Authorization", "Content-type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "ok")
	})

	apiV1 := r.Group("/api/v1")
	{
		authRoutes := apiV1.Group("/")
		authRoutes.Use(userhandler.VerifyJWT)

		apiV1File := authRoutes.Group("/file")
		{
			apiV1File.GET("/list-bucket", integrateHdl.ListBuckets)
			apiV1File.GET("/object", integrateHdl.GetFile)
			apiV1File.POST("/object", integrateHdl.UploadFile)
			apiV1File.DELETE("/object", integrateHdl.DeleteFile)
		}

		apiV1UserNoAuth := apiV1.Group("/user")
		{
			apiV1UserNoAuth.POST("/login", userHdl.Login)
			apiV1UserNoAuth.POST("/sign-up", userHdl.SignIn)
		}
	}

	return r
}
