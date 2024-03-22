package controllers

import (
	"fmt"
	"net/http"

	"github.com/Waris-Shaik/todo-backend/initializers"
	"github.com/Waris-Shaik/todo-backend/models"
	"github.com/gin-gonic/gin"
)

func CreateTodo(ctx *gin.Context) {

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

	// Parse the request body to get post data
	var todo models.Todo
	if err := ctx.Bind(&todo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Check if requied fields are empty
	if todo.Title == "" || todo.Description == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "please fill all required fields",
		})
		return
	}

	// Set the userID for the todo
	todo.UserID = userID
	userObj := models.UserLite{
		ID:       user.(models.User).ID,
		UserName: user.(models.User).UserName,
		Email:    user.(models.User).Email,
	}
	todo.User = userObj

	// Create the todo in the database
	result := initializers.DB.Create(&todo)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to create todo",
		})
		return
	}

	// Return the response
	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Todo Successfully Created",
		"todo":    todo,
	})

}

func GetTodos(ctx *gin.Context) {
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

	// Retreive the todos
	var todos []models.Todo
	result := initializers.DB.Preload("User").Where("user_id = ?", userID).Find(&todos)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to fetch todos",
			"alert":   result.Error.Error(),
		})
		return
	}

	// Return the reponse
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"todos":   todos,
	})
}

func GetSingleTodo(ctx *gin.Context) {

	// Get the todoID from URL paramter
	todoID := ctx.Param("id")

	if todoID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "todo ID is required",
		})
		return
	}

	// Retreive todo from the database
	var todo models.Todo
	result := initializers.DB.Preload("User").First(&todo, todoID)
	if result.Error != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "todo not found",
		})
		return
	}

	// Return the retreived todo
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"todo":    todo,
	})
}

func UpdateTodo(ctx *gin.Context) {

	// Get todoIF from URL parameter
	todoID := ctx.Param("id")

	if todoID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "todo ID is required",
		})
		return
	}

	// Retreive todo from the database
	var todo models.Todo
	result := initializers.DB.Preload("User").First(&todo, todoID)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "todo not found",
		})
		return
	}

	// toggle update
	result = initializers.DB.Model(&todo).Update("completed", !todo.Completed)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "failed to update todo",
		})
		return
	}

	// Return the updated post
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "todo successfully updated",
		"todo":    todo,
	})

}

func EditTodo(ctx *gin.Context) {

	// Get the todoID from URL parameter
	todoID := ctx.Param("id")
	if todoID == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "todo ID is required",
		})
		return
	}

	// Retreive todo from the database
	var originalTodo models.Todo
	result := initializers.DB.Preload("User").First(&originalTodo, todoID)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "todo not found",
		})
		return
	}

	// Parse the request body to get edited data
	var editedTodo models.Todo
	if err := ctx.Bind(&editedTodo); err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request body",
		})
		return
	}

	// Check if any changes made
	if editedTodo.Title == "" && editedTodo.Description == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "no changes were made",
		})
		return
	}

	result = initializers.DB.Model(&originalTodo).Updates(models.Todo{Title: editedTodo.Title, Description: editedTodo.Description})
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "failed to edit todo",
		})
		return
	}

	// Return the updated todo in response
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "todo edited successfully",
		"todo":    originalTodo,
	})
}

func DeleteTodo(ctx *gin.Context) {

	// Get todoID from URL parameter

	todoID := ctx.Param("id")

	if todoID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "token ID is required",
		})
		return
	}

	// Retreive todo from the database
	var todo models.Todo
	result := initializers.DB.Preload("User").First(&todo, todoID)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "todo not found",
		})
		return
	}

	// Delete todo in the database
	result = initializers.DB.Where("ID = ?", todoID).Delete(&todo)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "failed to delete todo",
		})
		return
	}

	// Return the deleted tod in  response
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Todo deleted successfully",
		"todo":    todo,
	})

}
