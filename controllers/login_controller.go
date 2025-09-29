package controllers

import (
	"cloud_file_manager/dto"
	"cloud_file_manager/handlers"
	"cloud_file_manager/usecase"
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
	var user dto.UserLoginDto
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	if user.Email == "" || user.Password == "" {
		response := handlers.Response{
			Message: "Informações inválidas",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	loginUser, err := lc.userUsecase.Login(user)
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