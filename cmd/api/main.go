package main

import (
	"testGolang/internal/config"
	"testGolang/internal/handler"
	"testGolang/internal/router"
	"testGolang/internal/service"
)

func main() {
	// Подключение Config
	cfg := config.MustLoad()

	// Сборка роутера с внедрёнными сервисами
	userService := service.NewUserService()
	permissionsService := service.NewPermissionsService()
	vectorService := service.NewVectorMemoryService()
	chatHandler := handler.NewChatHandler(userService, permissionsService, vectorService)

	router := router.SetupRouter(chatHandler, cfg)
	router.Run(cfg.HTTPServer.Addr)

	//TODO: init logger: slog
	//TODO init storage
}
