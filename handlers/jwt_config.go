package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyToken (ctx *gin.Context) {
	tokenString := ctx.GetHeader("Authorization")
	if tokenString == "" {
		response := Response{
			Message: "É necessário token de autorização",
		}
		ctx.JSON(http.StatusBadRequest, response)
		ctx.AbortWithStatus(http.StatusBadRequest)
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		response := Response{
			Message: "Não foi possível analisar o token",
		}
		ctx.JSON(http.StatusInternalServerError, response)
		ctx.AbortWithStatus(http.StatusBadRequest)
	}

	if !token.Valid {
		response := Response{
			Message: "Token inválido",
		}
		ctx.JSON(http.StatusUnauthorized, response)
		ctx.AbortWithStatus(http.StatusBadRequest)
	}

	ctx.Next()
}

func CreateToken(username string, password string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})
	
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}