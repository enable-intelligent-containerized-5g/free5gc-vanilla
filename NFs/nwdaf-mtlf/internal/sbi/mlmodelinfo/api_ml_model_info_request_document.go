package mlmodelinfo

import (
	"net/http"

	"github.com/enable-intelligent-containerized-5g/openapi"
	"github.com/enable-intelligent-containerized-5g/openapi/models"
	"github.com/free5gc/nwdaf/internal/logger"
	"github.com/free5gc/nwdaf/internal/sbi/producer"
	"github.com/free5gc/util/httpwrapper"
	"github.com/gin-gonic/gin"
)

// HTTPNwdafMlModelProvisionSubscription - Creates a new subscription to receive notifications of ML model provisioning events
// TODO: Implement HTTPNwdafMlModelInfoRequest
func HTTPNwdafMlModelInfoRequest(c *gin.Context) {
	req := httpwrapper.NewRequest(c.Request, nil)
	// req.Params["nfInstanceID"] = c.Params.ByName("nfInstanceID")

	httpResponse := producer.HandleNwdafMlModelInfoRequest(req)

	responseBody, err := openapi.Serialize(httpResponse.Body, "application/json")
	if err != nil {
		logger.ManagementLog.Warnln(err)
		problemDetails := models.ProblemDetails{
			Status: http.StatusUnprocessableEntity,
			Cause:  "SYSTEM_FAILURE",
			Detail: err.Error(),
		}
		c.JSON(http.StatusUnprocessableEntity, problemDetails)
	} else {
		c.Data(httpResponse.Status, "application/json", responseBody)
	}
}
