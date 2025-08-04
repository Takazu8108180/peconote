package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/peconote/peconote/internal/domain/model"
	"github.com/peconote/peconote/internal/usecase"
)

type UserController struct {
	usecase usecase.UserUsecase
}

func NewUserController(u usecase.UserUsecase) *UserController {
	return &UserController{usecase: u}
}

func (c *UserController) GetUsers(ctx *gin.Context) {
	users, err := c.usecase.GetUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.usecase.CreateUser(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, user)
}
