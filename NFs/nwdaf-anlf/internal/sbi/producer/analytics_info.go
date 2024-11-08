package producer

import (
	"net/http"
	"time"

	"github.com/enable-intelligent-containerized-5g/openapi/Nnrf_NFDiscovery"
	"github.com/enable-intelligent-containerized-5g/openapi/models"
	"github.com/free5gc/nwdaf/internal/logger"
	"github.com/free5gc/nwdaf/internal/sbi/consumer"
	"github.com/free5gc/nwdaf/pkg/factory"
	"github.com/free5gc/util/httpwrapper"
)

type DefaultNfLoad struct {
	CpuUsage float64 `json:"cpu_usage,omitempty" yaml:"cpu_usage,omitempty" bson:"cpu_usage,omitempty" `
	MemUsage float64 `json:"mem_usage,omitempty" yaml:"mem_usage,omitempty" bson:"mem_usage,omitempty" `
	CpuLimit float64 `json:"cpu_limit,omitempty" yaml:"cpu_limit,omitempty" bson:"cpu_limit,omitempty" `
	MemLimit float64 `json:"mem_limit,omitempty" yaml:"mem_limit,omitempty" bson:"mem_limit,omitempty" `
	NfLoad   float64 `json:"nf_load,omitempty" yaml:"nf_load,omitempty" bson:"nf_load,omitempty" `
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

func filterMlModelInfo(mlModelInfoList *[]models.MlModelData, eventId *models.EventId, targetPeriod int64) (mlModelInfoFiltered []models.MlModelData) {
	for _, mlModelInfo := range *mlModelInfoList {
		if mlModelInfo.EventId == *eventId && mlModelInfo.TargetPeriod == targetPeriod {
			mlModelInfoFiltered = append(mlModelInfoFiltered, mlModelInfo)
		}
	}

	return mlModelInfoFiltered
}

func filterNfInstancesWithIpDuplicate(nfInstances *[]models.NfProfile) (mlNfProfileFiltered []models.NfProfile) {
	ipTracker := make(map[string]bool)

	// Filter the NfInstances
	for _, instance := range *nfInstances {
		if _, exists := ipTracker[instance.Ipv4Addresses[0]]; !exists {
			// If no exits IP ADD instance
			mlNfProfileFiltered = append(mlNfProfileFiltered, instance)
			ipTracker[instance.Ipv4Addresses[0]] = true
		}
	}

	return mlNfProfileFiltered
}

func parseTimeToSeconds(startTime *time.Time, endTime *time.Time) int64 {
	startTimeUnix := startTime.Unix()
	endTimeUnix := endTime.Unix()
	return endTimeUnix - startTimeUnix
}

// Get ML by NfType, Size, Accuracy
func getMlModelByProfile(mlmodels *[]models.MlModelData, nftype *models.NfType, accuracy *models.NwdafMlModelAccuracy) (mlmodel []models.MlModelData) {

	// Filter By NfType
	var amfModels []models.MlModelData
	for _, model := range *mlmodels {
		if model.NfType == *nftype {
			amfModels = append(amfModels, model)
		}
	}

	// Select the smallest models
	// find the smallestSize
	minSize := amfModels[0].Size
	for _, model := range amfModels {
		if model.Size < minSize {
			minSize = model.Size
		}
	}
	// logger.AniLog.Info("minSize: ", minSize)
	var smallestModels []models.MlModelData
	for _, model := range amfModels {
		if model.Size == minSize {
			smallestModels = append(smallestModels, model)
		}
	}

	// Filter by Accuracy
	definedAccuracies := []models.NwdafMlModelAccuracy{
		models.NwdafMlModelAccuracy_LOW,
		models.NwdafMlModelAccuracy_MEDIUM,
		models.NwdafMlModelAccuracy_HIGH,
	}
	if accuracy != nil && *accuracy != "" {
		// logger.AniLog.Info("Custom Acuracy: ", *accuracy)
		requestAccuracy := *accuracy

		for i, v := range definedAccuracies {
			if v == requestAccuracy {
				// Move the requestAccuracy to he begining
				definedAccuracies = append([]models.NwdafMlModelAccuracy{v}, append(definedAccuracies[:i], definedAccuracies[i+1:]...)...)
				break
			}
		}
	}
	// Search models by accuracy priority
	var priorityModels []models.MlModelData
	for _, priority := range definedAccuracies {
		for _, model := range smallestModels {
			if model.Accuracy == priority {
				priorityModels = append(priorityModels, model)
			}
		}
		// If we find models with the current priority, we exit the loop.
		if len(priorityModels) > 0 {
			break
		}
	}

	return priorityModels
}

func HandleAnalyticsInfoNfLoadMetricsNew(request *httpwrapper.Request, typePayload models.TypePayloadRequest) (response *httpwrapper.Response) {
	logger.AniLog.Info("Handle Analytics Info NFLoad Metrics Request")

	// Extract the context NWDAF configuration
	NrfUri := factory.NwdafConfig.Configuration.NrfUri
	// Extract the body of the request
	analyticsInfoDataRequest := request.Body.(models.NwdafAnalyticsInfoRequest)
	// extract event values from analyticsInfoDataRequest
	eventID := analyticsInfoDataRequest.EventId
	logger.AniLog.Info("Event ID: ", eventID)

	// Check if NRF URI is set
	if NrfUri == "" {
		return httpwrapper.NewResponse(http.StatusInternalServerError, nil, "NrfUri is not set")
	}

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

	// check the type of payload request
	var nfFilterdByTypePayload = []models.NfProfile{}
	switch typePayload {
	case models.TypePayloadRequest_NF_INSTANCES:
		logger.AniLog.Infof("typePayload: %s", models.TypePayloadRequest_NF_INSTANCES)

		// extract nfInstance values from analytics
		nfInstancesIds := analyticsInfoDataRequest.NfInstanceIds

		// Filter NF instances by ID
		nfFilterdByTypePayload = filterNfInstanceById(&nfInstances, nfInstancesIds)
		// idsById := []string{}
		// for _, instance := range nfFilterdByTypePayload {
		// 	idsById = append(idsById, instance.NfInstanceId) // Add the Id
		// }
		// logger.AniLog.Infof("Filtered NF Instances by ID(%d): %s\n", len(idsById), idsById)

	case models.TypePayloadRequest_NF_TYPES:
		logger.AniLog.Infof("typePayload: %s", models.TypePayloadRequest_NF_TYPES)

		// extract nfTypes values from analytics
		nfTypes := analyticsInfoDataRequest.NfTypes
		logger.AniLog.Info("NF Types: ", nfTypes)

		// Filter NF instances by NfType
		nfFilterdByTypePayload = filterNfInstanceByNfType(&nfInstances, nfTypes)
		// idsByNfType := []string{}
		// for _, instance := range nfFilterdByTypePayload {
		// 	idsByNfType = append(idsByNfType, instance.NfInstanceId) // Add the Id
		// }
		// logger.AniLog.Infof("Filtered NF Instances by NfType(%d): %s", len(idsByNfType), idsByNfType)

		// return the response
		// return httpwrapper.NewResponse(http.StatusAccepted, nil, fmt.Sprintf("EventId: %s,  NfIds: %s", eventID, idsByNfType))

	default:
		logger.AniLog.Warn("Unknown type payload")
		return httpwrapper.NewResponse(http.StatusBadRequest, nil, "Unknown payload type")
	}

	// Filter NF instances by IP
	nfInstancesFilteredByIP := filterNfInstancesWithIpDuplicate(&nfFilterdByTypePayload)
	idsByIp := []string{}
	ips := []string{}
	for _, instance := range nfInstancesFilteredByIP {
		idsByIp = append(idsByIp, instance.NfInstanceId) // Add Id
		ips = append(ips, instance.Ipv4Addresses[0])     // Add Ip
	}
	logger.AniLog.Infof("Filtered NF Instances by IP(%d): %s, with IPs(%d): %s", len(idsByIp), idsByIp, len(ips), ips)

	if len(nfInstancesFilteredByIP) <= 0 {
		logger.AniLog.Warn("NFs not found")
		return httpwrapper.NewResponse(http.StatusNotFound, nil, "NFs not found")
	}

	// Get analisys (Predict or stadistics)
	return GetAnaliticsMetrics(&analyticsInfoDataRequest, &eventID, &NrfUri, &nfInstancesFilteredByIP)
}

func GetAnaliticsMetrics(request *models.NwdafAnalyticsInfoRequest, eventID *models.EventId, NrfUri *string, nfInstances *[]models.NfProfile) *httpwrapper.Response {
	var responseNfLoad = models.NwdaAnalyticsInfoNfLoadResponse{}
	// Get StartTime and EndTime
	startTime := request.StartTime
	endTime := request.EndTime
	currentTime := time.Now()

	var defaultValues = DefaultNfLoad{
		CpuUsage: 0.3,
		MemUsage: 320,
		CpuLimit: 0.5,
		MemLimit: 350,
		NfLoad:   87.5,
	}

	// Predict metrics
	switch {
	case endTime.After(currentTime):
		logger.AniLog.Info("Predict metrics: EndTime is greater than now")

		var mtlfUri string

		// Search MTLF URI
		param := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
			// ServiceNames: optional.Interface{},
		}
		err := consumer.SearchMlModelInfoInstance(&mtlfUri, *NrfUri, models.NfType_NWDAF, models.NfType_NWDAF, param)
		if err != nil {
			logger.AniLog.Error("MTLF URI not found: ", err)
			return httpwrapper.NewResponse(http.StatusInternalServerError, nil, err)
		}

		logger.AniLog.Info("MTLF URI: ", mtlfUri)

		var mlModelInfoList []models.MlModelData

		// Send GetMlModelInfoList
		err = consumer.SendGetMlModelInfoList(&mlModelInfoList, mtlfUri)
		if err != nil {
			logger.AniLog.Error("Error getting Ml Model Info: ", err)
			return httpwrapper.NewResponse(http.StatusInternalServerError, nil, err)
		}

		if len(mlModelInfoList) <= 0 {
			logger.AniLog.Warn("Ml Model not found")
			return httpwrapper.NewResponse(http.StatusNotFound, nil, "Ml Model not found")
		}

		// Convert time to seconds
		targetPeriod := parseTimeToSeconds(startTime, endTime)

		logger.AniLog.Infof("targetPeriod: %d", targetPeriod)
		if targetPeriod <= 0 {
			return httpwrapper.NewResponse(http.StatusBadRequest, nil, "EndTime must be greater than StartTime")
		}

		// Filter ML Model Info
		// logger.AniLog.Infof("ModelInfoList: %v", mlModelInfoList)
		// logger.AniLog.Infof("filterMlModelInfo Params ->   eventID: %s, targetPeriod: %d", string(*eventID), targetPeriod)
		mlModelInfoFiltered := filterMlModelInfo(&mlModelInfoList, eventID, targetPeriod)
		if mlModelInfoFiltered == nil {
			logger.AniLog.Info("No Found MlModels for predictions")
			return httpwrapper.NewResponse(http.StatusNotFound, nil, "ML Model Info not found for predictions")
		}
		// logger.AniLog.Infof("Filtered ML Model Info(%d): %v", len(mlModelInfoFiltered), mlModelInfoFiltered)

		// For each profile: get ml model, and get analitics
		var NfLoadsAnalitics = []models.NwdaAnalyticsInfoNfLoad{}
		for _, profile := range *nfInstances {
			var NfLoad = models.NwdaAnalyticsInfoNfLoad{}
			nfType := profile.NfType

			// Get ML by NfType, Size, Accuracy
			selectedModels := getMlModelByProfile(&mlModelInfoFiltered, &nfType, &request.Accuracy)

			if len(selectedModels) >= 0 {
				logger.AniLog.Infof("Found the MlModel %v for NfType %s with nfInstanceId %s", selectedModels[0].URI, nfType, profile.NfInstanceId)

				// Get CPU and RAM  from Prometheus

				// Analize the CPU and RAM

				NfLoad = models.NwdaAnalyticsInfoNfLoad{
					NfInstanceId: profile.NfInstanceId,
					Accuracy:     selectedModels[0].Accuracy,
					NfType:       profile.NfType,
					CpuUsage:     defaultValues.CpuUsage,
					MemUsage:     defaultValues.MemUsage,
					CpuLimit:     defaultValues.CpuLimit,
					MemLimit:     defaultValues.MemLimit,
					NfLoad:       defaultValues.NfLoad,
					NfStatus:     profile.NfStatus,
					Confidence:   selectedModels[0].Confidence,
				}

				NfLoadsAnalitics = append(NfLoadsAnalitics, NfLoad)

			} else {
				logger.AniLog.Errorf("No Found a MlModel for NfType %s with nfInstanceId %s", nfType, profile.NfInstanceId)
			}

		}

		responseNfLoad = models.NwdaAnalyticsInfoNfLoadResponse{
			AnaliticsNfLoad: NfLoadsAnalitics,
			EventId:         *eventID,
		}

		// Return results
		return httpwrapper.NewResponse(http.StatusOK, nil, responseNfLoad)

	// Stadistics metrics
	case startTime.Before(currentTime) && endTime.Before(currentTime):
		logger.AniLog.Info("Stadistics metrics: EndTime is less than now")

		// For each profile: get data from Prometheus
		var NfLoadsAnalitics = []models.NwdaAnalyticsInfoNfLoad{}
		for _, profile := range *nfInstances {
			var NfLoad = models.NwdaAnalyticsInfoNfLoad{}

			// Get CPU and RAM  from Prometheus

			NfLoad = models.NwdaAnalyticsInfoNfLoad{
				NfInstanceId: profile.NfInstanceId,
				// Accuracy: selectedModels[0].Accuracy,
				NfType:   profile.NfType,
				CpuUsage: defaultValues.CpuUsage,
				MemUsage: defaultValues.MemUsage,
				CpuLimit: defaultValues.CpuLimit,
				MemLimit: defaultValues.MemLimit,
				NfLoad:   defaultValues.NfLoad,
				NfStatus: profile.NfStatus,
				// Confidence: selectedModels[0].Confidence,
			}

			NfLoadsAnalitics = append(NfLoadsAnalitics, NfLoad)
		}

		responseNfLoad = models.NwdaAnalyticsInfoNfLoadResponse{
			AnaliticsNfLoad: NfLoadsAnalitics,
			EventId:         *eventID,
		}

		// Return results
		return httpwrapper.NewResponse(http.StatusOK, nil, responseNfLoad)

	default:
		logger.AniLog.Error("Invalid time range")
		return httpwrapper.NewResponse(http.StatusBadRequest, nil, "EndTime must be greater than StartTime")
	}
}
