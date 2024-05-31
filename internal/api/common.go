package api

import (
	"errors"
	"fmt"
	"trellode-go/internal/models"
	"trellode-go/internal/utils/messages"

	"github.com/gin-gonic/gin"
)

func getContext(c *gin.Context) (models.Context, error) {
	lang, _ := c.Get("lang")
	langStr := fmt.Sprintf("%s", lang)

	userIdValue, _ := c.Get("userId")
	userId := fmt.Sprintf("%s", userIdValue)
	if userId == "" {
		return models.Context{}, errors.New(messages.GetMessage(langStr, "NoUserId"))
	}

	userTypeValue, _ := c.Get("userType")
	userType := fmt.Sprintf("%s", userTypeValue)

	return models.Context{
		UserId:   userId,
		UserType: userType,
		Lang:     langStr,
	}, nil
}
