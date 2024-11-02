package mlmodelprovision

import (
	"net/http"

	"github.com/enable-intelligent-containerized-5g/openapi"
	"github.com/enable-intelligent-containerized-5g/openapi/models"
	"github.com/free5gc/nwdaf/internal/logger"
	"github.com/gin-gonic/gin"
)

// HTTPNwdafMlModelProvisionSubscription - Creates a new subscription to receive notifications of ML model provisioning events
// TODO: Implement HTTPNwdafMlModelInfoRequest
func HTTPNwdafMlModelInfoRequest(c *gin.Context) {
	var nfMlModelInfoRequest []models.NwdafEvent

	requestBody, err := c.GetRawData()
	if err != nil {
		logger.CfgLog.Errorf("Get Request Body error: %+v", err)
		problemDetails := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		c.JSON(http.StatusInternalServerError, problemDetails)
		return
	}

	err = openapi.Deserialize(&nfMlModelInfoRequest, requestBody, "application/json")
	if err != nil {
		problemDetail := "[Request Body] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.CfgLog.Errorf(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	// req := httpwrapper.NewRequest(c.Request, nfMlModelInfoRequest)
	// req.Params[""] =
}
