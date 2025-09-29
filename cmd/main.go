package main

import (
	"cloud_file_manager/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
    server := gin.Default()

    UserController := controllers.NewUserController()

    server.GET("/ping", func(ctx *gin.Context) {
        ctx.JSON(200, gin.H{
            "message": "PONG",
        })
    })

    server.GET("/users", UserController.GetUsers)

    server.Run(":8000")
}