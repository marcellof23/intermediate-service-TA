package user

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	"github.com/intermediate-service-ta/helper"
	"github.com/intermediate-service-ta/internal/model"
	"github.com/intermediate-service-ta/internal/repository"
)

type UserHandler struct {
	userRepo repository.UserRepository
}

func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

type M map[string]interface{}

func generateJWT(c *gin.Context, userID int64, username string) (string, error) {
	secretKey, err := helper.GetJWTSecretFromContext(c) // Get secret key if exist
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userID
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 30).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))

	if err != nil {
		fmt.Errorf("something went wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func (hdl *UserHandler) SignIn(c *gin.Context) {
	var creds model.Credentials
	// Get the JSON body and decode into credentials
	err := c.ShouldBindJSON(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		c.AbortWithStatusJSON(400, err)
		return
	}

	var user model.User
	user = model.User{
		Username: creds.Username,
		Password: helper.GetHash([]byte(creds.Password)),
		Role:     model.Normal,
	}

	u, err := hdl.userRepo.FindByUsername(c, user.Username)
	if err == nil && u.ID != 0 {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "username already taken",
		})
		return
	}

	res, err := hdl.userRepo.Create(c, &user)
	if err != nil {
		c.AbortWithStatusJSON(500, err)
		return
	}
	user.GroupID = user.ID

	_, err = hdl.userRepo.Update(c, &user)
	if err != nil {
		c.AbortWithStatusJSON(500, err)
		return
	}

	tokenString, err := generateJWT(c, res.ID, res.Username)
	if err != nil {
		c.AbortWithStatusJSON(500, err)
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    res,
		"token":   tokenString,
	})
}

func (hdl *UserHandler) Login(c *gin.Context) {
	var creds model.Credentials
	// Get the JSON body and decode into credentials
	err := c.ShouldBindJSON(&creds)
	if err != nil {
		c.AbortWithStatusJSON(400, err)
		return
	}

	user, err := hdl.userRepo.FindByUsername(c, creds.Username)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "username or password is invalid",
		})
		return
	}

	userPass := []byte(creds.Password)
	dbPass := []byte(user.Password)

	passErr := bcrypt.CompareHashAndPassword(dbPass, userPass)
	if passErr != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "username or password is invalid",
		})
		return
	}
	tokenString, err := generateJWT(c, user.ID, user.Username)
	if err != nil {
		c.AbortWithStatusJSON(500, err)
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    user,
		"token":   tokenString,
	})
}

func (hdl *UserHandler) GetClients(c *gin.Context) {
	clients, err := helper.GetClientsFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"message": "no config clients in intermediate service",
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    clients,
	})
}
