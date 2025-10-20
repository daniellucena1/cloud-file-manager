package app

import (
	"cloud_file_manager/src/aws"
	"cloud_file_manager/src/config"
	"cloud_file_manager/src/controllers"
	"cloud_file_manager/src/database"
	"cloud_file_manager/src/repository"
	"cloud_file_manager/src/routes"
	"cloud_file_manager/src/usecase"
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
	presigner := s3.NewPresignClient(client)

	server := gin.Default()

	server.Use(config.CORSMiddleware())

	UserRepository := repository.NewUserRepository(dbConection)
	AwsService := aws.NewAwsService(client, presigner)
	AwsUsecase := usecase.NewAwsUsecase(AwsService)
	UserUsecase := usecase.NewUserUseCase(UserRepository, AwsService)
	UserController := controllers.NewUserController(UserUsecase)
	LoginController := controllers.NewLoginController(UserUsecase)
	AwsController := controllers.NewAwsController(AwsUsecase)

	routes.SetupRoutes(server, UserController, LoginController, AwsController)

	server.Run(":8000")

	return nil
}
