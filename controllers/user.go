package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Waris-Shaik/todo-backend/initializers"
	"github.com/Waris-Shaik/todo-backend/models"
	"github.com/Waris-Shaik/todo-backend/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type SafeUser struct {
	ID        interface{} `json:"_id"`
	Name      string      `json:"name"`
	UserName  string      `json:"username"`
	Email     string      `json:"email"`
	CreatedAt time.Time   `json:"created_at"`
	// Exclude Password field
}

func SignUp(ctx *gin.Context) {

	// Parse the request body to get user data
	var user models.User
	if err := ctx.Bind(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Check if required fields are empty
	if err := utils.ValidateUserData(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Check password criteria matches
	if err := utils.ValidatePassword(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Check if user exists
	if err := utils.CheckExistingUser(user.Email); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Set the hashedPassword to user object
	user.Password = string(hashedPassword)

	// Store the user in database
	if err := utils.CreateUser(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Generate the token
	token, err := utils.GenerateToken(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success":  false,
			"messsage": err.Error(),
		})
		return
	}

	// Send the cookie 🍪
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("token", token, int((time.Minute * 30).Seconds()), "", "", false, true)

	// Return the created user
	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User successfully created",
		"user":    user,
	})

}

func Login(ctx *gin.Context) {

	// Parse the request body to get user data
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := ctx.Bind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Check if required fields are empty
	if body.Email == "" || body.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "please fill all required fileds",
		})
		return
	}

	// Check if user not exists
	var user models.User

	result := initializers.DB.Where("email = ?", body.Email).First(&user)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "user not found please register",
		})
		return
	}

	// Verify password
	if err := utils.IsPasswordMatches(&body.Password, &user.Password); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Generate a JWT token
	token, err := utils.GenerateToken(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	fmt.Println("token is:", token)

	// Set the cookie
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("token", token, int((time.Minute * 30).Seconds()), "", "", false, true)

	message := fmt.Sprintf("Welcome back %v", user.Name)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": message,
	})

}

func Logout(ctx *gin.Context) {

	// Remove the token in cookies
	ctx.SetCookie("token", "", -1, "", "", false, true)

	// Return the response
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Successfully logged out",
	})
}

func Me(ctx *gin.Context) {
	// Retreive the user from the request context
	user, exists := ctx.Get("user")
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "user data not found in context",
		})
		return
	}

	// Assert it
	userData, ok := user.(models.User)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retreive user data from the context",
		})
	}

	// Create an instance of SafeUser struct
	safeUserData := SafeUser{
		ID:        userData.ID,
		Name:      userData.Name,
		UserName:  userData.UserName,
		Email:     userData.Email,
		CreatedAt: userData.CreatedAt,
	}

	// Return the response
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Welcome Back " + safeUserData.Name,
		"user":    safeUserData,
	})
}

func GetUsers(ctx *gin.Context) {

	// Retreive users from the database
	var users []models.User
	result := initializers.DB.Find(&users)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": result.Error.Error(),
		})
		return
	}

	// Return the users in response
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"users":   users,
	})
}

func UpdateUser(ctx *gin.Context) {

	// Extracr user information from context
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "user not found",
		})
		return
	}

	// Parser user ID
	userID := user.(models.User).ID

	// Retreive user from the database
	var existingUser models.User
	result := initializers.DB.First(&existingUser, userID)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "user not found",
		})
		return
	}

	// Parse the request body to get updated user
	var updateUser models.User
	if err := ctx.Bind(&updateUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "failed to parse request body",
		})
		return
	}

	// Check if required fields are not empty
	if updateUser.Name == "" && updateUser.UserName == "" && updateUser.Email == "" && updateUser.Password == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "no changes were made",
		})
		return
	}

	var showPassword string
	// Check criteria meets or not
	if updateUser.Password != "" && len(updateUser.Password) > 0 {
		// Check password criteria matches
		if err := utils.ValidatePassword(&updateUser); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		showPassword = updateUser.Password

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateUser.Password), 10)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		// Set the hashedPassword to updateUser.password
		updateUser.Password = string(hashedPassword)

	}

	// Update the user object
	result = initializers.DB.Model(&existingUser).Updates(models.User{Name: updateUser.Name, UserName: updateUser.UserName, Email: updateUser.Email, Password: updateUser.Password, UpdatedAt: time.Now()})
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "failed to update user",
		})
		return
	}

	safeUserData := struct {
		ID        uint      `json:"_id"`
		Name      string    `json:"name"`
		UserName  string    `json:"username"`
		Email     string    `json:"email"`
		Password  any       `json:"password"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}{
		ID:        existingUser.ID,
		Name:      existingUser.Name,
		UserName:  existingUser.UserName,
		Email:     existingUser.Email,
		Password:  showPassword,
		CreatedAt: existingUser.CreatedAt,
		UpdatedAt: existingUser.UpdatedAt,
	}

	// Return the updated user in response
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User successfully updated",
		"user":    safeUserData,
	})

}
