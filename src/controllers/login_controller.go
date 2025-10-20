package controllers

import (
	"cloud_file_manager/src/dto"
	"cloud_file_manager/src/handlers"
	"cloud_file_manager/src/usecase"
	"cloud_file_manager/src/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginController struct {
	userUsecase usecase.UserUsecase
}

func NewLoginController(usecase usecase.UserUsecase) LoginController {
	return LoginController{
		userUsecase: usecase,
	}
}

func (lc *LoginController) Login(ctx *gin.Context) {
	
	user, err := utils.DecodeJson[dto.UserLoginDto](ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Email == "" || user.Password == "" {
		response := handlers.Response{
			Message: "Informações inválidas",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	loginUser, err := lc.userUsecase.Login(*user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	if loginUser == nil {
		response := handlers.Response{
			Message: "Usuário não consta na base de dados",
		}
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	ctx.JSON(http.StatusOK, loginUser)
}