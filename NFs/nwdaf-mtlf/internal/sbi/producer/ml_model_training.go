package producer

import (
	"encoding/base64"
	"encoding/json"
	"errors"
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
	pcm_models "github.com/enable-intelligent-containerized-5g/openapi/PacketCaptureModule/models"
	"github.com/enable-intelligent-containerized-5g/openapi/models"
	nwdaf_util "github.com/enable-intelligent-containerized-5g/openapi/nwdaf/util"
	"github.com/free5gc/nwdaf/internal/logger"

	"github.com/free5gc/nwdaf/internal/sbi/consumer"
	"github.com/free5gc/nwdaf/internal/util"
	"github.com/free5gc/nwdaf/pkg/factory"
	"github.com/free5gc/util/httpwrapper"
)

func HandleMlModelTrainingNfLoadMetric(request *httpwrapper.Request) (response *httpwrapper.Response) {
	logger.MlModelTrainingLog.Info("Handle MlModelTrainingNfLoadMetricRequest")

	nwdafMlTrainingReq, ok := request.Body.(models.NwdafMlModelTrainingRequest)
	if !ok {
		problemDetails := models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Detail: "The request body is't type NwdafMlModelTrainingRequest",
		}
		logger.MlModelTrainingLog.Errorf(problemDetails.Detail)
		return httpwrapper.NewResponse(http.StatusForbidden, nil, problemDetails)
	}

	putData, created, problemDetails := MlModelTrainingNfLoadProcedure(nwdafMlTrainingReq)
	if created {
		return httpwrapper.NewResponse(http.StatusCreated, nil, putData)
	} else if problemDetails != nil {
		logger.MlModelTrainingLog.Errorf("Error training the Ml model: %s", problemDetails.Detail)
		return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}

	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Detail: "UNSPECIFIED",
	}

	logger.MlModelTrainingLog.Error("Ml model Training failed")
	return httpwrapper.NewResponse(http.StatusForbidden, nil, problemDetails)
}

func MlModelTrainingNfLoadProcedure(mlTrainingReq models.NwdafMlModelTrainingRequest) (models.MlModelTrainingResponse, bool, *models.ProblemDetails) {
	logger.MlModelTrainingLog.Info("Procedure MlModelTrainingProcedure")

	// Request data
	eventID := mlTrainingReq.EventId
	targetPeriod := mlTrainingReq.TargetPeriod
	nfType := mlTrainingReq.NfType
	datasetFileReq := mlTrainingReq.File

	// Variables
	currentTime := time.Now()
	dataPath := util.NwdafDefaultDataPath
	dataLabeledPath := util.NwdafDefaultDataLabeledPath
	dataRawPath := util.NwdafDefaultDataRawPath
	var datasetFile string
	var selectedDatasetFile string
	var baseName string
	var nameId string

	// Check the TargetPeriod
	if targetPeriod < 60 {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusBadRequest,
			Detail: "The difference between the start date and the end date must be greater than 60 seconds",
		}
		return models.MlModelTrainingResponse{}, false, problemDetails
	}


	var statusGettingData int32
	var errGettingData error
	if strings.TrimSpace(datasetFileReq.Data) == "" || strings.TrimSpace(datasetFileReq.Name) == "" {
		logger.MlModelTrainingLog.Info("There is not a csv file in the request")
		// Get Data from PCM
		statusGettingData, errGettingData = GetDataForNfLoadFromPcm(mlTrainingReq, currentTime)
	}else{
		logger.MlModelTrainingLog.Info("There is a csv file in the request")
		statusGettingData, errGettingData = GetDataForNfLoadFromUploadedFile(mlTrainingReq)
	}
	if errGettingData != nil {
		problemDetails := &models.ProblemDetails{
			Status: statusGettingData,
			Detail: errGettingData.Error(),
		}
		logger.MlModelTrainingLog.Error(problemDetails.Detail)
		return models.MlModelTrainingResponse{}, false, problemDetails
	}
	logger.MlModelTrainingLog.Infof("Getting data completed and saved in: %s", dataRawPath)


	// Processing Data
	statusProcessingData, errProcessingData := ProcessingDataForNfLoad(&datasetFile, &selectedDatasetFile, &baseName, &nameId, mlTrainingReq, currentTime)
	if errProcessingData != nil {
		problemDetails := &models.ProblemDetails{
			Status: statusProcessingData,
			Detail: errProcessingData.Error(),
		}
		logger.MlModelTrainingLog.Error(problemDetails.Detail)
		return models.MlModelTrainingResponse{}, false, problemDetails
	}
	logger.MlModelTrainingLog.Infof("Data processing completed and saved in: %s", dataLabeledPath+datasetFile)


	// Training Model
	statusTrainingModel, errTrainingModel := TrainingModelForNfLoad(baseName, nameId, datasetFile)
	if errTrainingModel != nil {
		problemDetails := &models.ProblemDetails{
			Status: statusTrainingModel,
			Detail: errTrainingModel.Error(),
		}
		logger.MlModelTrainingLog.Error(problemDetails.Detail)
		return models.MlModelTrainingResponse{}, false, problemDetails
	}
	logger.MlModelTrainingLog.Infoln("Ml Model Training completed")

	// Saving the model
	var mlModelCreated models.MlModelTrainingModelInfo
	errLoadModel := loadMlmodelInfoFromJson(&mlModelCreated, dataPath+util.NwdafDefaultModelInfoFile)
	if errLoadModel != nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Detail: "Error getting saved model information: " + errLoadModel.Error(),
		}
		logger.MlModelTrainingLog.Error(problemDetails.Detail)
		return models.MlModelTrainingResponse{}, false, problemDetails
	}

	// Get the figure
	imageBytes, errGettingFigure := os.ReadFile(mlModelCreated.FigureURI)
	if errGettingFigure != nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Detail: "Error getting the saved figure: " + errGettingFigure.Error(),
		}
		logger.MlModelTrainingLog.Error(problemDetails.Detail)
		return models.MlModelTrainingResponse{}, false, problemDetails
	}
	// Encode the figure
	figureSavedBase64 := base64.StdEncoding.EncodeToString(imageBytes)

	modelConfidence := models.MlModelDataConfidence{
		R2:           mlModelCreated.R2,
		MSE:          mlModelCreated.MSE,
		R2Cpu:        mlModelCreated.R2CPU,
		R2Mem:        mlModelCreated.R2Mem,
		R2Troughput:  mlModelCreated.R2Thrpt,
		MSECpu:       mlModelCreated.MSECPU,
		MSEMem:       mlModelCreated.MSEMem,
		MSETroughput: mlModelCreated.MSEThrpt,
	}

	mlModelInfo := models.MlModelData{
		EventId:      eventID,
		Name:         mlModelCreated.Name,
		Size:         mlModelCreated.Size,
		FigureURI:    mlModelCreated.FigureURI,
		TargetPeriod: targetPeriod,
		Confidence:   modelConfidence,
		URI:          mlModelCreated.URI,
		Accuracy:     models.SetAcuracy(modelConfidence.R2),
		NfType:       nfType,
	}

	mlModelSaveResponse, saved, errSave := util.SaveMlModelProcedure(mlModelInfo)
	if !saved || len(mlModelSaveResponse.MlModels) <= 0 {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Detail: "Error saving the Ml Model in  the DB: " + errSave.Detail,
		}
		logger.MlModelTrainingLog.Error(problemDetails.Detail)
		return models.MlModelTrainingResponse{}, false, problemDetails
	}

	mlModelSaved := mlModelSaveResponse.MlModels[0]

	var modelInfoResponse models.MlModelTrainingResponse = models.MlModelTrainingResponse{
		EventId:      mlModelSaved.EventId,
		Name:         mlModelSaved.Name,
		Size:         mlModelSaved.Size,
		FigureURI:    mlModelSaved.FigureURI,
		TargetPeriod: mlModelSaved.TargetPeriod,
		Confidence:   mlModelSaved.Confidence,
		URI:          mlModelSaved.URI,
		Accuracy:     mlModelSaved.Accuracy,
		NfType:       mlModelSaved.NfType,
		Figure:       figureSavedBase64,
	}

	return modelInfoResponse, true, nil
}

func selecDataset(dirPath string, start int64, baseName string) (newID nwdaf_util.PairNum, err error) {
	filesCsv, errLoadFiles := nwdaf_util.LoadCsvFiles(dirPath)
	var listNum []nwdaf_util.PairNum

	if errLoadFiles == nil {
		for _, file := range filesCsv {
			fileName := strings.TrimSuffix(file, ".csv")
			parts := strings.Split(fileName, "_")

			if len(parts) == 7 {
				baseNameFile := strings.Join(parts[:5], "_")
				lastValue1 := parts[len(parts)-2]
				lastValue2 := parts[len(parts)-1]

				num1, errNum1 := strconv.ParseInt(lastValue1, 10, 64)
				num2, errNum2 := strconv.ParseInt(lastValue2, 10, 64)

				// Have Error
				if errNum1 != nil && errNum2 != nil {
					return newID, fmt.Errorf("%d, %d not are digits", num1, num2)
				}

				// logger.MlModelTrainingLog.Warn(baseNameFile, "  ",baseName)
				if baseNameFile == baseName {
					listNum = append(listNum, nwdaf_util.PairNum{Start: num1, End: num2})
				}
			}
		}

		if len(listNum) > 0 {
			var minNum = nwdaf_util.PairNum{Start: math.MaxInt64, End: math.MaxInt64}
			var maxNum = nwdaf_util.PairNum{Start: math.MinInt64, End: math.MinInt64}
			var selectedDatasets []nwdaf_util.PairNum
			var filteredByMin []nwdaf_util.PairNum

			// Filter  datasets
			for _, num := range listNum {
				if num.Start < start && start < num.End {
					selectedDatasets = append(selectedDatasets, num)
				}
			}
			// logger.MlModelTrainingLog.Warn("filtered datasets: ", selectedDatasets)

			// Select de min StartTime
			for _, num := range selectedDatasets {
				if num.Start < minNum.Start {
					minNum = num
				}
			}

			// Filter by min StartTime
			for _, num := range selectedDatasets {
				if num.Start == minNum.Start {
					filteredByMin = append(filteredByMin, num)
				}
			}

			// Select the max EndTime
			if len(filteredByMin) > 0 {
				for _, num := range filteredByMin {
					if num.End > maxNum.End {
						maxNum = num
					}
				}
			}
			// logger.MlModelTrainingLog.Warn("filter max: ", maxNum)

			// Select the newID
			if maxNum.Start < start && start < maxNum.End {
				return maxNum, nil
			}
		}

	}

	return nwdaf_util.PairNum{}, fmt.Errorf("no found a dataset for: %s", baseName)
}

func loadMlmodelInfoFromJson(modelInfo *models.MlModelTrainingModelInfo, filePath string) (err error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("could not read the saved model Info")
	}

	// Verify if the file is empty
	if len(data) == 0 {
		return fmt.Errorf("the model info is empty")
	}

	// Try parse the JSON file
	err = json.Unmarshal(data, &modelInfo)
	if err != nil {
		return fmt.Errorf("failed to parse the model info")
	}

	// Verificar si el contenido tiene datos válidos
	if modelInfo.Size <= 0 || modelInfo.URI == "" || math.IsNaN(modelInfo.MSE) ||
		math.IsNaN(modelInfo.R2) || math.IsNaN(modelInfo.MSECPU) ||
		math.IsNaN(modelInfo.MSEMem) || math.IsNaN(modelInfo.MSEThrpt) ||
		math.IsNaN(modelInfo.R2CPU) || math.IsNaN(modelInfo.R2Mem) ||
		math.IsNaN(modelInfo.R2Thrpt) {
		return fmt.Errorf("model info is missing required fields")
	}

	return nil
}

func HandleSaveMlModel(request *httpwrapper.Request) *httpwrapper.Response {
	logger.MlModelTrainingLog.Info("Handle SaveMlModel")

	mlmodeldata, ok := request.Body.(models.MlModelData)
	if !ok {
		return httpwrapper.NewResponse(http.StatusForbidden, nil, "The request body is't type MlModelData")
	}

	putData, created, problemDetails := util.SaveMlModelProcedure(mlmodeldata)
	if created {
		// logger.MlModelTrainingLog.Info("SaveMlModel success")
		return httpwrapper.NewResponse(http.StatusCreated, nil, putData)
	} else if problemDetails != nil {
		// logger.MlModelTrainingLog.Errorf("SaveMlModel failed: %s", problemDetails.Detail)
		return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}

	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Detail: "UNSPECIFIED",
	}
	logger.MlModelTrainingLog.Error("SaveMlModel failed")
	return httpwrapper.NewResponse(http.StatusForbidden, nil, problemDetails)
}

func GetDataForNfLoadFromPcm(reqMlData models.NwdafMlModelTrainingRequest, currentTime time.Time) (int32, error) {
	logger.MlModelTrainingLog.Info("Getting data from Packet Capture Module")
	// Variables
	namespace := factory.NwdafConfig.Configuration.Namespace
	instancek8s := factory.NwdafConfig.Configuration.KsmInstance
	pcmUri := factory.NwdafConfig.Configuration.OamUri
	eventID := reqMlData.EventId
	targetPeriod := reqMlData.TargetPeriod
	nfType := reqMlData.NfType
	startTime := reqMlData.StartTime.UTC()

	// File Paths
	dataRawPath := util.NwdafDefaultDataRawPath
	menUsageFile := util.NwdafDefaultMenUsageFile
	cpuUsageFile := util.NwdafDefaultCpuUsageFile
	totalThroughputFile := util.NwdafDefaultTotalThroughputFile

	logger.MlModelTrainingLog.Infof("Event ID: %s, tp: %d, NF: %s, StartTime: %s, %s", eventID, targetPeriod, nfType, startTime, currentTime)

	NrfUri := factory.NwdafConfig.Configuration.NrfUri
	if NrfUri == "" {
		return http.StatusInternalServerError, errors.New("NrfUri is not set")
	}

	// Running Pods
	runningPods, errPods := pcm.GetRunningPods(instancek8s, namespace, "", currentTime, pcmUri)
	if errPods != nil {
		return http.StatusInternalServerError, fmt.Errorf("error getting running pods from Packet Capture module: %s", errPods.Error())
	}

	param := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		// ServiceNames: optional.Interface{},
	}
	// Var to store all NF instances
	var nfInstances []models.NfProfile
	// Search all NF instances
	err := consumer.SearchAllNfInstance(&nfInstances, NrfUri, nfType, models.NfType_NWDAF, param)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error getting %s NfInstances: %s", nfType, err.Error())
	}

	if len(nfInstances) <= 0 {
		return http.StatusNotFound, fmt.Errorf("nf type %s not found", nfType)
	}

	// Select the firts profile
	profile := nfInstances[0]
	var podName string
	containerName := profile.ContainerName

	// Getting data from Prometheus
	logger.MlModelTrainingLog.Info("Getting data from Prometheus")
	foundPod := pcm_models.FindPodByContainer(runningPods, containerName)
	if foundPod != nil {
		podName = foundPod.Pod
	} else {
		return http.StatusNotFound, fmt.Errorf("no pod found for the specified container: %s", containerName)
	}

	// Get CPU and RAM  from Ml Model Training
	cpuUsageAverageRange, errCpu := pcm.GetCpuUsageAverageRange(namespace, podName, containerName, targetPeriod, 0, startTime, currentTime, pcmUri)
	memUsageAverageRange, errMem := pcm.GetMemUsageAverageRange(namespace, podName, containerName, targetPeriod, 0, startTime, currentTime, pcmUri)
	totalThroughputRange, errtotalThrougput := pcm.GetThroughputAverageRange(namespace, podName, targetPeriod, 0, pcm_models.MetricType_TOTAL_THROUGPUT_AVERAGE, startTime, currentTime, pcmUri)
	cpuLimit, errLimCpu := pcm.GetResourceLimit(namespace, podName, containerName, pcm_models.PrometheusUnit_CORE, currentTime, pcmUri)
	memLimit, errLimMem := pcm.GetResourceLimit(namespace, podName, containerName, pcm_models.PrometheusUnit_BYTE, currentTime, pcmUri)

	if errCpu != nil || errMem != nil || errLimCpu != nil || errLimMem != nil || errtotalThrougput != nil {
		return http.StatusInternalServerError, fmt.Errorf("error getting data from Packet capture module: %s, %s, %s, %s", errCpu, errMem, errLimCpu, errLimMem)
	}

	logger.MlModelTrainingLog.Info("Saving data")
	cpuLimitValue := cpuLimit[0]
	memLimitValue := memLimit[0]
	nwdaf_util.DivideValues(&cpuUsageAverageRange, cpuLimitValue.Value)
	nwdaf_util.DivideValues(&memUsageAverageRange, memLimitValue.Value)
	pcm_models.UpdateContainerNameInPrometheusResultList(&totalThroughputRange, containerName)

	// Llamar a la función para escribir el JSON
	pathCpuUsage := dataRawPath + cpuUsageFile
	errToCsvCpu := nwdaf_util.SaveToJson(pathCpuUsage, cpuUsageAverageRange)
	if errToCsvCpu != nil {
		return http.StatusInternalServerError, fmt.Errorf("error: %s", errToCsvCpu)
	} else {
		logger.MlModelTrainingLog.Infof("CpuUsage saved in %s (%d rows)", pathCpuUsage, len(cpuUsageAverageRange))
	}

	// Llamar a la función para escribir el JSON
	pathMemUsage := dataRawPath + menUsageFile
	errToCsvMem := nwdaf_util.SaveToJson(pathMemUsage, memUsageAverageRange)
	if errToCsvMem != nil {
		return http.StatusInternalServerError, fmt.Errorf("error: %s", errToCsvMem)
	} else {
		logger.MlModelTrainingLog.Infof("MemUsage saved in %s (%d rows)", pathMemUsage, len(memUsageAverageRange))
	}

	// Llamar a la función para escribir el JSON
	pathTotalThroughput := dataRawPath + totalThroughputFile
	errToCsvThroughput := nwdaf_util.SaveToJson(pathTotalThroughput, totalThroughputRange)
	if errToCsvThroughput != nil {
		return http.StatusInternalServerError, fmt.Errorf("error: %s", errToCsvThroughput)
	} else {
		logger.MlModelTrainingLog.Infof("CpuUsage saved in %s (%d rows)", pathTotalThroughput, len(totalThroughputRange))
	}

	return 0, nil
}

func GetDataForNfLoadFromUploadedFile(reqMlData models.NwdafMlModelTrainingRequest) (int32, error) {
	logger.MlModelTrainingLog.Info("Getting data from uploaded file")

	return 0, nil
}

func ProcessingDataForNfLoad(datasetFile *string, selectedDatasetFile *string, baseName *string, nameId *string, reqMlData models.NwdafMlModelTrainingRequest, currentTime time.Time) (int32, error) {
	logger.MlModelTrainingLog.Info("Processing data")
	// Variables
	newDataset := reqMlData.NewDataset
	startTimeSeconds := reqMlData.StartTime.Unix()
	currentTimeSeconds := currentTime.Unix()
	cpuColumn := string(pcm_models.MetricType_CPU_USAGE_AVERAGE)
	memColumn := string(pcm_models.MetricType_MEMORY_USAGE_AVERAGE)
	thrptColumn := string(pcm_models.MetricType_TOTAL_THROUGPUT_AVERAGE)
	eventID := reqMlData.EventId
	targetPeriod := reqMlData.TargetPeriod
	nfType := reqMlData.NfType

	// File Paths
	dataPath := util.NwdafDefaultDataPath
	dataRawPath := util.NwdafDefaultDataRawPath
	menUsageFile := util.NwdafDefaultMenUsageFile
	cpuUsageFile := util.NwdafDefaultCpuUsageFile
	totalThroughputFile := util.NwdafDefaultTotalThroughputFile
	pathDataProcessingScript := util.NwdafDefaultDataProcessingScriptPath
	dataPreprocessedPath := util.NwdafDefaultDataPreprocessedPath
	dataProcessedPath := util.NwdafDefaultDataProcessedPath
	dataLabeledPath := util.NwdafDefaultDataLabeledPath
	// Build the datasetName
	*baseName = fmt.Sprintf("%s_%s_%ds", eventID, nfType, targetPeriod)
	*nameId = fmt.Sprintf("%d_%d", startTimeSeconds, currentTimeSeconds)
	baseNameDataset := fmt.Sprintf("dataset_%s", *baseName)
	*datasetFile = fmt.Sprintf("%s_%s.csv", baseNameDataset, *nameId)

	// Select a suitable dataset
	*selectedDatasetFile = *datasetFile
	if !newDataset {
		idSeconds, errSelectDataset := selecDataset(dataLabeledPath, startTimeSeconds, baseNameDataset)
		if errSelectDataset != nil {
			logger.MlModelTrainingLog.Warnf("No suitable dataset found for %s", *datasetFile)
		} else {
			// Define the selected dataset
			*selectedDatasetFile = fmt.Sprintf("%s_%d_%d.csv", baseNameDataset, idSeconds.Start, idSeconds.End)
			logger.MlModelTrainingLog.Warnf("Selected Dataset for (%s): %s", *datasetFile, *selectedDatasetFile)
			// Set de dataset name for the data
			*nameId = fmt.Sprintf("%d_%d", idSeconds.Start, currentTimeSeconds)
			*datasetFile = fmt.Sprintf("%s_%s.csv", baseNameDataset, *nameId)
		}
	}

	// Run processing data script
	cmd := exec.Command("python3", pathDataProcessingScript, dataPath,
		dataRawPath, dataPreprocessedPath,
		dataProcessedPath, dataLabeledPath,
		cpuUsageFile, menUsageFile, totalThroughputFile, *datasetFile, *selectedDatasetFile,
		cpuColumn, memColumn, thrptColumn)

	// Get the output and error
	outputProcess, errProcess := cmd.CombinedOutput()
	if errProcess != nil {
		return http.StatusInternalServerError, fmt.Errorf("error processing data to Ml Model Training. %s", string(outputProcess))
	}

	return 0, nil
}

func TrainingModelForNfLoad(baseName string, nameId string, datasetFile string) (int32, error) {
	logger.MlModelTrainingLog.Info("Training Ml Model")

	// Variables
	timeSteps := factory.NwdafConfig.Configuration.MlModelTrainingInfo.TimeSteps
	fullBaseName := fmt.Sprintf("%s_%s", baseName, nameId)
	modelTrainingScriptPath := util.NwdafDefaultModelTrainingScriptPath
	modelsPath := util.NwdafDefaultModelsPath
	figuresPath := util.NwdafDefaultFiguresPath
	modelInfo := util.NwdafDefaultModelInfoFile
	modelInfoList := util.NwdafDefaultModelInfoListFile
	cpuColumn := string(pcm_models.MetricType_CPU_USAGE_AVERAGE)
	memColumn := string(pcm_models.MetricType_MEMORY_USAGE_AVERAGE)
	thrptColumn := string(pcm_models.MetricType_TOTAL_THROUGPUT_AVERAGE)

	// File Paths
	dataPath := util.NwdafDefaultDataPath
	dataLabeledPath := util.NwdafDefaultDataLabeledPath

	// Run Ml model training script
	cmdTraining := exec.Command("python3", modelTrainingScriptPath,
		modelsPath, dataPath, dataLabeledPath,
		figuresPath, datasetFile, modelInfo,
		modelInfoList, cpuColumn, memColumn, thrptColumn,
		fullBaseName, strconv.FormatInt(timeSteps, 10))

	// Get the output and error
	outputTraining, errTraining := cmdTraining.CombinedOutput()
	if errTraining != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in Ml Model Training. %s", string(outputTraining))
	}
	if strings.TrimSpace(string(outputTraining)) != "" {
		logger.MlModelTrainingLog.Warn(string(outputTraining))
	}

	return 0, nil
}
