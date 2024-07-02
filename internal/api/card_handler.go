package api

import (
	"net/http"
	"trellode-go/internal/models"
	"trellode-go/internal/utils/logging"
	"trellode-go/internal/utils/messages"

	toolbox_api "github.com/epfl-si/go-toolbox/api"
	"github.com/gin-gonic/gin"
)

func (s *server) getCard(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "GetContextFailure"), err.Error(), "", nil))
		return
	}

	id := c.Param("id")

	card, severity, err := s.cardService.GetCard(context, id)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(severity, toolbox_api.MakeError(c, "", severity, messages.GetMessage(context.Lang, "GetCardFailure"), err.Error(), "", nil))
		return
	}
	c.JSON(http.StatusOK, card)
}

func (s *server) createCard(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "GetContextFailure"), err.Error(), "", nil))
		return
	}

	var card models.Card
	if err := c.BindJSON(&card); err == nil {
		list, severity, err := s.cardService.CreateCard(context, &card)
		if err != nil {
			logging.LogError(s.Log, c, err.Error())
			c.JSON(severity, toolbox_api.MakeError(c, "", severity, messages.GetMessage(context.Lang, "CreateCardFailure"), err.Error(), "", nil))
			return
		}
		c.JSON(severity, list)
	} else {
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "InvalidJson"), err.Error(), "", nil))
	}
}

func (s *server) updateCard(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "GetContextFailure"), err.Error(), "", nil))
		return
	}

	var card models.Card
	id := c.Param("id")
	if err := c.BindJSON(&card); err == nil {
		if id != card.ID {
			c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "IdNotMatching"), "", "", nil))
		}
		severity, err := s.cardService.UpdateCard(context, &card)
		if err != nil {
			logging.LogError(s.Log, c, err.Error())
			c.JSON(severity, toolbox_api.MakeError(c, "", severity, messages.GetMessage(context.Lang, "UpdateCardFailure"), err.Error(), "", nil))
			return
		}
		c.JSON(severity, nil)
	} else {
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "InvalidJson"), err.Error(), "", nil))
	}
}

func (s *server) deleteCard(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "GetContextFailure"), err.Error(), "", nil))
		return
	}

	id := c.Param("id")

	severity, err := s.cardService.DeleteCard(context, id)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(severity, toolbox_api.MakeError(c, "", severity, messages.GetMessage(context.Lang, "DeleteCardFailure"), err.Error(), "", nil))
		return
	}

	c.JSON(severity, nil)
}
