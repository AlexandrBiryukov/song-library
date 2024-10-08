// main.go
package main

import (
	"song-library/config"
	"song-library/database"
	"song-library/handlers"
	"song-library/logger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// @title Song Library API
// @version 1.0
// @description API для онлайн-библиотеки песен.

// @host localhost:8080
// @BasePath /
func main() {
	// Инициализация логгера
	logger.Init()

	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Log.Fatal(err)
	}

	// Инициализация базы данных
	err = database.InitDB(cfg)
	if err != nil {
		logger.Log.Fatal(err)
	}

	// Создание Echo instance
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Маршруты
	e.POST("/songs", handlers.AddSong(cfg))
	e.GET("/songs", handlers.GetSongs)
	e.GET("/songs/:id/text", handlers.GetSongText)
	e.PUT("/songs/:id", handlers.UpdateSong)
	e.DELETE("/songs/:id", handlers.DeleteSong)

	// Запуск сервера
	e.Logger.Fatal(e.Start(":8080"))
}
