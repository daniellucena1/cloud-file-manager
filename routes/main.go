package routes

import (
	"cloud_file_manager/controllers"
	"cloud_file_manager/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes (
	server *gin.Engine, 
	UserController controllers.UserController, 
	LoginController controllers.LoginController,
	AwsController controllers.AwsController,
) {

	// PING
	server.GET("/ping", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
					"message": "PONG",
			})
	})

	// User routes
	users := server.Group("/users")
	users.GET("/", UserController.GetUsers)
	users.GET("/:id", handlers.VerifyToken, UserController.GetUserById)
	users.POST("/", UserController.CreateUser)

	// Login routes
	login := server.Group("/login")
	login.POST("/", LoginController.Login)

	// Aws routes
	aws := server.Group("/aws")
	aws.POST("/bucket", handlers.VerifyToken, AwsController.CreateBucket)
	aws.GET("/bucket", handlers.VerifyToken, AwsController.ListBuckets)
}