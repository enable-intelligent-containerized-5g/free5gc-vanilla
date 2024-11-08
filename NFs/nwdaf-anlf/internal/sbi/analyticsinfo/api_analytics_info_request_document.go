package analyticsinfo

import (
	"errors"
	"fmt"
	"net/http"

	// "github.com/antihax/optional"
	"github.com/enable-intelligent-containerized-5g/openapi"
	"github.com/enable-intelligent-containerized-5g/openapi/models"
	"github.com/free5gc/nwdaf/internal/logger"
	"github.com/free5gc/nwdaf/internal/sbi/producer"
	"github.com/free5gc/util/httpwrapper"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// HTTPNwdafAnalyticsInfoRequest - Creates a new subscription to receive notifications of ML model provisioning events
func HTTPNwdafAnalyticsInfoRequest(c *gin.Context) {
	logger.AniLog.Info("HTTP NwdafAnalyticsInfoRequest")
	// nfAnalyticsInfoRequest := models.NewNwdafAnalyticsInfoRequest()
	nfAnalyticsInfoRequest := models.NwdafAnalyticsInfoRequest{}

	// Get Request Body
	requestBody, err := c.GetRawData()
	if err != nil {
		logger.AniLog.Errorln("Get Request Body error: ", err)
		problemDetails := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		c.JSON(http.StatusInternalServerError, problemDetails)
		return
	}

	// Deserialize Request Body
	err = openapi.Deserialize(&nfAnalyticsInfoRequest, requestBody, "application/json")
	if err != nil {
		problemDetail := "[Request Body] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.AniLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	// Initialize validater
	validate := validator.New()

	// Validate struct and return error if any
	if err := validate.Struct(nfAnalyticsInfoRequest); err != nil {
		problemDetails := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: err.Error(),
		}
		logger.AniLog.Errorln("Validation Error: ", err)
		c.JSON(http.StatusBadRequest, problemDetails)
		return
	}

	// Validate EventId and return the event name
	analyticsID, err := isValidEvent(nfAnalyticsInfoRequest.EventId)
	if err != nil {
		problemDetail := "[Invalid event] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Missing mandatory parameter",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.AniLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	// Validate NfInstanceIds or NfTypes and return the typePayload
	// TODO: Implementar la validación de NfSetId o NfTypes
	typePayload, err := isValidNfInstanceIdsOrNfTypes(nfAnalyticsInfoRequest.NfInstanceIds, nfAnalyticsInfoRequest.NfTypes)
	if err != nil {
		problemDetail := "[Invalid NfSetId or NfTypes] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Missing mandatory parameter",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.AniLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	// Create a new request with the information of the request
	req := httpwrapper.NewRequest(c.Request, nfAnalyticsInfoRequest)

	// Call the producer to send the request
	var rsp *httpwrapper.Response

	if analyticsID == string(models.EventId_NF_LOAD) {
		// Execute the NF Load event
		rsp = producer.HandleAnalyticsInfoNfLoadMetricsNew(req, typePayload)
	} else {
		bodyMessage := fmt.Sprintf("The eventId %s is not implemented", analyticsID)
		rsp = httpwrapper.NewResponse(http.StatusNotImplemented, nil, bodyMessage)
		logger.AniLog.Warnf(bodyMessage)
	}

	responseBody, err := openapi.Serialize(rsp.Body, "application/json")
	if err != nil {
		logger.AniLog.Errorln(err)
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

// isValidNfInstanceIdsOrNfTypes verifica si el NfInstanceIds o NfTypes
func isValidNfInstanceIdsOrNfTypes(nfInstanceIds []string, nfTypes []models.NfType) (typePayload models.TypePayloadRequest, err error) {
	if nfInstanceIds == nil && nfTypes== nil {
		err = errors.New("please provide a valid NfInstanceIds or NfTypes")
	}

	if len(nfInstanceIds) == 0 && len(nfTypes) == 0 {
		err = errors.New("please provide a valid NfInstanceIds or NfTypes")
	}

	if len(nfTypes) > 0 {
		for _, nfType := range nfTypes {
			if valid, exists := models.ValidNfType[nfType]; exists || valid {
				typePayload = "NF_TYPES"
			} else {
				err = errors.New("please provide a valid NfType")
			}
		}
	} else {
		typePayload = "NF_INSTANCES"
	}

	return typePayload, err
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
