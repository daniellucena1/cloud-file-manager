package app

import (
	"cloud_file_manager/aws"
	"cloud_file_manager/config"
	"cloud_file_manager/controllers"
	"cloud_file_manager/database"
	"cloud_file_manager/repository"
	"cloud_file_manager/routes"
	"cloud_file_manager/usecase"
	"context"
	"log"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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

	cfg, err := awsConfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(cfg)

	server := gin.Default()
	
	UserRepository := repository.NewUserRepository(dbConection)
	AwsService := aws.NewAwsService(client)
	AwsUsecase := usecase.NewAwsUsecase(AwsService)
	UserUsecase := usecase.NewUserUseCase(UserRepository)
	UserController := controllers.NewUserController(UserUsecase)
	LoginController := controllers.NewLoginController(UserUsecase)
	AwsController := controllers.NewAwsController(AwsUsecase)

	routes.SetupRoutes(server, UserController, LoginController, AwsController)

	server.Run(":8000")

	return nil
}