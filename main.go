package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"chatbot-rag/base_common/constants"
	"chatbot-rag/handler"
	"chatbot-rag/handler/middlewares"
	"chatbot-rag/internal/domains/user"
	"chatbot-rag/internal/infrastructure/repository"
	"github.com/milvus-io/milvus-sdk-go/v2/client"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func initDatabase() (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)

	db, err := gorm.Open(mysql.Open(viper.GetString("DB_SOURCE")), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	return db, nil
}

func initMilvus(host, port string, timeout time.Duration) (client.Client, error) {
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	milvusAddress := fmt.Sprintf("%s:%s", host, port)

	config := client.Config{
		Address:  milvusAddress,
		Username: viper.GetString(constants.MilvusUserName),
		Password: viper.GetString(constants.MilvusPassword),
	}

	c, err := client.NewClient(ctx, config)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		err = ctx.Err()
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("connection timeout")
		} else if errors.Is(err, context.Canceled) {
			return nil, fmt.Errorf("context canceled")
		} else {
			return nil, err
		}

	default:
		milvusState, err := c.CheckHealth(ctx)
		if err != nil {
			return nil, err
		}
		if !milvusState.IsHealthy {
			logrus.Error(fmt.Sprintf("milvus v2 at %v is unhealthy", milvusAddress))
		}
	}

	return c, nil
}

func main() {
	viper.AddConfigPath("build")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return
	}

	db, err := initDatabase()
	if err != nil {
		logrus.Fatal("failed to open database:", err)
		return
	}

	// milvus
	milvusClient, err := initMilvus(viper.GetString(constants.MilvusHost), viper.GetString(constants.MilvusPort),
		30*time.Second)
	if err != nil {
		logrus.Fatal("failed to open milvus:", err)
		return
	}

	// open ai
	openAI := openai.NewClient(viper.GetString("OPENAI_API_KEY"))

	// repository
	milvusRepo := repository.NewMilvusRepository(milvusClient, 30*time.Second)
	embeddingRepo := repository.NewEmbeddingRepository(db)

	// init milvus collection and bucket
	err = milvusRepo.InitCollection(viper.GetString(constants.MilvusCollection))
	if err != nil {
		logrus.Errorf(constants.FormatTaskErr, "milvusRepo.InitCollection", err)
		return
	}

	userService := user.NewUserService(openAI, milvusRepo, embeddingRepo)

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

	r.Static("/static", "./static")

	routes := map[string]string{
		"": "chat.html",
	}

	for path, file := range routes {
		r.GET(path, func(c *gin.Context) {
			c.File("./static/" + file)
		})
	}

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
