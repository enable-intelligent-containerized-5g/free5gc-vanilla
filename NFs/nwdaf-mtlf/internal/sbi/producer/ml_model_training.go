package producer

import (
	"encoding/base64"
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
	packetcapturemodule "github.com/enable-intelligent-containerized-5g/openapi/PacketCaptureModule"
	"github.com/enable-intelligent-containerized-5g/openapi/models"
	pcm_models "github.com/enable-intelligent-containerized-5g/openapi/PacketCaptureModule/models"
	nwdaf_util "github.com/enable-intelligent-containerized-5g/openapi/nwdaf/util"
	"github.com/free5gc/nwdaf/internal/logger"

	"github.com/free5gc/nwdaf/internal/sbi/consumer"
	"github.com/free5gc/nwdaf/internal/util"
	"github.com/free5gc/nwdaf/pkg/factory"
	"github.com/free5gc/util/httpwrapper"
)

type PairNum struct {
	Start int64
	End   int64
}

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
		Detail:  "UNSPECIFIED",
	}

	logger.MlModelTrainingLog.Error("Ml model Training failed")
	return httpwrapper.NewResponse(http.StatusForbidden, nil, problemDetails)
}

func MlModelTrainingNfLoadProcedure(mlTrainingReq models.NwdafMlModelTrainingRequest) (models.MlModelTrainingResponse, bool, *models.ProblemDetails) {
	logger.MlModelTrainingLog.Info("Procedure MlModelTrainingProcedure")

	currentTime := time.Now()
	namespace := factory.NwdafConfig.Configuration.Namespace
	instancek8s := factory.NwdafConfig.Configuration.KsmInstance
	pcmUri := factory.NwdafConfig.Configuration.OamUri

	eventID := mlTrainingReq.EventId
	targetPeriod := mlTrainingReq.TargetPeriod
	nfType := mlTrainingReq.NfType
	startTime := mlTrainingReq.StartTime.UTC()
	newDataset := mlTrainingReq.NewDataset
	startTimeSeconds := mlTrainingReq.StartTime.Unix()
	currentTimeSeconds := currentTime.Unix()

	// Check the TargetPeriod
	if targetPeriod < 60 {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusBadRequest,
			Detail:  "The difference between the start date and the end date must be greater than 60 seconds",
		}
		return models.MlModelTrainingResponse{}, false, problemDetails
	}

	// formattedStartTime := mlTrainingReq.StartTime.Format("2006-01-02_15-04-05")
	// formattedCurrentTime := currentTime.Format("2006-01-02_15-04-05.000000000")

	logger.MlModelTrainingLog.Infof("Event ID: %s, tp: %d, NF: %s, StartTime: %s, %s", eventID, targetPeriod, nfType, startTime, currentTime)

	NrfUri := factory.NwdafConfig.Configuration.NrfUri
	if NrfUri == "" {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Detail:  "NrfUri is not set",
		}
		return models.MlModelTrainingResponse{}, false, problemDetails
	}

	// Running Pods
	runningPods, errPods := packetcapturemodule.GetRunningPods(instancek8s, namespace, "", currentTime, pcmUri)
	if errPods != nil {
		problemDetails := models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Detail: fmt.Sprintf("Error getting running pods from Packet Capture module: %s", errPods.Error()),
		}
		logger.MlModelTrainingLog.Errorf(problemDetails.Detail)
		return models.MlModelTrainingResponse{}, false, &problemDetails
	}

	param := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		// ServiceNames: optional.Interface{},
	}
	// Var to store all NF instances
	var nfInstances []models.NfProfile
	// Search all NF instances
	err := consumer.SearchAllNfInstance(&nfInstances, NrfUri, nfType, models.NfType_NWDAF, param)
	if err != nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Detail:  fmt.Sprintf("Error getting %s NfInstances: %s", nfType, err.Error()),
		}
		logger.MlModelTrainingLog.Error(problemDetails.Detail)
		return models.MlModelTrainingResponse{}, false, problemDetails
	}

	if len(nfInstances) <= 0 {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Detail:  fmt.Sprintf("Nf type %s not found", nfType),
		}
		logger.MlModelTrainingLog.Error(problemDetails.Detail)
		return models.MlModelTrainingResponse{}, false, problemDetails
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
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Detail:  fmt.Sprintf("No pod found for the specified container: %s", containerName),
		}
		logger.MlModelTrainingLog.Error(problemDetails.Detail)
		return models.MlModelTrainingResponse{}, false, problemDetails
	}

	// Get CPU and RAM  from Ml Model Training
	cpuUsageAverageRange, errCpu := packetcapturemodule.GetCpuUsageAverageRange(namespace, podName, containerName, targetPeriod, 0, startTime, currentTime, pcmUri)
	memUsageAverageRange, errMem := packetcapturemodule.GetMemUsageAverageRange(namespace, podName, containerName, targetPeriod, 0, startTime, currentTime, pcmUri)
	cpuLimit, errLimCpu := packetcapturemodule.GetResourceLimit(namespace, podName, containerName, pcm_models.PrometheusUnit_CORE, currentTime, pcmUri)
	memLimit, errLimMem := packetcapturemodule.GetResourceLimit(namespace, podName, containerName, pcm_models.PrometheusUnit_BYTE, currentTime, pcmUri)
	cpuLimitValue := cpuLimit[0]
	memLimitValue := memLimit[0]

	if errCpu != nil || errMem != nil || errLimCpu != nil || errLimMem != nil {
		problemDetails := models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Detail: fmt.Sprintf("Error getting data from Packet capture module: %s, %s, %s, %s", errCpu, errMem, errLimCpu, errLimMem),
		}
		logger.MlModelTrainingLog.Errorf(problemDetails.Detail)
		return models.MlModelTrainingResponse{}, false, &problemDetails
	}

	logger.MlModelTrainingLog.Info("Saving data")
	nwdaf_util.DivideValues(&cpuUsageAverageRange, cpuLimitValue.Value)
	nwdaf_util.DivideValues(&memUsageAverageRange, memLimitValue.Value)

	// // Data paths
	dataPath := util.NwdafDefaultDataPath
	dataRawPath := util.NwdafDefaultDataRawPath
	menUsageFile := util.NwdafDefaultMenUsageFile
	cpuUsageFile := util.NwdafDefaultCpuUsageFile

	// Llamar a la función para escribir el JSON
	pathCpuUsage := dataRawPath + cpuUsageFile
	errToCsvCpu := nwdaf_util.SaveToJson(pathCpuUsage, cpuUsageAverageRange)
	if errToCsvCpu != nil {
		logger.MlModelTrainingLog.Error("Error: ", errToCsvCpu)
	} else {
		logger.MlModelTrainingLog.Infof("CpuUsage saved in %s (%d rows)", pathCpuUsage, len(cpuUsageAverageRange))
	}

	// Llamar a la función para escribir el JSON
	pathMemUsage := dataRawPath + menUsageFile
	errToCsvMem := nwdaf_util.SaveToJson(pathMemUsage, memUsageAverageRange)
	if errToCsvMem != nil {
		logger.MlModelTrainingLog.Error("Error: ", errToCsvMem)
	} else {
		logger.MlModelTrainingLog.Infof("MemUsage saved in %s (%d rows)", pathMemUsage, len(memUsageAverageRange))
	}

	// Processing data
	logger.MlModelTrainingLog.Info("Processing data")
	cpuColumn := string(pcm_models.MetricType_CPU_USAGE_AVERAGE)
	memColumn := string(pcm_models.MetricType_MEMORY_USAGE_AVERAGE)
	pathDataProcessingScript := util.NwdafDefaultDataProcessingScriptPath
	dataPreprocessedPath := util.NwdafDefaultDataPreprocessedPath
	dataProcessedPath := util.NwdafDefaultDataProcessedPath
	dataLabeledPath := util.NwdafDefaultDataLabeledPath

	// Build the datasetName
	baseName := fmt.Sprintf("%s_%s_%ds", eventID, nfType, targetPeriod)
	nameId := fmt.Sprintf("%d_%d", startTimeSeconds, currentTimeSeconds)
	baseNameDataset := fmt.Sprintf("dataset_%s", baseName)
	datasetFile := fmt.Sprintf("%s_%s.csv", baseNameDataset, nameId)
	// datasetFile = "dataset_NF_LOAD_AMF_60s_1731787200_1731825367.csv"

	// Select a suitable dataset
	selectedDatasetFile := datasetFile
	if !newDataset {
		idSeconds, errSelectDataset := selecDataset(dataLabeledPath, startTimeSeconds, baseNameDataset)
		if errSelectDataset != nil {
			logger.MlModelTrainingLog.Warnf("No suitable dataset found for '%s'\n", datasetFile)
		} else {
			// Define the selected dataset
			selectedDatasetFile = fmt.Sprintf("%s_%d_%d.csv", baseNameDataset, idSeconds.Start, idSeconds.End)
			logger.MlModelTrainingLog.Warnf("Selected Dataset for (%s): %s", datasetFile, selectedDatasetFile)
			// Set de dataset name for the data
			nameId = fmt.Sprintf("%d_%d.csv", idSeconds.Start, currentTimeSeconds)
			datasetFile = fmt.Sprintf("%s_%s.csv", baseNameDataset, nameId)
		}
	}

	// Run processing data script
	cmd := exec.Command("python3", pathDataProcessingScript, dataPath,
		dataRawPath, dataPreprocessedPath,
		dataProcessedPath, dataLabeledPath,
		cpuUsageFile, menUsageFile, datasetFile, selectedDatasetFile,
		cpuColumn, memColumn)

	// Get the output and error
	outputProcess, errProcess := cmd.CombinedOutput()
	if errProcess != nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Detail:  fmt.Sprintf("Error processing data to Ml Model Training. %s", string(outputProcess)),
		}
		logger.MlModelTrainingLog.Error(problemDetails.Detail)
		return models.MlModelTrainingResponse{}, false, problemDetails
	}
	logger.MlModelTrainingLog.Infof("Data processing completed and saved in: %s", dataLabeledPath+datasetFile)

	// Training Model
	logger.MlModelTrainingLog.Info("Training Ml Model")
	timeSteps := factory.NwdafConfig.Configuration.MlModelTrainingInfo.TimeSteps
	fullBaseName := fmt.Sprintf("%s_%s", baseName, nameId)
	// fullBaseName = "NF_LOAD_AMF_60s_1731787200_1731825367"
	modelTrainingScriptPath := util.NwdafDefaultModelTrainingScriptPath
	modelsPath := util.NwdafDefaultModelsPath
	figuresPath := util.NwdafDefaultFiguresPath
	modelInfo := util.NwdafDefaultModelInfoFile
	modelInfoList := util.NwdafDefaultModelInfoListFile
	// Run Ml model training script
	cmdTraining := exec.Command("python3", modelTrainingScriptPath,
		modelsPath, dataPath, dataLabeledPath,
		figuresPath, datasetFile, modelInfo,
		modelInfoList, cpuColumn, memColumn,
		fullBaseName, strconv.FormatInt(timeSteps, 10))
	// Get the output and error
	outputTraining, errTraining := cmdTraining.CombinedOutput()
	if errTraining != nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Detail:  fmt.Sprintf("Error in Ml Model Training. %s", string(outputTraining)),
		}
		logger.MlModelTrainingLog.Error(problemDetails.Detail)
		return models.MlModelTrainingResponse{}, false, problemDetails
	}
	if strings.TrimSpace(string(outputTraining)) != "" {
		logger.MlModelTrainingLog.Warn(string(outputTraining))
	}
	logger.MlModelTrainingLog.Infoln("Ml Model Training completed")

	// Save the model
	var mlModelCreated models.MlModelTrainingModelInfo

	errLoadModel := loadMlmodelInfoFromJson(&mlModelCreated, dataPath+util.NwdafDefaultModelInfoFile)
	if errLoadModel != nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Detail:  "Error getting saved model information: " + errLoadModel.Error(),
		}
		logger.MlModelTrainingLog.Error(problemDetails.Detail)
		return models.MlModelTrainingResponse{}, false, problemDetails
	}

	// Get the figure
	imageBytes, errGettingFigure := os.ReadFile(mlModelCreated.FigureURI)
	if errGettingFigure != nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Detail:  "Error getting the saved figure: " + errGettingFigure.Error(),
		}
		logger.MlModelTrainingLog.Error(problemDetails.Detail)
		return models.MlModelTrainingResponse{}, false, problemDetails
	}
	// Encode the figure
	figureSavedBase64 := base64.StdEncoding.EncodeToString(imageBytes)

	modelConfidence := models.MlModelDataConfidence{
		R2:     mlModelCreated.R2,
		MSE:    mlModelCreated.MSE,
		R2Cpu:  mlModelCreated.R2CPU,
		R2Mem:  mlModelCreated.R2Mem,
		MSECpu: mlModelCreated.MSECPU,
		MSEMem: mlModelCreated.MSEMem,
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
			Detail:  "Error saving the Ml Model in  the DB: " + errSave.Detail,
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

func selecDataset(dirPath string, start int64, baseName string) (newID PairNum, err error) {
	filesCsv, errLoadFiles := nwdaf_util.LoadCsvFiles(dirPath)
	var listNum []PairNum

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
					listNum = append(listNum, PairNum{Start: num1, End: num2})
				}
			}
		}

		if len(listNum) > 0 {
			var minNum = PairNum{Start: math.MaxInt64, End: math.MaxInt64}
			var maxNum = PairNum{Start: math.MinInt64, End: math.MinInt64}
			var selectedDatasets []PairNum
			var filteredByMin []PairNum

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

	return PairNum{}, fmt.Errorf("no found a dataset for: %s", baseName)
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
		math.IsNaN(modelInfo.MSEMem) || math.IsNaN(modelInfo.R2CPU) ||
		math.IsNaN(modelInfo.R2Mem) {
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
		Detail:  "UNSPECIFIED",
	}
	logger.MlModelTrainingLog.Error("SaveMlModel failed")
	return httpwrapper.NewResponse(http.StatusForbidden, nil, problemDetails)
}
