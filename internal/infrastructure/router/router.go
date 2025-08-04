package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/peconote/peconote/internal/infrastructure/persistence"
	"github.com/peconote/peconote/internal/interfaces/controller"
	"github.com/peconote/peconote/internal/usecase"
)

func NewRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	userRepo := persistence.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userController := controller.NewUserController(userUsecase)

	r.GET("/users", userController.GetUsers)
	r.POST("/users", userController.CreateUser)

	return r
}
