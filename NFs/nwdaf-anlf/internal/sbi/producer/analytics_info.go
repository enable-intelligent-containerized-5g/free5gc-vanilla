package producer

import (
	"math"
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

// SubtractSeconds subtracts seconds from a given date
func SubtractSeconds(date time.Time, seconds int64) time.Time {
    return date.Add(-time.Duration(seconds) * time.Second)
}


// Get ML by NfType, Size, Accuracy
func getMlModelByProfile(mlmodels *[]models.MlModelData, nftype *models.NfType, accuracy *models.NwdafMlModelAccuracy) (mlmodel []models.MlModelData) {

	if len(*mlmodels) == 0 || mlmodels == nil {
		logger.AniLog.Error("No Found MlModels")
		return nil
	}

	// Filter By NfType
	var amfModels []models.MlModelData
	for _, model := range *mlmodels {
		if model.NfType == *nftype {
			amfModels = append(amfModels, model)
		}
	}

	if len(amfModels) == 0 || amfModels == nil {
		return nil
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

	if len(smallestModels) == 0 || smallestModels == nil {
		return nil
	}

	// Filter by Accuracy
	definedAccuracies := []models.NwdafMlModelAccuracy{
		models.NwdafMlModelAccuracy_LOW,
		models.NwdafMlModelAccuracy_MEDIUM,
		models.NwdafMlModelAccuracy_HIGH,
	}
	// definedAccuracies := models.NewNwdafMlModelAccuracyPriority()

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

	if len(priorityModels) == 0 || priorityModels == nil {
		return nil
	}

	return priorityModels
}

func HandleAnalyticsInfoNfLoadMetrics(request *httpwrapper.Request, typePayload models.TypePayloadRequest) (response *httpwrapper.Response) {
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
	nfFilterdByTypePayload := []models.NfProfile{}
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
	// idsByIp := []string{}
	// ips := []string{}
	// for _, instance := range nfInstancesFilteredByIP {
	// 	idsByIp = append(idsByIp, instance.NfInstanceId) // Add Id
	// 	ips = append(ips, instance.Ipv4Addresses[0])     // Add Ip
	// }
	// logger.AniLog.Infof("Filtered NF Instances by IP(%d): %s, with IPs(%d): %s", len(idsByIp), idsByIp, len(ips), ips)

	if len(nfInstancesFilteredByIP) <= 0 {
		logger.AniLog.Warn("NFs not found")
		return httpwrapper.NewResponse(http.StatusNotFound, nil, "NFs not found")
	}

	// Get analisys (Predict or statistics)
	return GetAnaliticsNfLoadProcedure(&analyticsInfoDataRequest, &eventID, &NrfUri, &nfInstancesFilteredByIP)
}

func GetAnaliticsNfLoadProcedure(request *models.NwdafAnalyticsInfoRequest, eventID *models.EventId, NrfUri *string, nfInstances *[]models.NfProfile) *httpwrapper.Response {
	var responseNfLoad = models.NwdafAnalyticsInfoNfLoadResponse{}
	var analysisType models.AnalysisType

	// Get StartTime and EndTime and convert to UTC
	// example startTime in UTC-5: "2024-11-13T12:00:00-05:00",
	// example endTime in UTC-5: "2024-11-13T12:01:00-05:00",
	startTime := request.StartTime.UTC()
	endTime := request.EndTime.UTC()
	currentTime := time.Now()

	logger.AniLog.Infof("StartTime: %s, EndTime: %s, CurrentTime: %s", startTime, endTime, currentTime)

	var DefaultNfLoad = models.ResourcesNfLoad{
		CpuLoad: 0,
		MemLoad: 0,
	}

	var defaultValues = models.NwdafAnalyticsInfoNfLoad{
		CpuUsage: 0,
		MemUsage: 0,
		CpuLimit: 0,
		MemLimit: 0,
		NfLoad:   DefaultNfLoad,
	}

	// Convert time to seconds
	targetPeriod := parseTimeToSeconds(&startTime, &endTime)

	logger.AniLog.Infof("targetPeriod: %d", targetPeriod)
	if targetPeriod <= 0 {
		return httpwrapper.NewResponse(http.StatusBadRequest, nil, "EndTime must be greater than StartTime")
	}

	// Convert time to seconds
	offSet := parseTimeToSeconds(&endTime, &currentTime)

	logger.AniLog.Infof("offSet: %d", offSet)
	if targetPeriod <= 0 {
		return httpwrapper.NewResponse(http.StatusBadRequest, nil, "EndTime must be smaller than CurrentTime")
	}

	namespace := factory.NwdafConfig.Configuration.Namespace
	instancek8s := factory.NwdafConfig.Configuration.KsmInstance

	// Predict metrics
	switch {
	case endTime.After(currentTime):
		logger.AniLog.Info("Predict metrics: EndTime is greater than now")
		analysisType = models.AnalysisType_PREDICTIONS

		// Running Pods
		runningPods := consumer.GetRunningPods(instancek8s, namespace, "", currentTime)

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

		// mtlfUri = "http://127.0.0.1:4201"

		// Send GetMlModelInfoList
		err = consumer.SendGetMlModelInfoList(&mlModelInfoList, mtlfUri)
		if err != nil {
			logger.AniLog.Error("Error getting Ml Model Info: ", err)
			problemDetails := models.ProblemDetails{
				Title:  "Error getting Ml Model Info",
				Status: http.StatusInternalServerError,
				Detail: err.Error(),
			}
			return httpwrapper.NewResponse(http.StatusInternalServerError, nil, problemDetails)
		}
		if len(mlModelInfoList) <= 0 {
			logger.AniLog.Warn("Ml Model not found")
			return httpwrapper.NewResponse(http.StatusNotFound, nil, "Ml Model not found")
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
		NfLoadsAnalitics := []models.NwdafAnalyticsInfoNfLoad{}
		for _, profile := range *nfInstances {
			var NfLoad = models.NwdafAnalyticsInfoNfLoad(defaultValues)
			nfType := profile.NfType

			// Get ML by NfType, Size, Accuracy
			selectedModels := getMlModelByProfile(&mlModelInfoFiltered, &nfType, &request.Accuracy)

			if len(selectedModels) <= 0 || selectedModels == nil {
				logger.AniLog.Errorf("No Found a MlModel for the NfType %s with nfInstanceId %s", nfType, profile.NfInstanceId)
				continue
			}

			logger.AniLog.Infof("Found the MlModel %v for the NfType %s with nfInstanceId %s", selectedModels[0].URI, nfType, profile.NfInstanceId)

			var podName string
			// containerName := util.GetPodNameFromIpv4(profile.Ipv4Addresses[0])[0]
			containerName := profile.ContainerName

			foundPod := models.FindPodByContainer(runningPods, containerName)

			if foundPod != nil {
				podName = foundPod.Pod
			} else {
				logger.AniLog.Infof("No pod found for the specified container: %s", containerName)
				continue
			}

			// Get CPU and RAM  from Prometheus
			var numSamples int64 = 4
			newStartTime := SubtractSeconds(currentTime, targetPeriod*(numSamples-1)) // Subtarct secons to curenntime
			logger.AniLog.Warnf("numSamples: %d, newStartTime: %s, currentTime: %s", numSamples, newStartTime, currentTime)

			cpuUsageAverageRange := consumer.GetCpuUsageAverageRange(namespace, podName, containerName, targetPeriod, 0, newStartTime, currentTime)

			for _, mt := range cpuUsageAverageRange {
				// Convertir el timestamp a segundos
				// seconds := int64(ts.Timestamp / 1000)
				// Crear un objeto time.Time
				// t := time.Unix(seconds, 0)
				// Imprimir la fecha y hora en formato UTC
				// logger.AniLog.Warn("Time: ", ts.Timestamp, " -> ", t.UTC())
				logger.AniLog.Warn("Metric: ", mt)
			}
			// logger.AniLog.Warn("cpuUsageAverageRange: ", cpuUsageAverageRange)

			// Analize the CPU and RAM

			NfLoad = models.NwdafAnalyticsInfoNfLoad{
				NfInstanceId: profile.NfInstanceId,
				Accuracy:     selectedModels[0].Accuracy,
				NfType:       profile.NfType,
				Pod:          "",
				Container:    "",
				CpuUsage:     defaultValues.CpuUsage,
				MemUsage:     defaultValues.MemUsage,
				CpuLimit:     defaultValues.CpuLimit,
				MemLimit:     defaultValues.MemLimit,
				NfLoad:       defaultValues.NfLoad,
				NfStatus:     profile.NfStatus,
				Confidence:   selectedModels[0].Confidence.R2,
			}

			NfLoadsAnalitics = append(NfLoadsAnalitics, NfLoad)
		}

		responseNfLoad = models.NwdafAnalyticsInfoNfLoadResponse{
			EventId:         *eventID,
			AnalysisType:    analysisType,
			TargetPeriod:     targetPeriod,
			OffSet:          offSet,
			AnaliticsNfLoad: NfLoadsAnalitics,
		}

		// Return results
		return httpwrapper.NewResponse(http.StatusOK, nil, responseNfLoad)




	// Statistics metrics
	case startTime.Before(currentTime) && endTime.Before(currentTime):
		logger.AniLog.Info("Statistics metrics: EndTime is less than now")
		analysisType = models.AnalysisType_STATISTICS

		// Running Pods
		// logger.AniLog.Warn("Start time: ", startTime)
		runningPods := consumer.GetRunningPods(instancek8s, namespace, "", endTime)
		// logger.UtilLog.Warn("Running pods: ", runningPods)


		// For each profile: get data from Prometheus
		NfLoadsAnalitics := []models.NwdafAnalyticsInfoNfLoad{}
		for _, profile := range *nfInstances {
			var NfLoad = models.NwdafAnalyticsInfoNfLoad(defaultValues)
			var podName string
			// containerName := util.GetPodNameFromIpv4(profile.Ipv4Addresses[0])[0]
			containerName := profile.ContainerName

			foundPod := models.FindPodByContainer(runningPods, containerName)

			if foundPod != nil {
				podName = foundPod.Pod
			} else {
				logger.AniLog.Infof("No pod found for the specified container: %s", containerName)
				continue
			}

			logger.AniLog.Warn("End time: ", endTime)
			// logger.AniLog.Infof("NAMESPACE: %s,POD: %s, CONTAINER: %s", namespace, podName, containerName)
			cpuUsageAverage := consumer.GetCpuUsageAverage(namespace, podName, containerName, targetPeriod, 0, endTime)[0]
			memUsageAverage := consumer.GetMemUsageAverage(namespace, podName, containerName, targetPeriod, 0, endTime)[0]
			cpuLimit := consumer.GetResourceLimit(namespace, podName, containerName, models.PrometheusUnit_CORE, endTime)[0]
			memLimit := consumer.GetResourceLimit(namespace, podName, containerName, models.PrometheusUnit_BYTE, endTime)[0]

			// logger.AniLog.Infof("Cpu Usage: %f, MenUsage: %f, CpuLimit: %f, MemLimit: %f", cpuUsageAverage.Value, memUsageAverage.Value, cpuLimit.Value, memLimit.Value)

			logger.AniLog.Warn("Timestamp: ", cpuUsageAverage)

			var nfLoad = models.ResourcesNfLoad{
				CpuLoad: getPercentil(cpuUsageAverage.Value, cpuLimit.Value),
				MemLoad: getPercentil(memUsageAverage.Value, memLimit.Value),
			}

			NfLoad = models.NwdafAnalyticsInfoNfLoad{
				NfInstanceId: profile.NfInstanceId,
				Pod:          podName,
				Container:    containerName,
				NfType:       profile.NfType,
				CpuUsage:     cpuUsageAverage.Value,
				MemUsage:     memUsageAverage.Value,
				CpuLimit:     cpuLimit.Value,
				MemLimit:     memLimit.Value,
				NfLoad:       nfLoad,
				NfStatus:     profile.NfStatus,
			}

			NfLoadsAnalitics = append(NfLoadsAnalitics, NfLoad)
		}

		responseNfLoad = models.NwdafAnalyticsInfoNfLoadResponse{
			EventId:         *eventID,
			AnalysisType:    analysisType,
			TargetPeriod:     targetPeriod,
			OffSet:          offSet,
			AnaliticsNfLoad: NfLoadsAnalitics,
		}

		// Return results
		return httpwrapper.NewResponse(http.StatusOK, nil, responseNfLoad)

	default:
		logger.AniLog.Error("Invalid time range")
		return httpwrapper.NewResponse(http.StatusBadRequest, nil, "EndTime must be greater than StartTime")
	}
}

func getPercentil(value float64, limit float64) float64 {
	load := value / limit
	if math.IsNaN(value) || limit == 0 {
		load = 0
	}
	return load
}
