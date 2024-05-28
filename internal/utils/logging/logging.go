package logging

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LogInfo(logger *zap.Logger, ctx *gin.Context) {
	reqMethod := ctx.Request.Method
	reqUri := ctx.Request.RequestURI
	statusCode := ctx.Writer.Status()

	// Request IP
	clientIP := ctx.GetHeader("X-Forwarded-For")
	if clientIP == "" {
		clientIP = ctx.ClientIP()
	}

	// if UUID provided by KrakenD, use it, otherwise take the one generated locally
	uuid := ctx.GetHeader("X-Krakend-UUID")
	if uuid == "" {
		val, _ := ctx.Get("uuid")
		uuid = fmt.Sprintf("%v", val)
	}
	userIdValue, _ := ctx.Get("userId")
	userId := ""
	if userIdValue != nil {
		userId = fmt.Sprintf("%v", userIdValue)
	}

	logger.Info("",
		zap.String("event.dataset", os.Getenv("API_NAME")),
		zap.String("http.request.method", reqMethod),
		zap.String("url.path", reqUri),
		zap.Int("http.response.status_code", statusCode),
		zap.String("client.address", clientIP),
		zap.String("user.id", userId),
		zap.String("uuid", fmt.Sprintf("%v", uuid)),
	)
}

func LogError(logger *zap.Logger, ctx *gin.Context, message string) {
	reqMethod := ctx.Request.Method
	reqUri := ctx.Request.RequestURI
	statusCode := ctx.Writer.Status()

	// Request IP
	clientIP := ctx.ClientIP()

	uuid, _ := ctx.Get("uuid")
	userIdValue, _ := ctx.Get("userId")
	userId := ""
	if userIdValue != nil {
		userId = fmt.Sprintf("%v", userIdValue)
	}

	logger.Error(message,
		zap.String("event.dataset", os.Getenv("API_NAME")),
		zap.String("http.request.method", reqMethod),
		zap.String("url.path", reqUri),
		zap.Int("http.response.status_code", statusCode),
		zap.String("client.address", clientIP),
		zap.String("user.id", userId),
		zap.String("uuid", fmt.Sprintf("%v", uuid)),
	)
}

func LogCustom(logger *zap.Logger, level string, method string, uri string, status int, body string, msg string) {
	if level == "info" {
		logger.Info(msg,
			zap.String("event.dataset", os.Getenv("API_NAME")),
			zap.String("http.request.method", method),
			zap.String("url.path", uri),
			zap.Int("http.response.status_code", status),
		)
	}
	if level == "error" {
		logger.Error(msg,
			zap.String("event.dataset", os.Getenv("API_NAME")),
			zap.String("http.request.method", method),
			zap.String("url.path", uri),
			zap.Int("http.response.status_code", status),
		)
	}
}
