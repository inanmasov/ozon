package service

import (
	"ozon/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func StartService() {
	// Создаем новый экземпляр сервера Gin
	s := gin.Default()

	// Определяем маршруты для запросов на URL
	s.POST("/shorten", handlers.ShortenURLHandler)
	s.GET("/:shortURL", handlers.RedirectHandler)

	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	// Запускаем сервер на порту
	s.Run(":" + viper.GetString("port"))
}
