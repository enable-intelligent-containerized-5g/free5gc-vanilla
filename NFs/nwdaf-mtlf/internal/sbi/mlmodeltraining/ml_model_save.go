package mlmodeltraining

import (
	"github.com/enable-intelligent-containerized-5g/openapi"
	"github.com/enable-intelligent-containerized-5g/openapi/models"
	"github.com/free5gc/nwdaf/internal/logger"
	"github.com/free5gc/nwdaf/internal/sbi/producer"
	"github.com/free5gc/util/httpwrapper"
	"github.com/gin-gonic/gin"
	"net/http"

)

func HTTPSaveMlModel(c *gin.Context) {
    logger.MlModelTrainingLog.Info("HTTP SaveMlModel")
	var modelInfo models.MlModelData

	// Step 1: Retrieve HTTP request body
	requestBody, err := c.GetRawData()
	if err != nil {
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	// Step 2: Convert requestBody to MlModelInfo model
	err = openapi.Deserialize(&modelInfo, requestBody, "application/json")
	if err != nil {
		problemDetail := "[Request Body] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	req := httpwrapper.NewRequest(c.Request, modelInfo)
	httpResponse := producer.HandleSaveMlModel(req)
	responseBody, err := openapi.Serialize(httpResponse.Body, "application/json")
	if err != nil {
		logger.MlModelTrainingLog.Errorln(err)
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
