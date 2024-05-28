package api

import (
	"errors"
	"fmt"
	"strconv"
	"trellode-go/internal/models"
	"trellode-go/internal/utils/messages"

	"github.com/gin-gonic/gin"
)

func getContext(c *gin.Context) (models.Context, error) {
	lang, _ := c.Get("lang")
	langStr := fmt.Sprintf("%s", lang)

	userIdValue, _ := c.Get("userId")
	userIdStr := fmt.Sprintf("%s", userIdValue)
	if userIdStr == "" {
		return models.Context{}, errors.New(messages.GetMessage(langStr, "NoUserId"))
	}

	userId, _ := strconv.Atoi(userIdStr)

	userTypeValue, _ := c.Get("userType")
	userType := fmt.Sprintf("%s", userTypeValue)

	return models.Context{
		UserId:   userId,
		UserType: userType,
		Lang:     langStr,
	}, nil
}
