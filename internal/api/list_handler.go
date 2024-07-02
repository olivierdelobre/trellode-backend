package api

import (
	"net/http"
	"trellode-go/internal/models"
	"trellode-go/internal/utils/logging"
	"trellode-go/internal/utils/messages"

	toolbox_api "github.com/epfl-si/go-toolbox/api"
	"github.com/gin-gonic/gin"
)

func (s *server) getList(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "GetContextFailure"), err.Error(), "", nil))
		return
	}

	id := c.Param("id")

	list, severity, err := s.listService.GetList(context, id)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(severity, toolbox_api.MakeError(c, "", severity, messages.GetMessage(context.Lang, "GetListFailure"), err.Error(), "", nil))
		return
	}
	c.JSON(http.StatusOK, list)
}

func (s *server) createList(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "GetContextFailure"), err.Error(), "", nil))
		return
	}

	var list models.List
	if err := c.BindJSON(&list); err == nil {
		board, severity, err := s.listService.CreateList(context, &list)
		if err != nil {
			logging.LogError(s.Log, c, err.Error())
			c.JSON(severity, toolbox_api.MakeError(c, "", severity, messages.GetMessage(context.Lang, "CreateListFailure"), err.Error(), "", nil))
			return
		}
		c.JSON(severity, board)
	} else {
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "InvalidJson"), err.Error(), "", nil))
	}
}

func (s *server) updateList(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "GetContextFailure"), err.Error(), "", nil))
		return
	}

	var list models.List
	id := c.Param("id")
	if err := c.BindJSON(&list); err == nil {
		if id != list.ID {
			c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "IdNotMatching"), "", "", nil))
		}
		severity, err := s.listService.UpdateList(context, &list)
		if err != nil {
			logging.LogError(s.Log, c, err.Error())
			c.JSON(severity, toolbox_api.MakeError(c, "", severity, messages.GetMessage(context.Lang, "UpdateListFailure"), err.Error(), "", nil))
			return
		}
		c.JSON(severity, nil)
	} else {
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "InvalidJson"), err.Error(), "", nil))
	}
}

type ReorderBody struct {
	IDsOrdered string `json:"idsordered"`
}

func (s *server) updateCardsOrder(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "GetContextFailure"), err.Error(), "", nil))
		return
	}

	id := c.Param("id")
	var body ReorderBody
	if err := c.BindJSON(&body); err == nil {
		if body.IDsOrdered == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "idsordered is required"})
			return
		}
		severity, err := s.listService.UpdateCardsOrder(context, id, body.IDsOrdered)
		if err != nil {
			logging.LogError(s.Log, c, err.Error())
			c.JSON(severity, toolbox_api.MakeError(c, "", severity, messages.GetMessage(context.Lang, "UpdateCardsOrderFailure"), err.Error(), "", nil))
			return
		}
		c.JSON(severity, nil)
	} else {
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "InvalidJson"), err.Error(), "", nil))
	}
}

func (s *server) deleteList(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "GetContextFailure"), err.Error(), "", nil))
		return
	}

	id := c.Param("id")

	severity, err := s.listService.DeleteList(context, id)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(severity, toolbox_api.MakeError(c, "", severity, messages.GetMessage(context.Lang, "DeleteListFailure"), err.Error(), "", nil))
		return
	}

	c.JSON(severity, nil)
}
