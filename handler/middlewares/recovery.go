package middlewares

import (
	"net/http"
	"runtime/debug"

	"chatbot-rag/base_common/constants"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logrus.Error(string(debug.Stack()))
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": constants.SomethingWentWrong,
				})
			}
		}()

		c.Next()
	}
}
