package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"trellode-go/internal/utils/logging"
)

func LoggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("uuid", uuid.New())

		ctx.Next()

		logging.LogInfo(logger, ctx)

		ctx.Next()
	}
}
