package router

import (
	"testGolang/internal/config"
	"testGolang/internal/handler"
	"testGolang/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(chatHandler *handler.ChatHandler, cfg *config.Config) *gin.Engine {
	// Не используем gin.Default() чтобы контролировать middleware самостоятельно.
	r := gin.New()

	// Basic middlewares: logger + recovery. Можно заменить кастомным логгером.
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.TimeoutMiddleware(cfg.HTTPServer.Timeout))

	v1 := r.Group("/api/v1")
	{
		v1.GET("/chat/:id/summary", chatHandler.GetChatSummary)

	}

	return r
}
