package middlewares

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Waris-Shaik/todo-backend/initializers"
	"github.com/Waris-Shaik/todo-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func IsAuthenticated(ctx *gin.Context) {

	// Get the JWT token from cookies
	tokenString, err := ctx.Cookie("token")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "please login",
		})
		return
	}

	// fmt.Println("Token String is:", tokenString)

	// Parse and validate JWT token

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		// Ensure the signing method is valid
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["Authorization"])
		}

		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil || !token.Valid {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "invalid token",
		})
		return
	}

	// Extract user ID from claims and fetch it from the user data
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "invalid token claims",
		})
		return
	}

	userID, ok := claims["_id"]
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "invalid user id in token",
		})
		return
	}

	// Retreive user from the database
	var user models.User
	result := initializers.DB.First(&user, userID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "user not found",
		})
		return
	}

	// Attach the user information to the request context
	ctx.Set("user", user)

	// Proceed to the next middleware or route handler
	ctx.Next()

}
