package controllers

import (
	"cloud_file_manager/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userController struct {
	userUsecase usecase.UserUsecase
}

func NewUserController(usecase usecase.UserUsecase) userController {
	return userController{
		userUsecase: usecase,
	}
}

func (p *userController) GetUsers(ctx *gin.Context) {
	
	users, err := p.userUsecase.GetUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusOK, users)
}