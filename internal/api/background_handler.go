package api

import (
	"net/http"
	"strconv"
	"trellode-go/internal/models"
	"trellode-go/internal/utils/logging"
	"trellode-go/internal/utils/messages"

	toolbox_api "github.com/epfl-si/go-toolbox/api"
	"github.com/gin-gonic/gin"
)

func (s *server) getBackground(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "GetContextFailure"), err.Error(), "", nil))
		return
	}

	idValue := c.Param("id")
	id, err := strconv.Atoi(idValue)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "NotAnInt"), err.Error(), "", nil))
		return
	}

	background, severity, err := s.backgroundService.GetBackground(context, id)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(severity, toolbox_api.MakeError(c, "", severity, messages.GetMessage(context.Lang, "GetBackgroundFailure"), err.Error(), "", nil))
		return
	}
	c.JSON(http.StatusOK, background)
}

func (s *server) getBackgrounds(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "GetContextFailure"), err.Error(), "", nil))
		return
	}

	backgrounds, severity, err := s.backgroundService.GetBackgrounds(context)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(severity, toolbox_api.MakeError(c, "", severity, messages.GetMessage(context.Lang, "GetBackgroundsFailure"), err.Error(), "", nil))
		return
	}
	c.JSON(http.StatusOK, backgrounds)
}

func (s *server) createBackground(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "GetContextFailure"), err.Error(), "", nil))
		return
	}

	var background models.Background
	if err := c.BindJSON(&background); err == nil {
		id, severity, err := s.backgroundService.CreateBackground(context, background.Data)
		if err != nil {
			logging.LogError(s.Log, c, err.Error())
			c.JSON(severity, toolbox_api.MakeError(c, "", severity, messages.GetMessage(context.Lang, "CreateBackgroundFailure"), err.Error(), "", nil))
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": id})
	} else {
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "InvalidJson"), err.Error(), "", nil))
	}
}

func (s *server) deleteBackground(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "GetContextFailure"), err.Error(), "", nil))
		return
	}

	idValue := c.Param("id")
	id, err := strconv.Atoi(idValue)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "NotAnInt"), err.Error(), "", nil))
		return
	}

	severity, err := s.backgroundService.DeleteBackground(context, id)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(severity, toolbox_api.MakeError(c, "", severity, messages.GetMessage(context.Lang, "DeleteBackgroundFailure"), err.Error(), "", nil))
		return
	}

	c.JSON(severity, nil)
}
