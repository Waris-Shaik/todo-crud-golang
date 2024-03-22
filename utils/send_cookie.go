package utils

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Waris-Shaik/todo-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func SendCookie(user *models.User, ctx *gin.Context) (*jwt.Token, error) {

	// Generate a jwt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"_id": user.ID,
		"exp": time.Now().Add(time.Minute * 30).Unix(), // Token expiration time
	})

	// Get JWT secret key
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	if len(secretKey) == 0 {
		return nil, fmt.Errorf("jwt secret key not found")
	}

	// Sign and get the encoded token as a string usig the jwt_secret
	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return nil, fmt.Errorf("failed to create token")
	}

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("token", tokenString, int((time.Minute * 30).Seconds()), "", "", false, true)
	fmt.Println("user.ID is", user.ID)

	return token, nil
}
