package api

import (
	"net/http"
	"trellode-go/internal/models"
	"trellode-go/internal/utils/logging"
	"trellode-go/internal/utils/messages"

	toolbox_api "github.com/epfl-si/go-toolbox/api"
	"github.com/gin-gonic/gin"
)

func (s *server) getChecklist(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "GetContextFailure"), err.Error(), "", nil))
		return
	}

	id := c.Param("id")

	checklist, severity, err := s.checklistService.GetChecklist(context, id)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(severity, toolbox_api.MakeError(c, "", severity, messages.GetMessage(context.Lang, "GetCheckListFailure"), err.Error(), "", nil))
		return
	}
	c.JSON(http.StatusOK, checklist)
}

func (s *server) createChecklist(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "GetContextFailure"), err.Error(), "", nil))
		return
	}

	var checklist models.Checklist
	if err := c.BindJSON(&checklist); err == nil {
		list, severity, err := s.checklistService.CreateChecklist(context, &checklist)
		if err != nil {
			logging.LogError(s.Log, c, err.Error())
			c.JSON(severity, toolbox_api.MakeError(c, "", severity, messages.GetMessage(context.Lang, "CreateCheckListFailure"), err.Error(), "", nil))
			return
		}
		c.JSON(severity, list)
	} else {
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "InvalidJson"), err.Error(), "", nil))
	}
}

func (s *server) updateChecklist(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "GetContextFailure"), err.Error(), "", nil))
		return
	}

	var checklist models.Checklist
	id := c.Param("id")
	if err := c.BindJSON(&checklist); err == nil {
		if id != checklist.ID {
			c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "IdNotMatching"), "", "", nil))
		}
		severity, err := s.checklistService.UpdateChecklist(context, &checklist)
		if err != nil {
			logging.LogError(s.Log, c, err.Error())
			c.JSON(severity, toolbox_api.MakeError(c, "", severity, messages.GetMessage(context.Lang, "UpdateCheckListFailure"), err.Error(), "", nil))
			return
		}
		c.JSON(severity, nil)
	} else {
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "InvalidJson"), err.Error(), "", nil))
	}
}

func (s *server) deleteChecklist(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "GetContextFailure"), err.Error(), "", nil))
		return
	}

	id := c.Param("id")

	severity, err := s.checklistService.DeleteChecklist(context, id)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(severity, toolbox_api.MakeError(c, "", severity, messages.GetMessage(context.Lang, "DeleteCheckListFailure"), err.Error(), "", nil))
		return
	}

	c.JSON(severity, nil)
}

func (s *server) getChecklistItem(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "GetContextFailure"), err.Error(), "", nil))
		return
	}

	id := c.Param("id")

	checklist, severity, err := s.checklistService.GetChecklistItem(context, id)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(severity, toolbox_api.MakeError(c, "", severity, messages.GetMessage(context.Lang, "GetCheckListItemFailure"), err.Error(), "", nil))
		return
	}
	c.JSON(http.StatusOK, checklist)
}

func (s *server) createChecklistItem(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "GetContextFailure"), err.Error(), "", nil))
		return
	}

	var checklistItem models.ChecklistItem
	if err := c.BindJSON(&checklistItem); err == nil {
		list, severity, err := s.checklistService.CreateChecklistItem(context, &checklistItem)
		if err != nil {
			logging.LogError(s.Log, c, err.Error())
			c.JSON(severity, toolbox_api.MakeError(c, "", severity, messages.GetMessage(context.Lang, "CreateCheckListItemFailure"), err.Error(), "", nil))
			return
		}
		c.JSON(severity, list)
	} else {
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "InvalidJson"), err.Error(), "", nil))
	}
}

func (s *server) updateChecklistItem(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "GetContextFailure"), err.Error(), "", nil))
		return
	}

	var checklistItem models.ChecklistItem
	id := c.Param("id")
	if err := c.BindJSON(&checklistItem); err == nil {
		if id != checklistItem.ID {
			c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "IdNotMatching"), "", "", nil))
		}
		severity, err := s.checklistService.UpdateChecklistItem(context, &checklistItem)
		if err != nil {
			logging.LogError(s.Log, c, err.Error())
			c.JSON(severity, toolbox_api.MakeError(c, "", severity, messages.GetMessage(context.Lang, "UpdateCheckListItemFailure"), err.Error(), "", nil))
			return
		}
		c.JSON(severity, nil)
	} else {
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "InvalidJson"), err.Error(), "", nil))
	}
}

func (s *server) deleteChecklistItem(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "GetContextFailure"), err.Error(), "", nil))
		return
	}

	id := c.Param("id")

	severity, err := s.checklistService.DeleteChecklistItem(context, id)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(severity, toolbox_api.MakeError(c, "", severity, messages.GetMessage(context.Lang, "DeleteCheckListItemFailure"), err.Error(), "", nil))
		return
	}

	c.JSON(severity, nil)
}

func (s *server) updateChecklistItemsOrder(c *gin.Context) {
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
		severity, err := s.checklistService.UpdateChecklistItemsOrder(context, id, body.IDsOrdered)
		if err != nil {
			logging.LogError(s.Log, c, err.Error())
			c.JSON(severity, toolbox_api.MakeError(c, "", severity, messages.GetMessage(context.Lang, "UpdateCheckListItemsOrderFailure"), err.Error(), "", nil))
			return
		}
		c.JSON(severity, nil)
	} else {
		c.JSON(http.StatusBadRequest, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(context.Lang, "InvalidJson"), err.Error(), "", nil))
	}
}
