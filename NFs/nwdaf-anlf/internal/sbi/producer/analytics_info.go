package producer

import (
	"fmt"
	"net/http"

	"github.com/enable-intelligent-containerized-5g/openapi/Nnrf_NFDiscovery"
	"github.com/enable-intelligent-containerized-5g/openapi/models"
	"github.com/free5gc/nwdaf/internal/context"
	"github.com/free5gc/nwdaf/internal/logger"
	"github.com/free5gc/nwdaf/internal/sbi/consumer"
	"github.com/free5gc/util/httpwrapper"
)

func getNfLoadMetrics(request *httpwrapper.Request, typePayload models.TypePayloadRequest) *httpwrapper.Response {
	// Extract the body of the request
	analyticsInfoData := request.Body.(models.NwdafAnalyticsInfoRequest)

	logger.CfgLog.Info("Analytics Info Data: ", analyticsInfoData)

	// extract event values from analyticsInfoData
	eventID := analyticsInfoData.EventId

	logger.CfgLog.Info("Event ID: ", eventID)

	logger.CfgLog.Info("Type Payload Request: ", typePayload)

	// check the type of payload request
	switch typePayload {
	case models.TypePayloadRequest_NF_INSTANCES:
		// extract nfInstance values from analytics
		nfInstancesIds := analyticsInfoData.NfInstanceIds
		param := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
			// ServiceNames: optional.Interface{},
		}
		allNfs, err := consumer.SearchAllNfInstance("http://127.0.0.1:30050", "", models.NfType_NWDAF, param)
		if err != nil {
			return httpwrapper.NewResponse(http.StatusInternalServerError, nil, err)
		}
		logger.CfgLog.Info("All NFs: ", allNfs)
		logger.CfgLog.Info("NF Instances IDs: ", nfInstancesIds)
		fmt.Println("All NFs: ", allNfs)
		// return the response
		return httpwrapper.NewResponse(http.StatusOK, nil, allNfs)

	case models.TypePayloadRequest_NF_TYPES:
		// extract nfTypes values from analytics
		nfTypes := analyticsInfoData.NfTypes
		logger.CfgLog.Info("NF Types: ", nfTypes)
		// return the response
		return httpwrapper.NewResponse(http.StatusAccepted, nil, eventID)

	default:
		logger.CfgLog.Warn("Unknown type payload")
		return httpwrapper.NewResponse(http.StatusBadRequest, nil, "Unknown payload type")
	}
}

func HandleAnalyticsInfoNfLoadMetrics(request *httpwrapper.Request, typePayload models.TypePayloadRequest) *httpwrapper.Response {
	logger.CfgLog.Info("[NWDAF-AnLF] Handle Analytics Info NFLoad Metrics Request")

	// Extract the context NWDAF configuration
	GlobalNwdafContext := context.GlobalNwdafContext
	// Check if NRF URI is set
	if GlobalNwdafContext == nil {
		logger.CfgLog.Error("NWDAF context is not set")
		return httpwrapper.NewResponse(http.StatusInternalServerError, nil, "NWDAF context is not set")
	}

	// Extract the body of the request
	analyticsInfoData := request.Body.(models.NwdafAnalyticsInfoRequest)

	// extract event values from analyticsInfoData
	eventID := analyticsInfoData.EventId
	logger.CfgLog.Info("Event ID: ", eventID)
	// check the type of payload request
	switch typePayload {
	case models.TypePayloadRequest_NF_INSTANCES:
		// extract nfInstance values from analytics
		nfInstancesIds := analyticsInfoData.NfInstanceIds
		param := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
			// ServiceNames: optional.Interface{},
		}
		allNfs, err := consumer.SearchAllNfInstance(GlobalNwdafContext.NrfUri, "", models.NfType_NWDAF, param)
		if err != nil {
			return httpwrapper.NewResponse(http.StatusInternalServerError, nil, err)
		}
		logger.CfgLog.Info("All NFs: ", allNfs)
		logger.CfgLog.Info("NF Instances IDs: ", nfInstancesIds)
		fmt.Println("All NFs: ", allNfs)
		// return the response
		return httpwrapper.NewResponse(http.StatusOK, nil, allNfs)
	case models.TypePayloadRequest_NF_TYPES:
		// extract nfTypes values from analytics
		nfTypes := analyticsInfoData.NfTypes
		logger.CfgLog.Info("NF Types: ", nfTypes)
		// return the response
		return httpwrapper.NewResponse(http.StatusAccepted, nil, eventID)
	default:
		logger.CfgLog.Warn("Unknown type payload")
		return httpwrapper.NewResponse(http.StatusBadRequest, nil, "Unknown payload type")
	}
}
