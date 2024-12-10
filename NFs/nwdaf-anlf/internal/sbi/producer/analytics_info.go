package producer

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/enable-intelligent-containerized-5g/openapi/Nnrf_NFDiscovery"
	pcm "github.com/enable-intelligent-containerized-5g/openapi/PacketCaptureModule"
	pcm_models "github.com/enable-intelligent-containerized-5g/openapi/PacketCaptureModule/pcm_models"
	"github.com/enable-intelligent-containerized-5g/openapi/models"
	nwdaf_util "github.com/enable-intelligent-containerized-5g/openapi/nwdaf/util"
	"github.com/free5gc/nwdaf/internal/logger"
	"github.com/free5gc/nwdaf/internal/sbi/consumer"
	"github.com/free5gc/nwdaf/internal/util"
	"github.com/free5gc/nwdaf/pkg/factory"
	"github.com/free5gc/util/httpwrapper"
)

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
		problemDetails := models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Detail: "NrfUri is not set",
		}
		logger.AniLog.Errorf(problemDetails.Detail)
		return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}

	param := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		// ServiceNames: optional.Interface{},
	}
	// Var to store all NF instances
	var nfInstances []models.NfProfile
	// Search all NF instances
	err := consumer.SearchAllNfInstance(&nfInstances, NrfUri, "", models.NfType_NWDAF, param)
	if err != nil {
		problemDetails := models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Detail: fmt.Sprintf("Error geting NfProfiles: %s", err.Error()),
		}
		logger.AniLog.Errorf(problemDetails.Detail)
		return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
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

	case models.TypePayloadRequest_NF_TYPES:
		logger.AniLog.Infof("typePayload: %s", models.TypePayloadRequest_NF_TYPES)

		// extract nfTypes values from analytics
		nfTypes := analyticsInfoDataRequest.NfTypes
		logger.AniLog.Info("NF Types: ", nfTypes)

		// Filter NF instances by NfType
		nfFilterdByTypePayload = filterNfInstanceByNfType(&nfInstances, nfTypes)

	default:
		problemDetails := models.ProblemDetails{
			Status: http.StatusBadRequest,
			Detail: fmt.Sprintf("Unknown payload type: the %s is not a valid payload", typePayload),
		}
		logger.AniLog.Errorf(problemDetails.Detail)
		return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}

	// Filter NF instances by IP
	nfInstancesFilteredByIP := filterNfInstancesWithIpDuplicate(&nfFilterdByTypePayload)

	if len(nfInstancesFilteredByIP) <= 0 {
		problemDetails := models.ProblemDetails{
			Status: http.StatusNotFound,
			Detail: "Error filtering NFs: NFs not found in " + NrfUri,
		}
		logger.AniLog.Errorf(problemDetails.Detail)
		return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
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

	// logger.AniLog.Infof("StartTime: %s, EndTime: %s, CurrentTime: %s", startTime, endTime, currentTime)

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
	targetPeriod := nwdaf_util.ParseTimeToSeconds(&startTime, &endTime)
	// Check the TargetPeriod
	if targetPeriod < 60 || targetPeriod <= 0 {
		problemDetails := models.ProblemDetails{
			Status: http.StatusBadRequest,
			Detail: "The difference between the start date and the end date must be greater than 60 seconds",
		}
		logger.AniLog.Errorf(problemDetails.Detail)
		return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}

	// Convert time to seconds
	offSet := nwdaf_util.ParseTimeToSeconds(&endTime, &currentTime)

	pcmUri := factory.NwdafConfig.Configuration.OamUri
	namespace := factory.NwdafConfig.Configuration.Namespace
	instancek8s := factory.NwdafConfig.Configuration.KsmInstance

	// Predict metrics
	var errList []error
	switch {
	case endTime.After(currentTime):
		logger.AniLog.Info("Predict metrics: EndTime is greater than now")
		analysisType = models.AnalysisType_PREDICTIONS

		// Running Pods
		runningPods, errPods := pcm.GetRunningPods(instancek8s, namespace, "", currentTime, pcmUri)
		if errPods != nil {
			problemDetails := models.ProblemDetails{
				Status: http.StatusInternalServerError,
				Detail: fmt.Sprintf("Error getting running pods from Packet Capture module: %s", errPods.Error()),
			}
			logger.AniLog.Errorf(problemDetails.Detail)
			return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
		}

		var mtlfUri string

		// Search MTLF URI
		param := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
			// ServiceNames: optional.Interface{},
		}
		err := consumer.SearchMlModelInfoInstance(&mtlfUri, *NrfUri, models.NfType_NWDAF, models.NfType_NWDAF, param)
		if err != nil {
			problemDetails := models.ProblemDetails{
				Status: http.StatusInternalServerError,
				Detail: fmt.Sprintf("MTLF URI not found: %s", err),
			}
			logger.AniLog.Errorf(problemDetails.Detail)
			return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
		}

		logger.AniLog.Info("MTLF URI: ", mtlfUri)

		var mlModelInfoList []models.MlModelData

		// Send GetMlModelInfoList
		err = consumer.SendGetMlModelInfoList(&mlModelInfoList, mtlfUri)
		if err != nil {
			problemDetails := models.ProblemDetails{
				Status: http.StatusInternalServerError,
				Detail: fmt.Sprintf("Error getting Ml Model Info: %s", err.Error()),
			}
			logger.AniLog.Error(problemDetails.Detail)
			return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
		}
		if len(mlModelInfoList) <= 0 {
			problemDetails := models.ProblemDetails{
				Status: http.StatusNotFound,
				Detail: "Ml Model not found",
			}
			logger.AniLog.Error(problemDetails.Detail)
			return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
		}

		mlModelInfoFiltered := filterMlModelInfo(&mlModelInfoList, eventID, targetPeriod)

		if mlModelInfoFiltered == nil {
			problemDetails := models.ProblemDetails{
				Status: http.StatusNotFound,
				Detail: "No Found MlModels for the predictions",
			}
			logger.AniLog.Error(problemDetails.Detail)
			return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
		}

		// For each profile: get ml model, and get analitics
		NfLoadsAnalitics := []models.NwdafAnalyticsInfoNfLoad{}
		for _, profile := range *nfInstances {
			var NfLoad = models.NwdafAnalyticsInfoNfLoad(defaultValues)
			nfType := profile.NfType

			// Get ML by NfType, Size, Accuracy
			selectedModels := getMlModelByProfile(&mlModelInfoFiltered, &nfType, &request.Accuracy)

			if len(selectedModels) <= 0 || selectedModels == nil {
				err := fmt.Errorf(" No Found a MlModel for the NfType %s with nfInstanceId %s", nfType, profile.NfInstanceId)
				logger.AniLog.Errorf(err.Error())
				errList = append(errList, err)

				continue
			}

			logger.AniLog.Infof("Found the MlModel %v for the NfType %s with nfInstanceId %s", selectedModels[0].URI, nfType, profile.NfInstanceId)

			var podName string
			containerName := profile.ContainerName

			foundPod := pcm_models.FindPodByContainer(runningPods, containerName)

			if foundPod != nil {
				podName = foundPod.Pod
			} else {
				err := fmt.Errorf(" No pod found for the specified container: %s", containerName)
				logger.AniLog.Infof(err.Error())
				errList = append(errList, err)

				continue
			}

			// Get CPU and RAM  from Prometheus
			var numSamples int64 = 4
			newStartTime := nwdaf_util.SubtractSeconds(currentTime, targetPeriod*(numSamples-1)) // Subtarct secons to curenntime
			logger.AniLog.Warnf("numSamples: %d, newStartTime: %s, currentTime: %s", numSamples, newStartTime, currentTime)

			cpuUsageAverageRange, errCpu := pcm.GetCpuUsageAverageRange(namespace, podName, containerName, targetPeriod, 0, newStartTime, currentTime, pcmUri)
			memUsageAverageRange, errMem := pcm.GetMemUsageAverageRange(namespace, podName, containerName, targetPeriod, 0, newStartTime, currentTime, pcmUri)
			totalThroughputRange, errtotalThrougput := pcm.GetThroughputAverageRange(namespace, podName, targetPeriod, 0, pcm_models.MetricType_TOTAL_THROUGPUT_AVERAGE, newStartTime, currentTime, pcmUri)
			cpuLimit, errLimCpu := pcm.GetResourceLimit(namespace, podName, containerName, pcm_models.PrometheusUnit_CORE, currentTime, pcmUri)
			memLimit, errLimMem := pcm.GetResourceLimit(namespace, podName, containerName, pcm_models.PrometheusUnit_BYTE, currentTime, pcmUri)

			if errCpu != nil || errMem != nil || errLimCpu != nil || errLimMem != nil || errtotalThrougput != nil {
				problemDetails := models.ProblemDetails{
					Status: http.StatusInternalServerError,
					Detail: fmt.Sprintf("Error getting data from Packet capture module: %s, %s, %s, %s", errCpu, errMem, errLimCpu, errLimMem),
				}
				logger.AniLog.Errorf(problemDetails.Detail)
				return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
			}

			logger.AniLog.Info("Saving data")
			cpuLimitValue := cpuLimit[0]
			memLimitValue := memLimit[0]
			nwdaf_util.DivideValues(&cpuUsageAverageRange, cpuLimitValue.Value)
			nwdaf_util.DivideValues(&memUsageAverageRange, memLimitValue.Value)
			pcm_models.UpdateContainerNameInPrometheusResultList(&totalThroughputRange, containerName)

			// Data paths
			dataPath := util.NwdafDefaultDataPath
			dataRawPath := util.NwdafDefaultDataRawPath
			menUsageFile := util.NwdafDefaultMenUsageFile
			cpuUsageFile := util.NwdafDefaultCpuUsageFile
			totalThroughputFile := util.NwdafDefaultTotalThroughputFile

			// Save Cpu usaje in a JSON file
			pathCpuUsage := dataRawPath + cpuUsageFile
			errToCsvCpu := nwdaf_util.SaveToJson(pathCpuUsage, cpuUsageAverageRange)
			if errToCsvCpu != nil {
				logger.AniLog.Error("Error: ", errToCsvCpu)
			} else {
				logger.AniLog.Infof("CpuUsage saved in %s (%d rows)", pathCpuUsage, len(cpuUsageAverageRange))
			}

			// Save Memory usage in a JSON file
			pathMemUsage := dataRawPath + menUsageFile
			errToCsvMem := nwdaf_util.SaveToJson(pathMemUsage, memUsageAverageRange)
			if errToCsvMem != nil {
				logger.AniLog.Error("Error: ", errToCsvMem)
			} else {
				logger.AniLog.Infof("MemUsage saved in %s (%d rows)", pathMemUsage, len(memUsageAverageRange))
			}

			// Llamar a la función para escribir el JSON
			pathTotalThroughput := dataRawPath + totalThroughputFile
			errToCsvThroughput := nwdaf_util.SaveToJson(pathTotalThroughput, totalThroughputRange)
			if errToCsvThroughput != nil {
				logger.MlModelTrainingLog.Error("Error: ", errToCsvThroughput)
			} else {
				logger.MlModelTrainingLog.Infof("CpuUsage saved in %s (%d rows)", pathTotalThroughput, len(totalThroughputRange))
			}

			// Processing data
			logger.AniLog.Info("Processing data")
			cpuColumn := string(pcm_models.MetricType_CPU_USAGE_AVERAGE)
			memColumn := string(pcm_models.MetricType_MEMORY_USAGE_AVERAGE)
			thrptColumn := string(pcm_models.MetricType_TOTAL_THROUGPUT_AVERAGE)
			pathDataProcessingScript := util.NwdafDefaultDataProcessingScriptPath
			dataPreprocessedPath := util.NwdafDefaultDataPreprocessedPath
			dataProcessedPath := util.NwdafDefaultDataProcessedPath
			dataLabeledPath := util.NwdafDefaultDataLabeledPath

			// Build the datasetName
			baseName := fmt.Sprintf("%s_%s_%ds", *eventID, nfType, targetPeriod)
			datasetFile := fmt.Sprintf("dataset_%s.csv", baseName)

			// Run processing data script
			cmd := exec.Command("python3", pathDataProcessingScript, dataPath,
				dataRawPath, dataPreprocessedPath,
				dataProcessedPath, dataLabeledPath,
				cpuUsageFile, menUsageFile, totalThroughputFile, datasetFile,
				cpuColumn, memColumn, thrptColumn)

			// Get the output and error
			outputProcess, errProcess := cmd.CombinedOutput()
			if errProcess != nil {
				problemDetails := &models.ProblemDetails{
					Status: http.StatusInternalServerError,
					Detail: fmt.Sprintf("Error in processing data %s. %s", *eventID, string(outputProcess)),
				}
				logger.AniLog.Error(problemDetails.Detail)
				return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
			}
			logger.AniLog.Infof("Data processing completed and saved in: %v", dataLabeledPath+datasetFile)

			// Analize the CPU and RAM

			logger.AniLog.Infof("predicting %s", *eventID)
			timeSteps := factory.NwdafConfig.Configuration.MlModelTrainingInfo.TimeSteps
			selectedModel := selectedModels[0]
			selectedModelUri := selectedModel.URI
			modelPredictionScriptPath := util.NwdafDefaultModelPredictionScriptPath
			modelsPath := util.NwdafDefaultModelsPath
			predictionsFile := util.NwdafDefaultModelPredictionFile
			// Run prediction script
			cmdPredicting := exec.Command("python3", modelPredictionScriptPath,
				modelsPath, dataPath, dataLabeledPath, datasetFile,
				predictionsFile, cpuColumn, memColumn, thrptColumn, selectedModelUri, strconv.FormatInt(timeSteps, 10))
			// Get the output and error
			outputPrediction, errTraining := cmdPredicting.CombinedOutput()
			if errTraining != nil {
				problemDetails := &models.ProblemDetails{
					Status: http.StatusInternalServerError,
					Detail: fmt.Sprintf("Error in predicting %s. %s: %s", *eventID, string(errTraining.Error()), outputPrediction),
				}
				logger.AniLog.Error(problemDetails.Detail)
				return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
			}
			if strings.TrimSpace(string(outputPrediction)) != "" {
				logger.AniLog.Warn(string(outputPrediction))
			}
			logger.AniLog.Infof("%s prediction completed and saved in %s", *eventID, dataPath+predictionsFile)

			// Load the predictions
			var nfLoadPred models.PredictionResult

			errLoadPrediction := loadPredictionInfoFromJson(&nfLoadPred, dataPath+util.NwdafDefaultModelPredictionFile)
			if errLoadPrediction != nil {
				problemDetails := &models.ProblemDetails{
					Status: http.StatusInternalServerError,
					Detail: "Error getting saved prediction information: " + errLoadPrediction.Error(),
				}
				logger.AniLog.Error(problemDetails.Detail)
				return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
			}

			// Calculate de real values of cpu and memory average
			realCpuAverage := cpuLimitValue.Value * nfLoadPred.CpuAverage
			realMemAverage := memLimitValue.Value * nfLoadPred.MemAverage
			realThroughput := 8 * nfLoadPred.Throughput

			nfLoadValues := models.ResourcesNfLoad{
				CpuLoad: nfLoadPred.CpuAverage,
				MemLoad: nfLoadPred.MemAverage,
			}

			NfLoad = models.NwdafAnalyticsInfoNfLoad{
				NfInstanceId: profile.NfInstanceId,
				Accuracy:     selectedModels[0].Accuracy,
				NfType:       profile.NfType,
				Pod:          podName,
				Container:    containerName,
				CpuUsage:     realCpuAverage,
				MemUsage:     realMemAverage,
				Throughput:   realThroughput,
				CpuLimit:     cpuLimitValue.Value,
				MemLimit:     memLimitValue.Value,
				NfLoad:       nfLoadValues,
				NfStatus:     profile.NfStatus,
				Confidence:   selectedModels[0].Confidence,
			}

			NfLoadsAnalitics = append(NfLoadsAnalitics, NfLoad)
		}

		responseNfLoad = models.NwdafAnalyticsInfoNfLoadResponse{
			EventId:         *eventID,
			AnalysisType:    analysisType,
			TargetPeriod:    targetPeriod,
			OffSet:          offSet,
			AnaliticsNfLoad: NfLoadsAnalitics,
		}

		if len(errList) > 0 && len(responseNfLoad.AnaliticsNfLoad) <= 0 {
			problemDetails := &models.ProblemDetails{
				Status: http.StatusNotFound,
				Detail: fmt.Sprintf("Not found data for the Analitics: %v", errList),
			}
			logger.AniLog.Error(problemDetails.Detail)
			return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
		}

		// Return results
		return httpwrapper.NewResponse(http.StatusOK, nil, responseNfLoad)

	// Statistics metrics
	case startTime.Before(currentTime) && endTime.Before(currentTime):
		logger.AniLog.Info("Statistics metrics: EndTime is less than now")
		analysisType = models.AnalysisType_STATISTICS

		// Running Pods
		runningPods, errPods := pcm.GetRunningPods(instancek8s, namespace, "", endTime, pcmUri)
		if errPods != nil {
			problemDetails := models.ProblemDetails{
				Status: http.StatusBadRequest,
				Detail: "Error getting running pods from Packet Capture module",
			}
			logger.AniLog.Errorf(problemDetails.Detail)
			return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
		}

		// For each profile: get data from Prometheus
		NfLoadsAnalitics := []models.NwdafAnalyticsInfoNfLoad{}
		for _, profile := range *nfInstances {
			var NfLoad = models.NwdafAnalyticsInfoNfLoad(defaultValues)
			var podName string
			// containerName := util.GetPodNameFromIpv4(profile.Ipv4Addresses[0])[0]
			containerName := profile.ContainerName

			foundPod := pcm_models.FindPodByContainer(runningPods, containerName)

			if foundPod != nil {
				podName = foundPod.Pod
			} else {
				err := fmt.Errorf(" No pod found for the specified container: %s", containerName)
				logger.AniLog.Infof(err.Error())
				errList = append(errList, err)

				continue
			}

			cpuUsageAverage, errCpu := pcm.GetCpuUsageAverage(namespace, podName, containerName, targetPeriod, 0, endTime, pcmUri)
			memUsageAverage, errMem := pcm.GetMemUsageAverage(namespace, podName, containerName, targetPeriod, 0, endTime, pcmUri)
			cpuLimit, errLimCpu := pcm.GetResourceLimit(namespace, podName, containerName, pcm_models.PrometheusUnit_CORE, endTime, pcmUri)
			memLimit, errLimMem := pcm.GetResourceLimit(namespace, podName, containerName, pcm_models.PrometheusUnit_BYTE, endTime, pcmUri)
			totalThroughput, errtotalThrougput := pcm.GetThroughputAverage(namespace, podName, targetPeriod, pcm_models.MetricType_TOTAL_THROUGPUT_AVERAGE, endTime, pcmUri)
			// upThroughput, errUpThroug := pcm.GetThroughputAverage(namespace, podName, targetPeriod, pcm_models.MetricType_UPLOAD_THROUGPUT_AVERAGE, endTime, pcmUri)
			// downThroughput, errDownThroug := pcm.GetThroughputAverage(namespace, podName, targetPeriod, pcm_models.MetricType_DOWNLOAD_THROUGPUT_AVERAGE, endTime, pcmUri)

			if errCpu != nil || errMem != nil || errLimCpu != nil || errLimMem != nil || errtotalThrougput != nil {
				problemDetails := models.ProblemDetails{
					Status: http.StatusInternalServerError,
					Detail: fmt.Sprintf("Error getting dsta from Packet capture module: %s, %s, %s, %s", errCpu, errMem, errLimCpu, errLimMem),
				}
				logger.AniLog.Errorf(problemDetails.Detail)
				return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
			}

			// Get the first value
			cpuAverageValue := cpuUsageAverage[0]
			memAverageValue := memUsageAverage[0]
			cpuLimitValue := cpuLimit[0]
			memLimitValue := memLimit[0]
			totalThroughputValue := totalThroughput[0]
			// upThroughputValue := upThroughput[0]
			// downThroughputValue := downThroughput[0]

			// Total Throughput
			// totalThroughput, errSumThrougput := nwdaf_util.SumThrougput(upThroughputValue.Value, downThroughputValue.Value)
			// if errSumThrougput != nil {
			// 	problemDetails := &models.ProblemDetails{
			// 		Status: http.StatusInternalServerError,
			// 		Detail: errSumThrougput.Error(),
			// 	}

			// 	logger.AniLog.Error(errSumThrougput.Error())
			// 	return httpwrapper.NewResponse(http.StatusBadRequest, nil, problemDetails)
			// }

			var nfLoad = models.ResourcesNfLoad{
				CpuLoad: nwdaf_util.GetPercentil(cpuAverageValue.Value, cpuLimitValue.Value),
				MemLoad: nwdaf_util.GetPercentil(memAverageValue.Value, memLimitValue.Value),
			}

			NfLoad = models.NwdafAnalyticsInfoNfLoad{
				NfInstanceId: profile.NfInstanceId,
				Pod:          podName,
				Container:    containerName,
				NfType:       profile.NfType,
				CpuUsage:     cpuAverageValue.Value,
				MemUsage:     memAverageValue.Value,
				Throughput:   totalThroughputValue.Value * 8,
				CpuLimit:     cpuLimitValue.Value,
				MemLimit:     memLimitValue.Value,
				NfLoad:       nfLoad,
				NfStatus:     profile.NfStatus,
			}

			NfLoadsAnalitics = append(NfLoadsAnalitics, NfLoad)
		}

		responseNfLoad = models.NwdafAnalyticsInfoNfLoadResponse{
			EventId:         *eventID,
			AnalysisType:    analysisType,
			TargetPeriod:    targetPeriod,
			OffSet:          offSet,
			AnaliticsNfLoad: NfLoadsAnalitics,
		}

		if len(errList) > 0 && len(responseNfLoad.AnaliticsNfLoad) <= 0 {
			problemDetails := &models.ProblemDetails{
				Status: http.StatusNotFound,
				Detail: fmt.Sprintf("Not found data for the Analitics: %v", errList),
			}
			logger.AniLog.Error(problemDetails.Detail)
			return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
		}

		// Return results
		return httpwrapper.NewResponse(http.StatusOK, nil, responseNfLoad)

	default:
		problemDetails := &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Detail: "EndTime must be greater than StartTime",
		}
		logger.AniLog.Error("Invalid time range: ", problemDetails.Detail)
		return httpwrapper.NewResponse(http.StatusBadRequest, nil, problemDetails)
	}
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

// Get ML by NfType, Size, Accuracy
func getMlModelByProfile(mlmodels *[]models.MlModelData, nftype *models.NfType, reqAccuracy *models.NwdafMlModelAccuracy) (mlmodel []models.MlModelData) {

	if len(*mlmodels) <= 0 || mlmodels == nil {
		logger.AniLog.Error("No Found MlModels")
		return nil
	}

	// Filter By NfType
	var nfModels []models.MlModelData
	for _, model := range *mlmodels {
		if model.NfType == *nftype {
			nfModels = append(nfModels, model)
		}
	}

	if len(nfModels) <= 0 || nfModels == nil {
		return nil
	}

	// Filter By Accuracy
	var accuracyModels []models.MlModelData
	if reqAccuracy != nil && *reqAccuracy != "" { // Request Accuracy
		for _, model := range nfModels {
			if model.Accuracy == *reqAccuracy {
				accuracyModels = append(accuracyModels, model)
			}
		}
	} else { // Default Accuracy Priority
		definedAccuracies := models.NewNwdafMlModelAccuracyPriority()
		for _, priority := range definedAccuracies {
			for _, model := range nfModels {
				if model.Accuracy == priority {
					accuracyModels = append(accuracyModels, model)
				}
			}
		}
	}

	if len(accuracyModels) <= 0 || accuracyModels == nil {
		return nil
	}

	// Select the smallest models
	minSize := accuracyModels[0].Size
	for _, model := range accuracyModels {
		if model.Size < minSize {
			minSize = model.Size
		}
	}
	var smallestModels []models.MlModelData
	for _, model := range accuracyModels {
		if model.Size == minSize {
			smallestModels = append(smallestModels, model)
		}
	}

	if len(smallestModels) == 0 || smallestModels == nil {
		return nil
	}

	return smallestModels
}

func loadPredictionInfoFromJson(nfLoadPred *models.PredictionResult, filePath string) (err error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("could not read the saved prediction Info")
	}

	// Verify if the file is empty
	if len(data) == 0 {
		return fmt.Errorf("the prediction info is empty")
	}

	// Try parse the JSON file
	err = json.Unmarshal(data, &nfLoadPred)
	if err != nil {
		return fmt.Errorf("failed to parse the prediction info")
	}

	// Verificar si el contenido tiene datos válidos
	if math.IsNaN(nfLoadPred.CpuAverage) || math.IsNaN(nfLoadPred.MemAverage) || math.IsNaN(nfLoadPred.Throughput) {
		return fmt.Errorf("prediction info is missing required fields")
	}

	return nil
}
