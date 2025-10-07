package controllers

import (
	"cloud_file_manager/src/dto"
	"cloud_file_manager/src/handlers"
	"cloud_file_manager/src/usecase"
	"cloud_file_manager/src/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AwsController struct {
	awsUsecase usecase.AwsUsecase
}

func NewAwsController(usecase usecase.AwsUsecase) AwsController {
	return AwsController{
		awsUsecase: usecase,
	}
}

func (ac *AwsController) CreateBucket(ctx *gin.Context) {
	claimsValue, exists := ctx.Get("claims")
	if !exists {
		response := handlers.Response {
			Message: "Não foi possível achar as informações do token",
		}
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	claims, ok := claimsValue.(jwt.MapClaims)
	if !ok {
			response := handlers.Response{
					Message: "Erro ao converter claims",
			}
			ctx.JSON(http.StatusInternalServerError, response)
			return
	}

	userId := int(claims["userId"].(float64))

	bucketName, err := utils.DecodeJson[dto.BucketNameDto](ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if bucketName.BucketName == "" {
		response := handlers.Response{
			Message: "É necessário o nome do bucket para sua criação",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	output, err := ac.awsUsecase.CreateBucket(userId, bucketName.BucketName)
	if err != nil {
		response := handlers.Response{
			Message: "Não foi possível criar o bucket, verifique o nome escolhido",
		}
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	ctx.JSON(http.StatusOK, output)
}

func (ac *AwsController) ListBuckets(ctx *gin.Context) {
	output, err := ac.awsUsecase.ListBuckets()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, output)
}	

func (ac *AwsController) ListBucketItems(ctx *gin.Context) {

	claimsValue, exists := ctx.Get("claims")
	if !exists {
		response := handlers.Response {
			Message: "Não foi possível achar as informações do token",
		}
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	claims, ok := claimsValue.(jwt.MapClaims)
	if !ok {
			response := handlers.Response{
					Message: "Erro ao converter claims",
			}
			ctx.JSON(http.StatusInternalServerError, response)
			return
	}

	userId := int(claims["userId"].(float64))

	output, err := ac.awsUsecase.ListBucketItems(userId)
	if err != nil {
		response := handlers.Response{
			Message: "Não foi possível listar os items do bucket",
		}
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	ctx.JSON(http.StatusOK, output)
}

func (ac *AwsController) GetObject(ctx *gin.Context) {
	claimsValue, exists := ctx.Get("claims")
	if !exists {
		response := handlers.Response {
			Message: "Não foi possível achar as informações do token",
		}
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	claims, ok := claimsValue.(jwt.MapClaims)
	if !ok {
			response := handlers.Response{
					Message: "Erro ao converter claims",
			}
			ctx.JSON(http.StatusInternalServerError, response)
			return
	}

	userId := int(claims["userId"].(float64))

	objectKey, err := utils.DecodeJson[dto.ObjectKeyDto](ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if objectKey.ObectKey == "" {
		response := handlers.Response{
			Message: "É necessário o caminho do arquivo",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	output, err := ac.awsUsecase.GetObject(userId, objectKey.ObectKey)
	if err != nil {
		response := handlers.Response{
			Message: "Não foi possível buscar o objeto, verifique o caminho do mesmo",
		}
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	ctx.JSON(http.StatusOK, output)
}

func (ac *AwsController) PutObject(ctx *gin.Context) {
	claimsValue, exists := ctx.Get("claims")
	if !exists {
		response := handlers.Response {
			Message: "Não foi possível achar as informações do token",
		}
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	claims, ok := claimsValue.(jwt.MapClaims)
	if !ok {
			response := handlers.Response{
					Message: "Erro ao converter claims",
			}
			ctx.JSON(http.StatusInternalServerError, response)
			return
	}

	userId := int(claims["userId"].(float64))

	objectKey, err := utils.DecodeJson[dto.ObjectKeyDto](ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if objectKey.ObectKey == "" {
		response := handlers.Response{
			Message: "É necessário o caminho do arquivo",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	output, err := ac.awsUsecase.PutObject(userId, objectKey.ObectKey)
	if err != nil {
		response := handlers.Response{
			Message: "Não foi possível buscar o objeto, verifique o caminho do arquivo",
		}
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	ctx.JSON(http.StatusOK, output)
}
