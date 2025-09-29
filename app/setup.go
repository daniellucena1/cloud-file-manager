package app

import (
	"cloud_file_manager/config"
	"cloud_file_manager/controllers"
	"cloud_file_manager/database"
	"cloud_file_manager/repository"
	"cloud_file_manager/usecase"

	"github.com/gin-gonic/gin"
)

func SetupAndRunApp() error {
	err := config.LoadENV()
	if err != nil {
		return err
	}

	dbConection, err := database.ConnectDB()
	if err != nil {
		return err
	}

	server := gin.Default()
	
	UserRepository := repository.NewUserRepository(dbConection)

	UserUsecase := usecase.NewUserUseCase(UserRepository)

	UserController := controllers.NewUserController(UserUsecase)

	server.GET("/ping", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
					"message": "PONG",
			})
	})

	server.GET("/users", UserController.GetUsers)

	server.Run(":8000")

	return nil
}