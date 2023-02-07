package api

import (
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/intermediate-service-ta/boot"
	dohandler "github.com/intermediate-service-ta/do"
	gcshandler "github.com/intermediate-service-ta/gcs"
	s3handler "github.com/intermediate-service-ta/s3"
)

type Sess struct {
	GCSSession *session.Session
	DOSSession *session.Session
	S3Session  *session.Session
}

func initSession(dep *boot.Dependencies) Sess {
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

	return Sess{
		GCSSession: gcsSession,
		DOSSession: doSession,
		S3Session:  s3Session,
	}
}

func InitRoutes(dep *boot.Dependencies) *gin.Engine {

	// init Handler
	gcsHdl := gcshandler.NewGCSHandler()
	dosHdl := dohandler.NewDOSHandler()
	s3Hdl := s3handler.NewS3Handler()

	// init blank engine
	r := gin.New()

	// init session
	sess := initSession(dep)

	// attach session to context
	r.Use(func(c *gin.Context) {
		c.Set("gcsSession", sess.GCSSession)
		c.Set("doSession", sess.DOSSession)
		c.Set("s3Session", sess.S3Session)
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
		apiV1.GET("/gcs/list-bucket", gcsHdl.ListGCSBuckets)
		apiV1.GET("/do/list-bucket", dosHdl.ListDOSBuckets)
		apiV1.GET("/s3/list-bucket", s3Hdl.ListS3Objects)
	}

	return r
}
