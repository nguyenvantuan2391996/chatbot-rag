package middlewares

import (
	"chatbot-rag/base_common/constants"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(constants.RequestIDField, uuid.NewString())
		ctx.Next()
	}
}
