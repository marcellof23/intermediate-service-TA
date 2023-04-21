package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"github.com/intermediate-service-ta/helper"
)

func VerifyJWT(c *gin.Context) {
	if c.Request.Header["Token"] != nil {
		secretKey, err := helper.GetJWTSecretFromContext(c) // Get secret key if exist
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"message": err.Error(),
			})
			return
		}

		token, err := jwt.Parse(c.Request.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return "", errors.New("unauthorized")
			}
			return []byte(secretKey), nil
		})

		// parsing errors result
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "You're Unauthorized",
			})
			return
		}
		// if there's a token
		if token.Valid {
			c.Next()
			return
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "You're Unauthorized due to invalid token",
			})
			return
		}
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "You're Unauthorized due to No token in the header",
		})
		return
	}
}
