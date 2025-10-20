package controllers

import (
	"cloud_file_manager/src/handlers"
	"cloud_file_manager/src/models"
	"cloud_file_manager/src/usecase"
	"cloud_file_manager/src/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUsecase usecase.UserUsecase
}

func NewUserController(usecase usecase.UserUsecase) UserController {
	return UserController{
		userUsecase: usecase,
	}
}

func (u *UserController) GetUsers(ctx *gin.Context) {
	
	users, err := u.userUsecase.GetUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusOK, users)
}

func (u *UserController) CreateUser(ctx *gin.Context) {

	user, err := utils.DecodeJson[models.User](ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	insertedUser, err := u.userUsecase.CreateUser(*user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, insertedUser)
}

func (u *UserController) GetUserById(ctx *gin.Context) {
	
	id := ctx.Param("id")
	if id == "" {
		response := handlers.Response{
			Message: "Id do produto não pode ser nulo",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	userId, err := strconv.Atoi(id)
	if err != nil {
		response := handlers.Response{
			Message: "Id do produto precisa ser um número",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	user, err := u.userUsecase.GetUserById(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	
	if user == nil {
		response := handlers.Response{
			Message: "Usuário não foi encontrado na base de dados",
		}
		ctx.JSON(http.StatusNotFound, response)
		return
	}
	ctx.JSON(http.StatusOK, user)
}