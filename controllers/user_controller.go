package controllers

import (
	"cloud_file_manager/models"
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
	users := []models.User{
		{
			ID: 1,
			Name: "Daniel Torres",
			Email: "daniel@exemplo.com",
			Password: "qualquercoisa123",
		},
	}

	ctx.JSON(http.StatusOK, users)
}