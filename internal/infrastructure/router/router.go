package router

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"

	adapterhandler "github.com/peconote/peconote/internal/adapter/handler"
	adapterrepo "github.com/peconote/peconote/internal/adapter/repository"
	"github.com/peconote/peconote/internal/infrastructure/persistence"
	"github.com/peconote/peconote/internal/interfaces/controller"
	"github.com/peconote/peconote/internal/usecase"
)

func NewRouter(gormDB *gorm.DB, sqlxDB *sqlx.DB) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery(), jsonLogger())

	userRepo := persistence.NewUserRepository(gormDB)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userController := controller.NewUserController(userUsecase)

	r.GET("/users", userController.GetUsers)
	r.POST("/users", userController.CreateUser)

	memoRepo := adapterrepo.NewMemoRepository(sqlxDB)
	memoUsecase := usecase.NewMemoUsecase(memoRepo)
	memoHandler := adapterhandler.NewMemoHandler(memoUsecase)

	r.POST("/api/memos", memoHandler.CreateMemo)

	return r
}

func jsonLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		m := map[string]interface{}{
			"method":     param.Method,
			"path":       param.Path,
			"status":     param.StatusCode,
			"latency_ms": param.Latency.Milliseconds(),
		}
		if v := param.Request.Context().Value("trace_id"); v != nil {
			if s, ok := v.(string); ok {
				m["trace_id"] = s
			}
		}
		b, _ := json.Marshal(m)
		return string(b) + "\n"
	})
}
