package producer

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/enable-intelligent-containerized-5g/openapi"
	"github.com/enable-intelligent-containerized-5g/openapi/models"
	"github.com/free5gc/nrf/internal/logger"
	"github.com/free5gc/util/httpwrapper"
	timedecode "github.com/free5gc/util/mapstruct"
	"github.com/free5gc/util/mongoapi"
)

func RemoveInactiveNfs() *models.ProblemDetails {
	logger.DiscoveryLog.Warnf("Function RemoveInactiveNfs started at: %s", time.Now())

	// Crear el objeto url.Values
	queryParameters := url.Values{}
	// Agregar dos parámetros clave-valor
	queryParameters.Add("target-nf-type", "")
	queryParameters.Add("requester-nf-type", "NRF")

	filter := buildFilter(queryParameters)

	nfProfilesRaw, err := mongoapi.RestfulAPIGetMany("NfProfile", filter)
	if err != nil {
		problemDetails := &models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		return problemDetails
	}

	// nfProfile data for response
	var nfProfilesStruct []models.NfProfile
	if err := timedecode.Decode(nfProfilesRaw, &nfProfilesStruct); err != nil {
		problemDetails := &models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		return problemDetails
	}

	for i, nfProfile := range nfProfilesStruct {
		// Get ID
		nfInstanceId := nfProfile.NfInstanceId

		services := *nfProfile.NfServices
		apiPrefix := services[0].ApiPrefix

		// Difine fullURI
		baseUri := apiPrefix + "/common-nfprofileprovition/v1"
		// api := services[0].ServiceName
		// versions := *services[0].Versions
		// version := versions[0].ApiVersionInUri
		endpoint := "nfprofileprovition/request"
		// fullURL := fmt.Sprintf("%s/%s/%s", baseUri, api, version)
		fullURL := fmt.Sprintf("%s/%s", baseUri, endpoint)

		if !validateNfActive(fullURL, nfInstanceId) {
			logger.DiscoveryLog.Infof("Delete NF #%d: %s", i, nfInstanceId)
			err := NFDeregisterProcedure(nfInstanceId)
			if err != nil {
				logger.DiscoveryLog.Errorf("Error deleting: %s", nfInstanceId)
			}
		} else {
			logger.DiscoveryLog.Infof("Ignore NF #%d: %s", i, nfInstanceId)
		}
	}

	logger.DiscoveryLog.Infof("Function completed at: %s", time.Now())

	return nil
}

func validateNfActive(url string, nfInstanceId string) bool {
  logger.DiscoveryLog.Infoln("URL: ", url)
	// Realizar la solicitud HTTP GET
	resp, err := http.Get(url)
	if err != nil {
		logger.DiscoveryLog.Error("Error al hacer la solicitud:", err)
		return false
	}
	defer resp.Body.Close()

	// Leer el cuerpo de la respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.DiscoveryLog.Error("Error al leer el cuerpo de la respuesta:", err)
		return false
	}

  // Deserializar el cuerpo de la respuesta
  var data httpwrapper.Response
  err = openapi.Deserialize(&data, body, "application/json")
  if err != nil {
    logger.DiscoveryLog.Error("Error al deserializar el cuerpo de la respuesta:", err)
    return false
  }
  nfId := data.Body

	// Imprimir el cuerpo de la respuesta
	logger.DiscoveryLog.Info("nfId: ", nfId)

	// Verificar si la respuesta coincide con el ID del NF
	return nfId == nfInstanceId
}
