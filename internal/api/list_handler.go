package api

import (
	"net/http"
	"trellode-go/internal/models"
	"trellode-go/internal/utils/logging"

	"github.com/gin-gonic/gin"
)

func (s *server) getList(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	id := c.Param("id")

	list, severity, err := s.listService.GetList(context, id)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(severity, gin.H{"detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (s *server) createList(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	var list models.List
	if err := c.BindJSON(&list); err == nil {
		board, severity, err := s.listService.CreateList(context, &list)
		if err != nil {
			logging.LogError(s.Log, c, err.Error())
			c.JSON(severity, gin.H{"detail": err.Error()})
			return
		}
		c.JSON(severity, board)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (s *server) updateList(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	var list models.List
	id := c.Param("id")
	if err := c.BindJSON(&list); err == nil {
		if id != list.ID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID in URL and body must match"})
		}
		severity, err := s.listService.UpdateList(context, &list)
		if err != nil {
			logging.LogError(s.Log, c, err.Error())
			c.JSON(severity, gin.H{"detail": err.Error()})
			return
		}
		c.JSON(severity, nil)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

type ReorderBody struct {
	IDsOrdered string `json:"idsordered"`
}

func (s *server) updateCardsOrder(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
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
			c.JSON(severity, gin.H{"detail": err.Error()})
			return
		}
		c.JSON(severity, nil)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (s *server) deleteList(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	id := c.Param("id")

	severity, err := s.listService.DeleteList(context, id)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(severity, gin.H{"detail": err.Error()})
		return
	}

	c.JSON(severity, nil)
}
