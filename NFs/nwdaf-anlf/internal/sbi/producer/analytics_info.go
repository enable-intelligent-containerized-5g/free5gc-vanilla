package producer

import (
	"fmt"
	"net/http"
	"time"

	"github.com/enable-intelligent-containerized-5g/openapi/Nnrf_NFDiscovery"
	"github.com/enable-intelligent-containerized-5g/openapi/models"
	"github.com/free5gc/nwdaf/internal/logger"
	"github.com/free5gc/nwdaf/internal/sbi/consumer"
	"github.com/free5gc/nwdaf/pkg/factory"
	"github.com/free5gc/util/httpwrapper"
)

func getNfLoadMetrics(request *httpwrapper.Request, typePayload models.TypePayloadRequest) *httpwrapper.Response {
	// Extract the body of the request
	analyticsInfoData := request.Body.(models.NwdafAnalyticsInfoRequest)

	logger.AniLog.Info("Analytics Info Data: ", analyticsInfoData)

	// extract event values from analyticsInfoData
	eventID := analyticsInfoData.EventId

	logger.AniLog.Info("Event ID: ", eventID)

	logger.AniLog.Info("Type Payload Request: ", typePayload)

	// check the type of payload request
	switch typePayload {
	case models.TypePayloadRequest_NF_INSTANCES:
		// extract nfInstance values from analytics
		nfInstancesIds := analyticsInfoData.NfInstanceIds
		param := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
			// ServiceNames: optional.Interface{},
		}

		var nfInstances []models.NfProfile

		err := consumer.SearchAllNfInstance(&nfInstances, "http://127.0.0.1:30050", "", models.NfType_NWDAF, param)
		if err != nil {
			return httpwrapper.NewResponse(http.StatusInternalServerError, nil, err)
		}

		logger.AniLog.Info("All NFs: ", nfInstances)
		logger.AniLog.Info("NF Instances IDs: ", nfInstancesIds)
		fmt.Println("All NFs: ", nfInstances)
		// return the response
		return httpwrapper.NewResponse(http.StatusOK, nil, nfInstances)

	case models.TypePayloadRequest_NF_TYPES:
		// extract nfTypes values from analytics
		nfTypes := analyticsInfoData.NfTypes
		logger.AniLog.Info("NF Types: ", nfTypes)
		// return the response
		return httpwrapper.NewResponse(http.StatusAccepted, nil, eventID)

	default:
		logger.AniLog.Warn("Unknown type payload")
		return httpwrapper.NewResponse(http.StatusBadRequest, nil, "Unknown payload type")
	}
}

func HandleAnalyticsInfoNfLoadMetrics(request *httpwrapper.Request, typePayload models.TypePayloadRequest) *httpwrapper.Response {
	logger.AniLog.Info("[NWDAF-AnLF] Handle Analytics Info NFLoad Metrics Request")

	// Extract the context NWDAF configuration
	NrfUri := factory.NwdafConfig.Configuration.NrfUri

	// Check if NRF URI is set
	if NrfUri == "" {
		return httpwrapper.NewResponse(http.StatusInternalServerError, nil, "NrfUri is not set")
	}

	// Extract the body of the request
	analyticsInfoData := request.Body.(models.NwdafAnalyticsInfoRequest)

	// extract event values from analyticsInfoData
	eventID := analyticsInfoData.EventId
	logger.AniLog.Info("Event ID: ", eventID)
	// check the type of payload request
	switch typePayload {
	case models.TypePayloadRequest_NF_INSTANCES:
		// extract nfInstance values from analytics
		nfInstancesIds := analyticsInfoData.NfInstanceIds
		param := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
			// ServiceNames: optional.Interface{},
		}

		// Var to store all NF instances
		var nfInstances []models.NfProfile

		// Search all NF instances
		err := consumer.SearchAllNfInstance(&nfInstances, NrfUri, "", models.NfType_NWDAF, param)
		if err != nil {
			return httpwrapper.NewResponse(http.StatusInternalServerError, nil, err)
		}

		// Get StartTime and EndTime
		startTime := analyticsInfoData.StartTime
		endTime := analyticsInfoData.EndTime
		currentTime := time.Now()

		// Predict metrics
		if endTime.After(currentTime) || startTime.Before(currentTime) {
			fmt.Println("EndTime is greater than now")

			var mtlfUri string

      // Search MLF URI
			err := consumer.SearchMlModelInfoInstance(&mtlfUri, NrfUri, models.NfType_NWDAF, models.NfType_NWDAF, param)
			if err != nil {
				logger.AniLog.Error("MLF URI not found: ", err)
				return httpwrapper.NewResponse(http.StatusInternalServerError, nil, err)
			}

      logger.AniLog.Info("MLF URI: ", mtlfUri)

      var mlModelInfoList []models.MlModelData

      // Send GetMlModelInfoList
      err = consumer.SendGetMlModelInfoList(&mlModelInfoList, mtlfUri)
      if err != nil {
        logger.AniLog.Error("ML Model Info List not found: ", err)
        return httpwrapper.NewResponse(http.StatusInternalServerError, nil, err)
      }

      return httpwrapper.NewResponse(http.StatusOK, nil, mlModelInfoList)
		}

		// Analyze metrics existence
		if startTime.Before(currentTime) || endTime.Before(currentTime) {
			fmt.Println("EndTime is less than now")
		}

		logger.AniLog.Info("All NFs: ", nfInstances)
		logger.AniLog.Info("NF Instances IDs: ", nfInstancesIds)
		fmt.Println("All NFs: ", nfInstances)
		// return the response
		return httpwrapper.NewResponse(http.StatusOK, nil, nfInstances)
	case models.TypePayloadRequest_NF_TYPES:
		// extract nfTypes values from analytics
		nfTypes := analyticsInfoData.NfTypes
		logger.AniLog.Info("NF Types: ", nfTypes)
		// return the response
		return httpwrapper.NewResponse(http.StatusAccepted, nil, eventID)
	default:
		logger.AniLog.Warn("Unknown type payload")
		return httpwrapper.NewResponse(http.StatusBadRequest, nil, "Unknown payload type")
	}
}
