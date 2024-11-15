package producer

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/enable-intelligent-containerized-5g/openapi/Nnrf_NFDiscovery"
	"github.com/enable-intelligent-containerized-5g/openapi/models"
	"github.com/free5gc/nwdaf/internal/logger"
	"github.com/free5gc/nwdaf/internal/sbi"
	"github.com/free5gc/nwdaf/internal/sbi/consumer"
	"github.com/free5gc/nwdaf/pkg/factory"
	"github.com/free5gc/util/httpwrapper"
)

func HandleMlModelTrainingNfLoadMetric(request *httpwrapper.Request) (response *httpwrapper.Response) {
	logger.MlModelTrainingLog.Info("Handle MlModelTrainingNfLoadMetricRequest")

	nwdafMlTrainingReq, ok := request.Body.(sbi.NwdafMlModelTrainingRequest)
	if !ok {
		return httpwrapper.NewResponse(http.StatusForbidden, nil, "The request body is't type NwdafMlModelTrainingRequest")
	}

	putData, created, problemDetails := MlModelTrainingProcedure(nwdafMlTrainingReq)
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

func MlModelTrainingProcedure(mlTrainingReq sbi.NwdafMlModelTrainingRequest) (models.MlModelDataResponse, bool, *models.ProblemDetails) {
	logger.MlModelTrainingLog.Info("Procedure MlModelTrainingProcedure")

	currentTime := time.Now()
	namespace := factory.NwdafConfig.Configuration.Namespace
	instancek8s := factory.NwdafConfig.Configuration.KsmInstance

	eventID := mlTrainingReq.EventId
	targetPeriod := mlTrainingReq.TargetPeriod
	nfType := mlTrainingReq.NfType
	startTime := mlTrainingReq.StartTime.UTC()

	logger.MlModelTrainingLog.Infof("Event ID: %s, tp: %d, NF: %s, StartTime: %s", eventID, targetPeriod, nfType, startTime)

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
		return models.MlModelDataResponse{}, false, problemDetails
	}

	if len(namespace) <= 0 {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  fmt.Sprintf("No %s type Nfs found", nfType),
		}
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
		return models.MlModelDataResponse{}, false, problemDetails
	}

	// Get CPU and RAM  from Ml Model Training
	cpuUsageAverageRange := consumer.GetCpuUsageAverageRange(namespace, podName, containerName, targetPeriod, 0, startTime, currentTime)
	memUsageAverageRange := consumer.GetMemUsageAverageRange(namespace, podName, containerName, targetPeriod, 0, startTime, currentTime)
	cpuLimit := consumer.GetResourceLimit(namespace, podName, containerName, consumer.PrometheusUnit_CORE, currentTime)[0]
	memLimit := consumer.GetResourceLimit(namespace, podName, containerName, consumer.PrometheusUnit_BYTE, currentTime)[0]

	logger.MlModelTrainingLog.Info("Saving data")
	divideValues(&cpuUsageAverageRange, cpuLimit.Value)
	divideValues(&memUsageAverageRange, memLimit.Value)

	basePath := "internal/sbi/producer/data/"
	menUsageFile := "memUsage.json"
	cpuUsageFile := "cpuUsage.json"

	// Llamar a la función para escribir el CSV
	pathCpuUsage := basePath+cpuUsageFile
	errToCsvCpu := saveToJson(pathCpuUsage, cpuUsageAverageRange)
	if errToCsvCpu != nil {
		logger.MlModelTrainingLog.Error("Error: ", errToCsvCpu)
	} else {
		logger.MlModelTrainingLog.Infof("CpuUsage saved in %s (%d rows)",pathCpuUsage, len(cpuUsageAverageRange))
	}

	// Llamar a la función para escribir el CSV
	pathMemUsage := basePath+menUsageFile
	errToCsvMem := saveToJson(pathMemUsage, memUsageAverageRange)
	if errToCsvMem != nil {
		logger.MlModelTrainingLog.Error("Error: ", errToCsvMem)
	} else {
		logger.MlModelTrainingLog.Infof("MemUsage saved in %s (%d rows)",pathMemUsage, len(memUsageAverageRange))
	}

	logger.MlModelTrainingLog.Info("Processing data")
	pathProcesingData := "internal/sbi/producer/ml_model_training_scripts/process_data.py"
	cmd := exec.Command("python3", pathProcesingData, basePath, cpuUsageFile, menUsageFile)
	// Obtener la salida y el error
	output, err := cmd.CombinedOutput()
	if err != nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  fmt.Sprintf("Error processing data to Ml Model Training (%v): %s", err, string(output)),
		}
		logger.MlModelTrainingLog.Error(problemDetails.Cause)
		return models.MlModelDataResponse{}, false, problemDetails
	}

	// Imprimir la salida del script de Python
	logger.MlModelTrainingLog.Info(string(output))
	// El modelo ha sido entrenado y evaluado
	logger.MlModelTrainingLog.Info("Data processing completed ")

	// pathTrainModel := "internal/sbi/producer/ml_model_training_scripts/train_model.py"
	// cmd := exec.Command("python3", pathTrainModel)

	// Obtener la salida y el error
	// output, err := cmd.CombinedOutput()
	// if err != nil {
	// 	logger.MlModelTrainingLog.Errorf("Error al ejecutar el script de Python: %v", err)
	// }

	// // Imprimir la salida del script de Python
	// logger.MlModelTrainingLog.Info("Salida del script de Python:")
	// logger.MlModelTrainingLog.Info(string(output))

	// // El modelo ha sido entrenado y evaluado
	// logger.MlModelTrainingLog.Info("Entrenamiento completado.")

	problemDetails := &models.ProblemDetails{
		Status: http.StatusOK,
		Cause:  "Working feature: " + podName,
	}

	return models.MlModelDataResponse{}, false, problemDetails
}

func divideValues(results *[]consumer.PrometheusResult, divisor float64) {
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
	// El segundo parámetro es el prefijo y el tercero el espacio de indentación
	indentedData, err := json.MarshalIndent(data, "", "    ") // Usa 4 espacios como indentación
	if err != nil {
		return err
	}

	// Escribir los datos con indentación al archivo
	_, err = file.Write(indentedData)
	return err
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

func writeCSV(filename string, metrics []consumer.PrometheusResult) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Escribir el encabezado en el archivo CSV
	header := []string{
		"Pod", "Container",
		"Timestamp1", "CpuUsage1",
		"Timestamp2", "CpuUsage2",
		"Timestamp3", "CpuUsage3",
		"Timestamp4", "CpuUsage4",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Iterar sobre los datos en bloques de 4
	for i := 0; i < len(metrics); i += 4 {
		if i+3 < len(metrics) {
			// Usar el pod y container del primer elemento del bloque como referencia
			pod := metrics[i].Pod
			container := metrics[i].Container

			// Crear una fila con Pod, Container y los siguientes 4 valores de CpuUsage
			row := []string{pod, container}

			// Agregar los valores de CpuUsage
			for j := i; j < i+4; j++ {
				row = append(row,
					strconv.FormatInt(int64(metrics[j].Timestamp), 10), // Convertir Timestamp a string
					fmt.Sprintf("%f", metrics[j].Value),                // Convertir el valor de CpuUsage en formato flotante
				)
			}

			// Escribir la fila en el archivo CSV
			if err := writer.Write(row); err != nil {
				return err
			}
		}
	}

	return nil
}

func writeModelToCSV(filename string, results []consumer.PrometheusResult) error {
	// Crear el archivo CSV
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// Crear un escritor CSV
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Escribir los encabezados del CSV
	err = writer.Write([]string{"Timestamp", "Namespace", "Pod", "Container", "Value"})
	if err != nil {
		return fmt.Errorf("error writing header: %v", err)
	}

	// Escribir los datos de los resultados
	for _, result := range results {
		record := []string{
			fmt.Sprintf("%f", result.Timestamp),
			fmt.Sprintf("%f", result.Value),
			string(result.MetricType),
			result.Namespace,
			result.Pod,
			result.Container,
			result.Phase,
			result.Uid,
		}
		err := writer.Write(record)
		if err != nil {
			return fmt.Errorf("error writing record: %v", err)
		}
	}

	return nil
}

// func prepareData(data []consumer.PrometheusResult) ([][]float64, []float64) {
// 	features := make([][]float64, len(data))
// 	labels := make([]float64, len(data))

// 	for i, d := range data {
// 		features[i] = []float64{d.CpuUsage1, d.CpuUsage2, d.CpuUsage3}
// 		labels[i] = d.CpuUsage4
// 	}
// 	return features, labels
// }
