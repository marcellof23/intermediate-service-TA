package api

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/intermediate-service-ta/boot"
)

func InitRoutes(dep *boot.Dependencies) *gin.Engine {

	r := gin.New()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length", "Authorization", "Content-type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	return r
}
