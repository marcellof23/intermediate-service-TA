package api

import (
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/intermediate-service-ta/boot"
	integratehandler "github.com/intermediate-service-ta/integrator-storage"
)

func initSession(dep *boot.Dependencies) boot.Sess {
	gcsSession := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("auto"),
		Endpoint:    aws.String("https://storage.googleapis.com"),
		Credentials: credentials.NewStaticCredentials(dep.Config().GoogleAccessKeyID, dep.Config().GoogleAccessKeySecret, ""),
	}))

	doSession := session.Must(session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(dep.Config().DigitalOceanAccessKeyID, dep.Config().DigitalOceanAccessKeySecret, ""),
		Endpoint:         aws.String("https://sgp1.digitaloceanspaces.com"),
		Region:           aws.String("us-east-1"),
		S3ForcePathStyle: aws.Bool(false), // // Configures to use subdomain/virtual calling format. Depending on your version, alternatively use o.UsePathStyle = false
	}))

	s3Session := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("ap-southeast-1"),
		Credentials: credentials.NewStaticCredentials(dep.Config().AmazonAccessKeyID, dep.Config().AmazonAccessKeySecret, ""),
	}))

	return boot.Sess{
		GCSSession: gcsSession,
		DOSSession: doSession,
		S3Session:  s3Session,
	}
}

func initClient(sess boot.Sess) boot.Client {
	gcsClient := s3.New(sess.GCSSession)
	dosClient := s3.New(sess.DOSSession)
	s3Client := s3.New(sess.S3Session)

	return boot.Client{
		GCSClient: gcsClient,
		DOSClient: dosClient,
		S3Client:  s3Client,
	}

}

func InitRoutes(dep *boot.Dependencies) *gin.Engine {

	// init Handler
	integrateHdl := integratehandler.NewIntegratorHandler()

	// init blank engine
	r := gin.New()

	// init session
	sess := initSession(dep)
	client := initClient(sess)

	// attach session to context
	r.Use(func(c *gin.Context) {
		c.Set("vdfsSession", sess)
		c.Set("vdfsClient", client)
	})

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
		apiV1.GET("/list-bucket", integrateHdl.ListBuckets)
		apiV1.POST("/upload-object", integrateHdl.UploadObject)
	}

	return r
}
