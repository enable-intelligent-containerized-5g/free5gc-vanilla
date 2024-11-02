package analyticsinfo

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	// "github.com/antihax/optional"
	"github.com/enable-intelligent-containerized-5g/openapi"
	"github.com/enable-intelligent-containerized-5g/openapi/Nnrf_NFDiscovery"
	"github.com/enable-intelligent-containerized-5g/openapi/models"
	"github.com/free5gc/nwdaf/internal/logger"
	"github.com/free5gc/nwdaf/internal/sbi/consumer"
	"github.com/free5gc/util/httpwrapper"
	"github.com/gin-gonic/gin"
)

// HTTPNwdafAnalyticsInfoRequest - Creates a new subscription to receive notifications of ML model provisioning events
// TODO: Implement HTTPNwdafAnalyticsInfoRequest
func HTTPNwdafAnalyticsInfoRequest(c *gin.Context) {
	var nfAnalyticsInfoRequest models.NwdafAnalyticsInfoRequest

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

	// searchAllNfs() // Search NFs
	searchServiceAnaliticsInfo() // Searh Service AnaliticsInfo

	if err := json.Unmarshal(requestBody, &nfAnalyticsInfoRequest); err != nil {
		logger.CfgLog.Errorf("Unmarshal error: %+v", err)
		c.JSON(http.StatusBadRequest, models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: "[Request Body] " + err.Error(),
		})
		return
	}

	err = openapi.Deserialize(&nfAnalyticsInfoRequest, requestBody, "application/json")
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

	// Validar el evento
	analyticsID, err := isValidEvent(nfAnalyticsInfoRequest.EventId)
	if err != nil {
		problemDetail := "[Invalid event] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Missing mandatory parameter",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.CfgLog.Errorf(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	// Validar NfSetId o NfTypes
	// TODO: Implementar la validación de NfSetId o NfTypes
	// err = isValidNfSetIdOrNfTypes(nfAnalyticsInfoRequest.NfSetIds, nfAnalyticsInfoRequest.NfTypes)
	// if err != nil {
	// 	problemDetail := "[Invalid NfSetId or NfTypes] " + err.Error()
	// 	rsp := models.ProblemDetails{
	// 		Title:  "Missing mandatory parameter",
	// 		Status: http.StatusBadRequest,
	// 		Detail: problemDetail,
	// 	}
	// 	logger.CfgLog.Errorf(problemDetail)
	// 	c.JSON(http.StatusBadRequest, rsp)
	// 	return
	// }

	if analyticsID == string(models.NwdafEvent_NF_LOAD) {
		fmt.Println("EventId", analyticsID)
	}

	req := httpwrapper.NewRequest(c.Request, nfAnalyticsInfoRequest)
	req.Params["analytics_id"] = c.Params.ByName("analytics_id")

	// Imprimir el request
	fmt.Println(req)
}

// isValidNfSetIdOrNfTypes verifica si el NfSetId o NfTypes
// func isValidNfSetIdOrNfTypes(nfSetId []string, nfTypes []models.NrfNfManagementNfType) (err error) {
// 	if len(nfSetId) == 0 && len(nfTypes) == 0 {
// 		err = errors.New("please provide a valid NfSetId or NfTypes")
// 		return err
// 	}

// 	// if valid, exists := models.ValidNrfNfManagementNfType[nfTypes]
// }

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

func searchAllNfs() error {
	param := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		// ServiceNames: optional.Interface{},
	}
	allNfs := consumer.SearchAllNfInstance("http://127.0.0.1:30050", "", models.NfType_NWDAF, param)
	// fmt.Println("allNfs", allNfs)

	return allNfs
}

func searchServiceAnaliticsInfo() error {
	param := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		// ServiceNames: optional.Interface{},
	}
	allNfs := consumer.SearchAnaliticsInfoInstance("http://127.0.0.1:30050", models.NfType_NWDAF, models.NfType_NWDAF, param)
	// fmt.Println("ServiceMlModelProvision: ", allNfs)

	return allNfs
}
