package main

import (
	"chatbot-rag/handler"
	"chatbot-rag/handler/middlewares"
	"chatbot-rag/internal/domains/user"
	"github.com/spf13/viper"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func main() {
	viper.AddConfigPath("build")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return
	}

	userService := user.NewUserService()

	h := handler.NewHandler(userService)

	r := gin.New()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(middlewares.Recover())

	api := r.Group("v1/api")
	{
		api.Use(middlewares.RequestID())
		api.POST("/index", h.Index)
		api.POST("/chat", h.Chat)
	}

	err = r.Run(":" + viper.GetString("PORT"))
	if err != nil {
		return
	}
}
