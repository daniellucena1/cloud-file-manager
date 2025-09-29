package controllers

import (
	"cloud_file_manager/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userController struct {

}

func NewUserController() userController {
	return userController{}
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