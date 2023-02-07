package do

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

type DOSHandler struct {
}

func NewDOSHandler() *DOSHandler {
	return &DOSHandler{}
}

func (hdl *DOSHandler) ListDOSBuckets(c *gin.Context) {
	dosSession, ok := c.MustGet("doSession").(*session.Session)
	if !ok {
		fmt.Println("Failed to get digital ocean session")
		return
	}

	s3Client := s3.New(dosSession)

	spaces, err := s3Client.ListBuckets(nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, b := range spaces.Buckets {
		fmt.Println(aws.StringValue(b.Name))
	}
}
