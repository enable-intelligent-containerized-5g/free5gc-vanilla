package producer

import (
	"fmt"
	"net/url"
	"time"

	"net/http"

	"github.com/enable-intelligent-containerized-5g/openapi/models"
	"github.com/free5gc/nrf/internal/logger"
	timedecode "github.com/free5gc/util/mapstruct"
	"github.com/free5gc/util/mongoapi"
)

func RemoveInactiveNfs() *models.ProblemDetails {
	logger.DiscoveryLog.Warnf("Function RemoveInactiveNfs started at: %s", time.Now())

	// Crear el objeto url.Values
	queryParameters := url.Values{}
	// Agregar dos parámetros clave-valor
	queryParameters.Add("target-nf-type", "NWDAF")
	queryParameters.Add("requester-nf-type", "NWDAF")

	var filter = buildFilter(queryParameters)

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

		// Get fullURI
		services := *nfProfile.NfServices
		baseUri := services[0].ApiPrefix
		api := services[0].ServiceName
		versions := *services[0].Versions
		version := versions[0].ApiVersionInUri
		fullURL := fmt.Sprintf("%s/%s/%s", baseUri, api, version)

		if !checkURL(fullURL) {
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

func checkURL(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
