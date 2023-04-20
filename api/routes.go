package api

import (
	"crypto/tls"
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

func newClient(endpoint, region, accessKeyID, secretKey string) *session.Session {
	cl := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	awsConfig := &aws.Config{
		Region:           aws.String("ap-southeast-1"),
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
	gcsSession := newClient("https://storage.googleapis.com", "auto", dep.Config().GoogleAccessKeyID, dep.Config().GoogleAccessKeySecret)
	doSession := newClient("https://sgp1.digitaloceanspaces.com", "us-east-1", dep.Config().DigitalOceanAccessKeyID, dep.Config().DigitalOceanAccessKeySecret)
	s3Session := newClient("", "ap-southeast-1", dep.Config().AmazonAccessKeyID, dep.Config().AmazonAccessKeySecret)

	sessionMap := make(map[string]*session.Session)
	sessionMap["gcs"] = gcsSession
	sessionMap["dos"] = doSession
	sessionMap["s3"] = s3Session

	return boot.Sess{
		SessionMap: sessionMap,
	}
}

func initClient(sess boot.Sess) boot.Client {
	var clientMap = make(map[string]*s3.S3)
	for _, v := range boot.Clients {
		clientMap[v] = s3.New(sess.SessionMap[v])
	}

	return boot.Client{
		ClientMap: clientMap,
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
		c.Set("db", dep.DB())
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
		apiV1.DELETE("/delete-object", integrateHdl.DeleteObject)
	}

	return r
}
