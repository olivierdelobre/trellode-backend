package api

import (
	"net/http"
	"strconv"
	"trellode-go/internal/models"
	"trellode-go/internal/utils/logging"

	"github.com/gin-gonic/gin"
)

func (s *server) getBackground(c *gin.Context) {
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

	background, severity, err := s.backgroundService.GetBackground(context, id)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(severity, gin.H{"detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, background)
}

func (s *server) getBackgrounds(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	backgrounds, severity, err := s.backgroundService.GetBackgrounds(context)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(severity, gin.H{"detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, backgrounds)
}

/*
	func (s *server) createBackground(c *gin.Context) {
		context, err := getContext(c)
		if err != nil {
			logging.LogError(s.Log, c, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
			return
		}

		file, header, err := c.Request.FormFile("file")
		if err != nil {
			logging.LogError(s.Log, c, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
			return
		}
		log.Println(header.Filename)
		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, file); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"detail": err.Error()})
			return
		}

		id, severity, err := s.backgroundService.CreateBackground(context, buf.Bytes())
		if err != nil {
			logging.LogError(s.Log, c, err.Error())
			c.JSON(severity, gin.H{"detail": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"id": id})
	}
*/
func (s *server) createBackground(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	var background models.Background
	if err := c.BindJSON(&background); err == nil {
		id, severity, err := s.backgroundService.CreateBackground(context, background.Data)
		if err != nil {
			logging.LogError(s.Log, c, err.Error())
			c.JSON(severity, gin.H{"detail": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": id})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (s *server) deleteBackground(c *gin.Context) {
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

	severity, err := s.backgroundService.DeleteBackground(context, id)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(severity, gin.H{"detail": err.Error()})
		return
	}

	c.JSON(severity, nil)
}
