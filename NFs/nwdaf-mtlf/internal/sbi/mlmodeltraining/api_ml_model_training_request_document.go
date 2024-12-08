package mlmodeltraining

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/enable-intelligent-containerized-5g/openapi"
	"github.com/enable-intelligent-containerized-5g/openapi/models"
	"github.com/free5gc/nwdaf/internal/logger"
	"github.com/free5gc/nwdaf/internal/sbi/producer"
	"github.com/free5gc/util/httpwrapper"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// HTTPNwdafAnalyticsInfoRequest - Creates a new subscription to receive notifications of ML model provisioning events
func HTTPNwdafMlModelTrainingRequest(c *gin.Context) {
	logger.MlModelTrainingLog.Info("HTTP HTTPNwdafMlModelTrainingRequest")
	mlModelTrainingRequest := models.NwdafMlModelTrainingRequest{}

	// Get Request Body
	requestBody, err := c.GetRawData()
	if err != nil {
		logger.MlModelTrainingLog.Errorln("Get Request Body error: ", err)
		problemDetails := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		c.JSON(http.StatusInternalServerError, problemDetails)
		return
	}



	logger.MlModelTrainingLog.Warn("After Get Request Body")

	// Deserialize Request Body
	err = openapi.Deserialize(&mlModelTrainingRequest, requestBody, "application/json")
	logger.MlModelTrainingLog.Warn(mlModelTrainingRequest, " error: ", err)

	if err != nil {
		problemDetail := "[Request Body] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.MlModelTrainingLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	logger.MlModelTrainingLog.Warn("After Deserialize Request Body")

	// Initialize validater
	validate := validator.New()

	// Validate struct and return error if any
	if err := validate.Struct(mlModelTrainingRequest); err != nil {
		problemDetails := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: err.Error(),
		}
		logger.MlModelTrainingLog.Errorln("Validation Error: ", err)
		c.JSON(http.StatusBadRequest, problemDetails)
		return
	}

	// Validate EventId and return the event name
	analyticsID, err := isValidEvent(mlModelTrainingRequest.EventId)
	if err != nil {
		problemDetail := "[Invalid event] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Missing mandatory parameter",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.MlModelTrainingLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	// Create a new request with the information of the request
	req := httpwrapper.NewRequest(c.Request, mlModelTrainingRequest)

	// Call the producer to send the request
	var rsp *httpwrapper.Response

	if analyticsID == string(models.EventId_NF_LOAD) {
		// Execute the NF Load event
		rsp = producer.HandleMlModelTrainingNfLoadMetric(req)
	} else {
		bodyMessage := fmt.Sprintf("The eventId %s is not implemented", analyticsID)
		rsp = httpwrapper.NewResponse(http.StatusNotImplemented, nil, bodyMessage)
		logger.MlModelTrainingLog.Warnf(bodyMessage)
	}

	responseBody, err := openapi.Serialize(rsp.Body, "application/json")
	if err != nil {
		logger.MlModelTrainingLog.Errorln(err)
		problemDetails := models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "SYSTEM_FAILURE",
			Detail: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, problemDetails)
	} else {
		c.Data(rsp.Status, "application/json", responseBody)
	}
}

// isValidEvent verifica si el evento es válido y retorna el nombre del evento
func isValidEvent(event models.EventId) (analyticsID string, err error) {
	if valid, exists := models.ValidEventIds[event]; exists && valid {
		// Retorna el nombre del evento si es valido
		analyticsID = string(event)
	} else {
		// Renorna un error si el evento no es valido
		err = errors.New("please provide a valid event")
	}
	return analyticsID, err
}
