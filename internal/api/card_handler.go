package api

import (
	"net/http"
	"strconv"
	"trellode-go/internal/models"
	"trellode-go/internal/utils/logging"

	"github.com/gin-gonic/gin"
)

func (s *server) getCard(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	idValue := c.Param("id")
	id, err := strconv.Atoi(idValue)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	card, severity, err := s.cardService.GetCard(context, id)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(severity, gin.H{"detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, card)
}

func (s *server) createCard(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	var card models.Card
	if err := c.BindJSON(&card); err == nil {
		list, severity, err := s.cardService.CreateCard(context, &card)
		if err != nil {
			logging.LogError(s.Log, c, err.Error())
			c.JSON(severity, gin.H{"detail": err.Error()})
			return
		}
		c.JSON(severity, list)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (s *server) updateCard(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	var card models.Card
	id := c.Param("id")
	if err := c.BindJSON(&card); err == nil {
		if id != strconv.Itoa(card.ID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID in URL and body must match"})
		}
		severity, err := s.cardService.UpdateCard(context, &card)
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

func (s *server) deleteCard(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	idValue := c.Param("id")
	id, err := strconv.Atoi(idValue)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	severity, err := s.cardService.DeleteCard(context, id)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(severity, gin.H{"detail": err.Error()})
		return
	}

	c.JSON(severity, nil)
}
