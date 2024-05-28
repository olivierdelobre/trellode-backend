package api

import (
	"fmt"
	"net/http"
	"os"
	"trellode-go/internal/models"

	"github.com/gin-gonic/gin"
)

type LivenessResponse struct {
	Items LivenessApi `json:"items"`
}

type LivenessApi struct {
	API LivenessList `json:"api"`
}

type LivenessList struct {
	List LivenessItem `json:"list"`
}

type LivenessItem struct {
	Critical int    `json:"critical"`
	Action   string `json:"action"`
	Label    string `json:"label"`
	Status   string `json:"status"`
}

func (s *server) getHealthcheck(c *gin.Context) {
	//s.Log.Info("[Monitoring - getHealthcheck] called")

	c.JSON(http.StatusOK, gin.H{"status": "ok"})

}

func (s *server) getLiveness(c *gin.Context) {
	format := c.Query("format")

	errorLabel := ""

	getListOK := true
	_, _, err := s.boardService.GetBoards(models.Context{})
	if err != nil {
		errorLabel = "Failed to get boards: " + err.Error()
		getListOK = false
	}

	if format == "metrics" {
		status := 0
		if getListOK {
			status = 1
		}
		output := fmt.Sprintf(`## HELP trellodeapi_status Lists API status: 1=OK, 0=KO
# TYPE trellodeapi_status gauge
trellodeapi_status{component="global", line="%s"} %d`, os.Getenv("ENVIRONMENT"), status)
		c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(output))
		return
	}

	if format == "json" {
		getListStatus := "ok"
		if !getListOK {
			getListStatus = "ko"
		}

		response := LivenessResponse{
			LivenessApi{
				LivenessList{
					LivenessItem{
						Critical: 1,
						Action:   "GetList",
						Label:    "Get list by ID",
						Status:   getListStatus,
					},
				},
			},
		}
		c.JSON(http.StatusOK, response)
		return
	}

	if errorLabel != "" {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": errorLabel})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
