package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Waris-Shaik/todo-backend/controllers"
	"github.com/Waris-Shaik/todo-backend/initializers"
	"github.com/Waris-Shaik/todo-backend/middlewares"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func main() {

	// PORT
	PORT := os.Getenv("PORT")

	// PORT error
	if PORT == "" {
		log.Fatal("Error PORT is not defined in .env File")
	}

	// router
	router := gin.Default()

	// routes
	router.POST("/api/v1/users/signup", controllers.SignUp)
	router.POST("/api/v1/users/login", controllers.Login)
	router.GET("/api/v1/users/logout", controllers.Logout)
	router.GET("/api/v1/users/me", middlewares.IsAuthenticated, controllers.Me)
	router.GET("/api/v1/users/all", middlewares.IsAuthenticated, controllers.GetUsers)
	router.PATCH("/api/v1/users/updatemyprofile", middlewares.IsAuthenticated, controllers.UpdateUser)
	router.POST("/api/v1/todos/new", middlewares.IsAuthenticated, controllers.CreateTodo)
	router.GET("/api/v1/todos/my", middlewares.IsAuthenticated, controllers.GetTodos)
	router.GET("/api/v1/todos/:id", middlewares.IsAuthenticated, controllers.GetSingleTodo)
	router.PATCH("/api/v1/todos/:id", middlewares.IsAuthenticated, controllers.UpdateTodo)
	router.PUT("/api/v1/todos/:id", middlewares.IsAuthenticated, controllers.EditTodo)
	router.DELETE("/api/v1/todos/:id", middlewares.IsAuthenticated, controllers.DeleteTodo)

	// Server listening
	fmt.Println("Server is listening on PORT:", PORT, "⚡⚡⚡")

	// Server error
	if err := router.Run(":" + PORT); err != nil {
		log.Fatal("Failed to connect to server")
	}

}
