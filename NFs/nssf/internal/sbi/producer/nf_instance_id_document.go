/*
 * NSSF NSSAI Availability
 *
 * NSSF NSSAI Availability Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package producer

import (
	"net/http"

	"github.com/free5gc/nssf/internal/logger"
	"github.com/free5gc/nssf/internal/plugin"
	"github.com/enable-intelligent-containerized-5g/openapi/models"
	"github.com/free5gc/util/httpwrapper"
)

// HandleNSSAIAvailabilityDelete - Deletes an already existing S-NSSAIs per TA
// provided by the NF service consumer (e.g AMF)
func HandleNSSAIAvailabilityDelete(request *httpwrapper.Request) *httpwrapper.Response {
	logger.Nssaiavailability.Infof("Handle NSSAIAvailabilityDelete")

	nfID := request.Params["nfId"]

	problemDetails := NSSAIAvailabilityDeleteProcedure(nfID)

	if problemDetails != nil {
		return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	return httpwrapper.NewResponse(http.StatusNoContent, nil, nil)
}

// HandleNSSAIAvailabilityPatch - Updates an already existing S-NSSAIs per TA
// provided by the NF service consumer (e.g AMF)
func HandleNSSAIAvailabilityPatch(request *httpwrapper.Request) *httpwrapper.Response {
	logger.Nssaiavailability.Infof("Handle NSSAIAvailabilityPatch")

	nssaiAvailabilityUpdateInfo := request.Body.(plugin.PatchDocument)
	nfID := request.Params["nfId"]

	// TODO: Request NfProfile of NfId from NRF
	//       Check if NfId is valid AMF and obtain AMF Set ID
	//       If NfId is invalid, return ProblemDetails with code 404 Not Found
	//       If NF consumer is not authorized to update NSSAI availability, return ProblemDetails with code 403 Forbidden

	response, problemDetails := NSSAIAvailabilityPatchProcedure(nssaiAvailabilityUpdateInfo, nfID)

	if response != nil {
		return httpwrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return httpwrapper.NewResponse(http.StatusForbidden, nil, problemDetails)
}

// HandleNSSAIAvailabilityPut - Updates/replaces the NSSF
// with the S-NSSAIs the NF service consumer (e.g AMF) supports per TA
func HandleNSSAIAvailabilityPut(request *httpwrapper.Request) *httpwrapper.Response {
	logger.Nssaiavailability.Infof("Handle NSSAIAvailabilityPut")

	nssaiAvailabilityInfo := request.Body.(models.NssaiAvailabilityInfo)
	nfID := request.Params["nfId"]

	response, problemDetails := NSSAIAvailabilityPutProcedure(nssaiAvailabilityInfo, nfID)

	if response != nil {
		return httpwrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return httpwrapper.NewResponse(http.StatusForbidden, nil, problemDetails)
}
