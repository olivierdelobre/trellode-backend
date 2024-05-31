package api

import (
	"net/http"
	"trellode-go/internal/models"
	"trellode-go/internal/utils/logging"

	"github.com/gin-gonic/gin"
)

func (s *server) getBoard(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	id := c.Param("id")

	board, severity, err := s.boardService.GetBoard(context, id)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(severity, gin.H{"detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, board)
}

func (s *server) getBoards(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	archivedValue := c.Query("archived")
	archived := false
	if archivedValue == "1" {
		archived = true
	}

	boards, severity, err := s.boardService.GetBoards(context, archived)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(severity, gin.H{"detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, boards)
}

func (s *server) createBoard(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	var board models.Board
	if err := c.BindJSON(&board); err == nil {
		board, severity, err := s.boardService.CreateBoard(context, &board)
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

func (s *server) updateBoard(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	var board models.Board
	id := c.Param("id")

	if err := c.BindJSON(&board); err == nil {
		if id != board.ID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID in URL and body must match"})
		}
		severity, err := s.boardService.UpdateBoard(context, id, &board)
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

func (s *server) deleteBoard(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	id := c.Param("id")

	severity, err := s.boardService.DeleteBoard(context, id)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(severity, gin.H{"detail": err.Error()})
		return
	}

	c.JSON(severity, nil)
}
