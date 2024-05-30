package api

import (
	"net/http"
	"strconv"
	"trellode-go/internal/utils/logging"

	"github.com/gin-gonic/gin"
)

func (s *server) getLogs(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	boardIdValue := c.Query("boardid")
	boardId, err := strconv.Atoi(boardIdValue)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	backgrounds, severity, err := s.logService.GetLogs(context, boardId)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(severity, gin.H{"detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, backgrounds)
}
