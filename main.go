package main

import (
	"log"
	"os"
	"time"

	"chatbot-rag/base_common/constants"
	"chatbot-rag/base_common/milvus"
	"chatbot-rag/handler"
	"chatbot-rag/handler/middlewares"
	"chatbot-rag/internal/domains/user"
	"chatbot-rag/internal/infrastructure/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	sdkMilvus "github.com/milvus-io/milvus-sdk-go/milvus"
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
	conn, err := milvus.NewMilvusClient(viper.GetString("MILVUS_HOST"), viper.GetString("MILVUS_PORT"), constants.DefaultTimeout)
	if err != nil {
		logrus.Errorf(constants.FormatTaskErr, "NewMilvusClient", err)
		panic(err)
	}

	isHasCollection, err := conn.HasCollection(viper.GetString("MILVUS_COLLECTION"))
	if err != nil {
		logrus.Errorf(constants.FormatTaskErr, "HasCollection", err)
		panic(err)
	}

	if !isHasCollection {
		err = conn.CreateCollection(viper.GetString("MILVUS_COLLECTION"), 1536, 1024, sdkMilvus.IP)
		if err != nil {
			panic(err)
		}

		err = conn.CreateIndex(viper.GetString("MILVUS_COLLECTION"), 1024, sdkMilvus.IVFFLAT)
		if err != nil {
			panic(err)
		}
	}

	// open ai
	openAI := openai.NewClient(viper.GetString("OPENAI_API_KEY"))

	// repository
	embeddingRepo := repository.NewEmbeddingRepository(db)

	userService := user.NewUserService(conn, openAI, embeddingRepo)

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
