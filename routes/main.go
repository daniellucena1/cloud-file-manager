package routes

import (
	"cloud_file_manager/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes (server *gin.Engine, UserController controllers.UserController) {

	// PING
	server.GET("/ping", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
					"message": "PONG",
			})
	})

	// User routes
	users := server.Group("/users")
	users.GET("/", UserController.GetUsers)
	users.GET("/:id", UserController.GetUserById)
	users.POST("/", UserController.CreateUser)
}