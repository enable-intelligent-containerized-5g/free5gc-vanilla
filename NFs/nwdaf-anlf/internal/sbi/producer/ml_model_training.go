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
	"github.com/enable-intelligent-containerized-5g/openapi/models"
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
		return httpwrapper.NewResponse(http.StatusForbidden, nil, "The request body is't type NwdafMlModelTrainingRequest")
	}

	putData, created, problemDetails := MlModelTrainingNfLoadProcedure(nwdafMlTrainingReq)
	if created {
		// logger.MlModelTrainingLog.Info("SaveMlModel success")
		return httpwrapper.NewResponse(http.StatusCreated, nil, putData)
	} else if problemDetails != nil {
		// logger.MlModelTrainingLog.Errorf("SaveMlModel failed: %s", problemDetails.Cause)
		return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}

	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}

	logger.MlModelTrainingLog.Error("SaveMlModel failed")
	return httpwrapper.NewResponse(http.StatusForbidden, nil, problemDetails)
}

func MlModelTrainingNfLoadProcedure(mlTrainingReq models.NwdafMlModelTrainingRequest) (models.MlModelDataResponse, bool, *models.ProblemDetails) {
	logger.MlModelTrainingLog.Info("Procedure MlModelTrainingProcedure")

	currentTime := time.Now()
	namespace := factory.NwdafConfig.Configuration.Namespace
	instancek8s := factory.NwdafConfig.Configuration.KsmInstance

	eventID := mlTrainingReq.EventId
	targetPeriod := mlTrainingReq.TargetPeriod
	nfType := mlTrainingReq.NfType
	startTime := mlTrainingReq.StartTime.UTC()
	newDataset := mlTrainingReq.NewDataset
	startTimeSeconds := mlTrainingReq.StartTime.Unix()
	currentTimeSeconds := currentTime.Unix()

	// formattedStartTime := mlTrainingReq.StartTime.Format("2006-01-02_15-04-05")
	// formattedCurrentTime := currentTime.Format("2006-01-02_15-04-05.000000000")

	logger.MlModelTrainingLog.Infof("Event ID: %s, tp: %d, NF: %s, StartTime: %s, %s", eventID, targetPeriod, nfType, startTime, currentTime)

	NrfUri := factory.NwdafConfig.Configuration.NrfUri
	if NrfUri == "" {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "NrfUri is not set",
		}
		return models.MlModelDataResponse{}, false, problemDetails
	}

	// Running Pods
	runningPods := consumer.GetRunningPods(instancek8s, namespace, "", currentTime)
	// logger.MlModelTrainingLog.Warn(runningPods)

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
			Cause:  fmt.Sprintf("Error getting %s NfInstances: %s", nfType, err.Error()),
		}
		logger.MlModelTrainingLog.Error(problemDetails.Cause)
		return models.MlModelDataResponse{}, false, problemDetails
	}

	if len(nfInstances) <= 0 {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  fmt.Sprintf("No %s type Nfs found", nfType),
		}
		logger.MlModelTrainingLog.Error(problemDetails.Cause)
		return models.MlModelDataResponse{}, false, problemDetails
	}

	// Select the firts profile
	profile := nfInstances[0]
	var podName string
	containerName := profile.ContainerName

	// Getting data from Prometheus
	logger.MlModelTrainingLog.Info("Getting data from Prometheus")
	foundPod := findPodByContainer(runningPods, containerName)
	if foundPod != nil {
		podName = foundPod.Pod
	} else {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  fmt.Sprintf("No pod found for the specified container: %s", containerName),
		}
		logger.MlModelTrainingLog.Error(problemDetails.Cause)
		return models.MlModelDataResponse{}, false, problemDetails
	}

	// Get CPU and RAM  from Ml Model Training
	cpuUsageAverageRange := consumer.GetCpuUsageAverageRange(namespace, podName, containerName, targetPeriod, 0, startTime, currentTime)
	memUsageAverageRange := consumer.GetMemUsageAverageRange(namespace, podName, containerName, targetPeriod, 0, startTime, currentTime)
	cpuLimit := consumer.GetResourceLimit(namespace, podName, containerName, models.PrometheusUnit_CORE, currentTime)[0]
	memLimit := consumer.GetResourceLimit(namespace, podName, containerName, models.PrometheusUnit_BYTE, currentTime)[0]

	logger.MlModelTrainingLog.Info("Saving data")
	divideValues(&cpuUsageAverageRange, cpuLimit.Value)
	divideValues(&memUsageAverageRange, memLimit.Value)

	// Data paths
	dataPath := util.NwdafDefaultDataPath
	dataRawPath := util.NwdafDefaultDataRawPath
	menUsageFile := util.NwdafDefaultMenUsageFile
	cpuUsageFile := util.NwdafDefaultCpuUsageFile

	// Llamar a la función para escribir el JSON
	pathCpuUsage := dataRawPath + cpuUsageFile
	errToCsvCpu := saveToJson(pathCpuUsage, cpuUsageAverageRange)
	if errToCsvCpu != nil {
		logger.MlModelTrainingLog.Error("Error: ", errToCsvCpu)
	} else {
		logger.MlModelTrainingLog.Infof("CpuUsage saved in %s (%d rows)", pathCpuUsage, len(cpuUsageAverageRange))
	}

	// Llamar a la función para escribir el JSON
	pathMemUsage := dataRawPath + menUsageFile
	errToCsvMem := saveToJson(pathMemUsage, memUsageAverageRange)
	if errToCsvMem != nil {
		logger.MlModelTrainingLog.Error("Error: ", errToCsvMem)
	} else {
		logger.MlModelTrainingLog.Infof("MemUsage saved in %s (%d rows)", pathMemUsage, len(memUsageAverageRange))
	}

	// Processing data
	logger.MlModelTrainingLog.Info("Processing data")
	cpuColumn := string(models.MetricType_CPU_USAGE_AVERAGE)
	memColumn := string(models.MetricType_MEMORY_USAGE_AVERAGE)
	pathDataProcessingScript := util.NwdafDefaultDataProcessingScriptPath
	dataPreprocessedPath := util.NwdafDefaultDataPreprocessedPath
	dataProcessedPath := util.NwdafDefaultDataProcessedPath
	dataLabeledPath := util.NwdafDefaultDataLabeledPath

	// Build the datasetName	
	baseName := fmt.Sprintf("%s_%s_%ds", eventID, nfType, targetPeriod)
	nameID := fmt.Sprintf("%d_%d", startTimeSeconds, currentTimeSeconds)
	baseNameDataset := fmt.Sprintf("dataset_%s", baseName)
	datasetFile := fmt.Sprintf("%s_%s.csv", baseNameDataset, nameID)

	// Select a suitable dataset
	selectedDatasetFile := datasetFile
	if !newDataset {
		nameID, errSelectDataset := selecDataset(dataLabeledPath, startTimeSeconds, baseNameDataset)
		if errSelectDataset != nil {
			logger.MlModelTrainingLog.Warnf("No suitable dataset found for '%s'\n", datasetFile)
		} else {
			selectedDatasetFile = fmt.Sprintf("%s_%s.csv", baseNameDataset, nameID)
			datasetFile = selectedDatasetFile
			logger.MlModelTrainingLog.Warnf("Selected Dataset for %s: %s", datasetFile, selectedDatasetFile)
		}
	}

	// Run processing data script
	cmd := exec.Command("python3", pathDataProcessingScript, dataPath, 
						dataRawPath, dataPreprocessedPath, 
						dataProcessedPath, dataLabeledPath, 
						cpuUsageFile, menUsageFile, datasetFile, 
						cpuColumn, memColumn)

	// Get teh output and error
	outputProcess, errProcess := cmd.CombinedOutput()
	if errProcess != nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  fmt.Sprintf("Error processing data to Ml Model Training. %s", string(outputProcess)),
		}
		logger.MlModelTrainingLog.Error(problemDetails.Cause)
		return models.MlModelDataResponse{}, false, problemDetails
	}
	logger.MlModelTrainingLog.Infof("Data processing completed and saved in: %s", dataLabeledPath + datasetFile)

	// Training Model
	logger.MlModelTrainingLog.Info("Training Ml Model")
	fullBaseName:= fmt.Sprintf("%s_%s", baseName, nameID)
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
								fullBaseName)
	// Get the output and error
	outputTraining, errTraining := cmdTraining.CombinedOutput()
	if errTraining != nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  fmt.Sprintf("Error in Ml Model Training. %s", string(outputTraining)),
		}
		logger.MlModelTrainingLog.Error(problemDetails.Cause)
		return models.MlModelDataResponse{}, false, problemDetails
	}
	logger.MlModelTrainingLog.Infoln("Ml Model Training completed")

	// Save the model
	

	problemDetails := &models.ProblemDetails{
		Status: http.StatusOK,
		Cause:  "Working feature: " + podName,
	}

	return models.MlModelDataResponse{}, false, problemDetails
}

func selecDataset(dirPath string, start int64, baseName string) (newID string, err error) {
	type PairNum struct {
		Start int64
		End   int64
	}
	filesCsv, errLoadFiles := loadCsvFiles(dirPath)
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
			var minNum = PairNum{
				Start: math.MaxInt64,
				// End: math.MaxInt64,
			}
			var maxNum = PairNum{
				Start: math.MinInt64,
				End:   math.MinInt64,
			}
			var listSelectedNum []PairNum

			// Select de min StartTime
			for _, num := range listNum {
				if num.Start < minNum.Start {
					minNum = num
				}
			}

			// Filter by min StartTime
			for _, num := range listNum {
				if num.Start == minNum.Start {
					listSelectedNum = append(listSelectedNum, num)
				}
			}
			// logger.MlModelTrainingLog.Warn("filter mi: ", listSelectedNum)

			// Select the max EndTime
			if len(listSelectedNum) > 0 {
				for _, num := range listSelectedNum {
					if num.End > maxNum.End {
						maxNum = num
					}
				}
			}
			// logger.MlModelTrainingLog.Warn("filter max: ", maxNum)

			// Select the newID
			if maxNum.Start < start {
				return fmt.Sprintf("%d_%d", maxNum.Start, maxNum.End), nil
			}
		}

	}

	return newID, fmt.Errorf(" No foun a dataset for: %s", baseName)
}

func divideValues(results *[]models.PrometheusResult, divisor float64) {
	if math.IsNaN(divisor) || divisor == 0 {
		divisor = 1
	}
	for i := range *results {
		(*results)[i].Value /= divisor
	}
}

// Función para guardar una estructura en un archivo JSON
func saveToJson(filename string, data interface{}) error {
	// Crear el archivo
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Convertir los datos a JSON con indentación
	indentedData, err := json.MarshalIndent(data, "", "    ") // Usa 4 espacios como indentación
	if err != nil {
		return err
	}

	// Escribir los datos con indentación al archivo
	_, err = file.Write(indentedData)
	return err
}

func loadCsvFiles(dirPath string) (files []string, err error) {
	filesDir, err := os.ReadDir(dirPath)
	if err != nil {
		return files, err
	}

	for _, file := range filesDir {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".csv") {
			files = append(files, file.Name())
		}
	}

	return files, nil
}

// func loadFromFile(filename string, data interface{}) error {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	decoder := json.NewDecoder(file)
// 	return decoder.Decode(data)
// }

// func writeCSV(filename string, metrics []consumer.PrometheusResult) error {
// 	file, err := os.Create(filename)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	writer := csv.NewWriter(file)
// 	defer writer.Flush()

// 	// Escribir el encabezado en el archivo CSV
// 	header := []string{
// 		"Pod", "Container",
// 		"Timestamp1", "CpuUsage1",
// 		"Timestamp2", "CpuUsage2",
// 		"Timestamp3", "CpuUsage3",
// 		"Timestamp4", "CpuUsage4",
// 	}
// 	if err := writer.Write(header); err != nil {
// 		return err
// 	}

// 	// Iterar sobre los datos en bloques de 4
// 	for i := 0; i < len(metrics); i += 4 {
// 		if i+3 < len(metrics) {
// 			// Usar el pod y container del primer elemento del bloque como referencia
// 			pod := metrics[i].Pod
// 			container := metrics[i].Container

// 			// Crear una fila con Pod, Container y los siguientes 4 valores de CpuUsage
// 			row := []string{pod, container}

// 			// Agregar los valores de CpuUsage
// 			for j := i; j < i+4; j++ {
// 				row = append(row,
// 					strconv.FormatInt(int64(metrics[j].Timestamp), 10), // Convertir Timestamp a string
// 					fmt.Sprintf("%f", metrics[j].Value),                // Convertir el valor de CpuUsage en formato flotante
// 				)
// 			}

// 			// Escribir la fila en el archivo CSV
// 			if err := writer.Write(row); err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	return nil
// }

// func writeModelToCSV(filename string, results []consumer.PrometheusResult) error {
// 	// Crear el archivo CSV
// 	file, err := os.Create(filename)
// 	if err != nil {
// 		return fmt.Errorf("error creating file: %v", err)
// 	}
// 	defer file.Close()

// 	// Crear un escritor CSV
// 	writer := csv.NewWriter(file)
// 	defer writer.Flush()

// 	// Escribir los encabezados del CSV
// 	err = writer.Write([]string{"Timestamp", "Namespace", "Pod", "Container", "Value"})
// 	if err != nil {
// 		return fmt.Errorf("error writing header: %v", err)
// 	}

// 	// Escribir los datos de los resultados
// 	for _, result := range results {
// 		record := []string{
// 			fmt.Sprintf("%f", result.Timestamp),
// 			fmt.Sprintf("%f", result.Value),
// 			string(result.MetricType),
// 			result.Namespace,
// 			result.Pod,
// 			result.Container,
// 			result.Phase,
// 			result.Uid,
// 		}
// 		err := writer.Write(record)
// 		if err != nil {
// 			return fmt.Errorf("error writing record: %v", err)
// 		}
// 	}

// 	return nil
// }

// func prepareData(data []consumer.PrometheusResult) ([][]float64, []float64) {
// 	features := make([][]float64, len(data))
// 	labels := make([]float64, len(data))

// 	for i, d := range data {
// 		features[i] = []float64{d.CpuUsage1, d.CpuUsage2, d.CpuUsage3}
// 		labels[i] = d.CpuUsage4
// 	}
// 	return features, labels
// }
