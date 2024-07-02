package api

import (
	"net/http"
	"trellode-go/internal/models"
	"trellode-go/internal/utils/logging"
	"trellode-go/internal/utils/messages"

	toolbox_api "github.com/epfl-si/go-toolbox/api"
	"github.com/gin-gonic/gin"
)

func (s *server) getBoard(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "GetContextFailure"), err.Error(), "", nil))
		return
	}

	id := c.Param("id")

	board, severity, err := s.boardService.GetBoard(context, id)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(severity, toolbox_api.MakeError(c, "", severity, messages.GetMessage(context.Lang, "GetBoardFailure"), err.Error(), "", nil))
		return
	}
	c.JSON(http.StatusOK, board)
}

func (s *server) getBoards(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "GetContextFailure"), err.Error(), "", nil))
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
		c.JSON(severity, toolbox_api.MakeError(c, "", severity, messages.GetMessage(context.Lang, "GetBoardsFailure"), err.Error(), "", nil))
		return
	}
	c.JSON(http.StatusOK, boards)
}

func (s *server) createBoard(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "GetContextFailure"), err.Error(), "", nil))
		return
	}

	var board models.Board
	if err := c.BindJSON(&board); err == nil {
		board, severity, err := s.boardService.CreateBoard(context, &board)
		if err != nil {
			logging.LogError(s.Log, c, err.Error())
			c.JSON(severity, toolbox_api.MakeError(c, "", severity, messages.GetMessage(context.Lang, "CreateBoardFailure"), err.Error(), "", nil))
			return
		}
		c.JSON(severity, board)
	} else {
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "InvalidJson"), err.Error(), "", nil))
	}
}

func (s *server) updateBoard(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "GetContextFailure"), err.Error(), "", nil))
		return
	}

	var board models.Board
	id := c.Param("id")

	if err := c.BindJSON(&board); err == nil {
		if id != board.ID {
			c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "IdNotMatching"), "", "", nil))
		}
		severity, err := s.boardService.UpdateBoard(context, id, &board)
		if err != nil {
			logging.LogError(s.Log, c, err.Error())
			c.JSON(severity, toolbox_api.MakeError(c, "", severity, messages.GetMessage(context.Lang, "UpdateBoardFailure"), err.Error(), "", nil))
			return
		}
		c.JSON(severity, board)
	} else {
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "InvalidJson"), err.Error(), "", nil))
	}
}

func (s *server) deleteBoard(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "GetContextFailure"), err.Error(), "", nil))
		return
	}

	id := c.Param("id")

	severity, err := s.boardService.DeleteBoard(context, id)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(severity, toolbox_api.MakeError(c, "", severity, messages.GetMessage(context.Lang, "DeleteBoardFailure"), err.Error(), "", nil))
		return
	}

	c.JSON(severity, nil)
}

func (s *server) updateListsOrder(c *gin.Context) {
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
		severity, err := s.boardService.UpdateListsOrder(context, id, body.IDsOrdered)
		if err != nil {
			logging.LogError(s.Log, c, err.Error())
			c.JSON(severity, toolbox_api.MakeError(c, "", severity, messages.GetMessage(context.Lang, "UpdateListsOrderFailure"), err.Error(), "", nil))
			return
		}
		c.JSON(severity, nil)
	} else {
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "InvalidJson"), err.Error(), "", nil))
	}
}
