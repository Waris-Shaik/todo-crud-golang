package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/Waris-Shaik/todo-backend/initializers"
	"github.com/Waris-Shaik/todo-backend/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func ValidateUserData(user *models.User) error {
	if user.Name == "" || user.UserName == "" || user.Email == "" || user.Password == "" {
		return fmt.Errorf("please fill all required fields")
	}
	return nil
}

func ValidatePassword(user *models.User) error {

	const (
		minPasswordLength = 6
		maxPsswordLength  = 13
	)

	if len(user.Password) < minPasswordLength || len(user.Password) > maxPsswordLength {
		if len(user.Password) < minPasswordLength {

			return fmt.Errorf("password should contain at least 6 characters")
		} else {

			return fmt.Errorf("password must not be greater than 13 characters")
		}
	}

	return nil
}

func CheckExistingUser(email string) error {
	// Retreive hthe user from the database
	var existingUser models.User
	if result := initializers.DB.Where("email = ?", email).First(&existingUser); result.Error == nil {
		return fmt.Errorf("user already exists please login")
	}

	return nil
}

func CreateUser(user *models.User) error {
	result := initializers.DB.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GenerateToken(user *models.User) (string, error) {
	// Generate a jwt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"_id": user.ID,
		"exp": time.Now().Add(time.Minute * 30).Unix(), // Token expiration time
	})

	// Get JWT secret key
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	if len(secretKey) == 0 {
		return "", fmt.Errorf("jwt secret key not found")
	}

	// Sign and get the encoded token as a string usig the jwt_secret
	return token.SignedString(secretKey)

}

func IsPasswordMatches(userPassword *string, existingUserPassword *string) error {

	err := bcrypt.CompareHashAndPassword([]byte(*existingUserPassword), []byte(*userPassword))
	if err != nil {
		return fmt.Errorf("password does not matches")
	}

	return nil

}
