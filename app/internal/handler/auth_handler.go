package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Ozielmarcos/mytodolist/app/internal/model"
	"github.com/Ozielmarcos/mytodolist/app/internal/repository"
	"github.com/Ozielmarcos/mytodolist/app/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RegisterHandler(c *gin.Context) {
	var user model.User
	c.BindJSON(&user)

	hash, _ := service.HashPassword(user.Password)
	user.Password = hash

	err := repository.CreateUser(user)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Erro ao criar usuário!"})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"name": user.Name, "email": user.Email})
}

func Login(c *gin.Context) {
	var input model.User

	c.BindJSON(&input)
	fmt.Printf("usuário %v", &input)

	user, err := repository.GetUserByEmail(input.Email)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Usuário não encontrado!"})
		fmt.Print(err)
		return
	}

	if !service.CheckPassword(user.Password, input.Password) {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Senha incorreta!"})
		return
	}

	token, _ := service.GenerateToken(user)
	refreshToken, _ := service.GenerateRefreshToken(user)

	c.IndentedJSON(http.StatusOK, gin.H{"user_id": user.ID, "token": token, "refresh_token": refreshToken})
}

func RefreshTokenHandler(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token requerido"})
		return
	}

	token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Refresh token inválido ou expirado!"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Claims inválidos!"})
		return
	}

	userId := fmt.Sprintf("%v", claims["user_id"])

	user := model.User{ID: userId}
	newAccessToken, _ := service.GenerateRefreshToken(user)

	c.JSON(http.StatusOK, gin.H{"token": newAccessToken})
}
