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

func getNfLoadMetrics() {
}

func filterNfInstanceById(nfIntances *[]models.NfProfile, nfInstanceIds []string) (nfInstancesFiltered []models.NfProfile) {
	for _, nfInstance := range *nfIntances {
		for _, nfInstanceId := range nfInstanceIds {
			if nfInstance.NfInstanceId == nfInstanceId {
				nfInstancesFiltered = append(nfInstancesFiltered, nfInstance)
			}
		}
	}
	return nfInstancesFiltered
}

func filterNfInstanceByNfType(nfIntances *[]models.NfProfile, nfTypes []models.NfType) (nfInstancesFiltered []models.NfProfile) {
	for _, nfInstance := range *nfIntances {
		for _, nfType := range nfTypes {
			if nfInstance.NfType == nfType {
				nfInstancesFiltered = append(nfInstancesFiltered, nfInstance)
			}
		}
	}
	return nfInstancesFiltered
}

func filterMlModelInfo(mlModelInfoList *[]models.MlModelData, eventId *models.EventId, accuracy *models.NwdafMlModelAccuracy, targetPeriod int64) (mlModelInfoFiltered []models.MlModelData) {
	for _, mlModelInfo := range *mlModelInfoList {
		if mlModelInfo.EventId == *eventId && mlModelInfo.Accuracy == *accuracy && mlModelInfo.TargetPeriod == targetPeriod {
			mlModelInfoFiltered = append(mlModelInfoFiltered, mlModelInfo)
		}
	}
	return mlModelInfoFiltered
}

func parseTimeToSeconds(startTime *time.Time, endTime *time.Time) int64 {
	startTimeUnix := startTime.Unix()
	endTimeUnix := endTime.Unix()
	return endTimeUnix - startTimeUnix
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

		// Filter NF instances by ID
		nfInstancesFiltered := filterNfInstanceById(&nfInstances, nfInstancesIds)

		logger.AniLog.Info("Filtered NF Instances: ", nfInstancesFiltered)

		// Get StartTime and EndTime
		startTime := analyticsInfoData.StartTime
		endTime := analyticsInfoData.EndTime
		currentTime := time.Now()

		// Predict metrics
		switch {
		case endTime.After(currentTime) || startTime.Before(currentTime):
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

			// Convert time to seconds
			targetPeriod := parseTimeToSeconds(startTime, endTime)

			// Filter ML Model Info
			mlModelInfoFiltered := filterMlModelInfo(&mlModelInfoList, &eventID, &analyticsInfoData.Accuracy, targetPeriod)
			if mlModelInfoFiltered == nil {
				return httpwrapper.NewResponse(http.StatusInternalServerError, nil, "ML Model Info not found")
			}

			logger.AniLog.Info("Filtered ML Model Info: ", mlModelInfoFiltered)

			return httpwrapper.NewResponse(http.StatusOK, nil, mlModelInfoFiltered)
		case startTime.Before(currentTime) || endTime.Before(currentTime):
			fmt.Println("EndTime is less than now")
			return httpwrapper.NewResponse(http.StatusOK, nil, nfInstances)
		default:
			fmt.Println("invalid time range")
			return httpwrapper.NewResponse(http.StatusBadRequest, nil, "Invalid time range")
		}

		// return the response
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
