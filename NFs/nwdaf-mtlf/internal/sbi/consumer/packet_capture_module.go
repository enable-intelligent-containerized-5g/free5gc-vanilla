package consumer

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"context"

	"github.com/enable-intelligent-containerized-5g/openapi/models"
	"github.com/free5gc/nwdaf/internal/logger"
	"github.com/free5gc/nwdaf/pkg/factory"

	"github.com/prometheus/client_golang/api" // go get github.com/prometheus/client_golang/api
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

func CreateClient() (apiClient v1.API, err error) {
	// Get PcmUri
	pcmUri := factory.NwdafConfig.Configuration.OamUri

	// Create the Prometheus Client
	client, err := api.NewClient(api.Config{
		Address: pcmUri,
	})
	if err != nil {
		return apiClient, fmt.Errorf(" Error creating Prometheus client: %s", err)
	}

	apiClient = v1.NewAPI(client)

	return apiClient, err
}

func ExecutePrometheusQuery(query string, metric models.MetricType, timeReq time.Time) (metrics []models.PrometheusResult) {
	// Get PcmUri
	logger.PcmLog.Infof("Getting %s value from Prometheus", metric)
	apiClient, errClient := CreateClient()
	if errClient != nil {
		logger.PcmLog.Error(errClient)
		return metrics
	}

	// Definir el contexto y el timeout para la consulta
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Realizar una consulta para obtener el uso de CPU en tiempo real UTC
	result, warnings, err := apiClient.Query(ctx, string(query), timeReq)
	if err != nil {
		logger.PcmLog.Errorf("Error in the Request: %v", err)
	}
	if len(warnings) > 0 {
		logger.PcmLog.Warnf("Warnings: %v", warnings)
	}

	metrics = ProcessPrometheusMetricResult(result, metric)

	if metrics == nil {
		var value models.PrometheusResult
		metrics = append(metrics, value)
	}

	return metrics
}

func ExecutePrometheusQueryRange(query string, metric models.MetricType, startTime time.Time, endTime time.Time, step time.Duration) (metrics []models.PrometheusResult) {
	// Get PcmUri
	logger.PcmLog.Infof("Getting %s range from Prometheus", metric)
	apiClient, errClient := CreateClient()
	if errClient != nil {
		logger.PcmLog.Error(errClient)
		return metrics
	}

	// Definir el contexto y el timeout para la consulta
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Realizar una consulta para obtener el uso de CPU en tiempo real UTC
	result, warnings, err := apiClient.QueryRange(ctx, query, v1.Range{Start: startTime, End: endTime, Step: step * time.Second})
	if err != nil {
		logger.PcmLog.Errorf("Error in the Request: %v", err)
	}
	if len(warnings) > 0 {
		logger.PcmLog.Warnf("Warnings: %v", warnings)
	}

	metrics = ProcessPrometheusMetricResult(result, metric)

	if metrics == nil {
		var value models.PrometheusResult
		metrics = append(metrics, value)
	}

	return metrics
}

func GetPodsByPhase(instance string, ns string, phase models.KubernetesPhase, timeReq time.Time) []models.PrometheusResult {
	var params = models.PrometheusQueryParams{
		Instance:  instance,
		Namespace: ns,
		Phase:     string(phase),
	}

	query := models.BuildPodsByStatusQuery(&params)
	metric := models.MetricType_POD_STATUS

	return ExecutePrometheusQuery(query, metric, timeReq)
}

func GetCpuUsageAverage(ns string, pod string, ctnr string, tp int64, offSet int64, timeReq time.Time) []models.PrometheusResult {
	var params = models.PrometheusQueryParams{
		Namespace:    ns,
		Pod:          pod,
		Container:    ctnr,
		TargetPeriod: BuiildTargetPeriod(tp),
		// Offset:       BuiildTargetPeriod(offSet),
	}

	query := models.BuildCpuUsageAverageQuery(&params)
	metric := models.MetricType_CPU_USAGE_AVERAGE

	return ExecutePrometheusQuery(query, metric, timeReq)
}

func GetCpuUsageAverageRange(ns string, pod string, ctnr string, tp int64, offSet int64, startTime time.Time, endTime time.Time) []models.PrometheusResult {
	var params = models.PrometheusQueryParams{
		Namespace:    ns,
		Pod:          pod,
		Container:    ctnr,
		TargetPeriod: BuiildTargetPeriod(tp),
		// Offset:       BuiildTargetPeriod(offSet),
	}

	query := models.BuildCpuUsageAverageQuery(&params)
	metric := models.MetricType_CPU_USAGE_AVERAGE

	return ExecutePrometheusQueryRange(query, metric, startTime, endTime, time.Duration(tp))
}

func GetMemUsageAverageRange(ns string, pod string, ctnr string, tp int64, offSet int64, startTime time.Time, endTime time.Time) []models.PrometheusResult {
	var params = models.PrometheusQueryParams{
		Namespace:    ns,
		Pod:          pod,
		Container:    ctnr,
		TargetPeriod: BuiildTargetPeriod(tp),
		// Offset:       BuiildTargetPeriod(offSet),
	}

	query := models.BuildMemUsageAverageQuery(&params)
	metric := models.MetricType_MEMORY_USAGE_AVERAGE

	return ExecutePrometheusQueryRange(query, metric, startTime, endTime, time.Duration(tp))
}

func GetMemUsageAverage(ns string, pod string, ctnr string, tp int64, offSet int64, timeReq time.Time) []models.PrometheusResult {
	var params = models.PrometheusQueryParams{
		Namespace:    ns,
		Pod:          pod,
		Container:    ctnr,
		TargetPeriod: BuiildTargetPeriod(tp),
		// Offset:       BuiildTargetPeriod(offSet),
	}

	query := models.BuildMemUsageAverageQuery(&params)
	metric := models.MetricType_MEMORY_USAGE_AVERAGE

	return ExecutePrometheusQuery(query, metric, timeReq)
}

func GetResourceLimit(ns string, pod string, ctnr string, unit models.PrometheusUnit, timeReq time.Time) []models.PrometheusResult {
	var params = models.PrometheusQueryParams{
		Namespace: ns,
		Pod:       pod,
		Container: ctnr,
		Unit:      string(unit),
	}

	query := models.BuildResourceLimitQuery(&params)
	var metric models.MetricType
	if unit == models.PrometheusUnit_CORE {
		metric = models.MetricType_CPU_LIMIT
	} else {
		metric = models.MetricType_MEMORY_LIMIT
	}

	return ExecutePrometheusQuery(query, metric, timeReq)
}

func GetResourceRequest(ns string, pod string, ctnr string, unit models.PrometheusUnit, timeReq time.Time) []models.PrometheusResult {
	var params = models.PrometheusQueryParams{
		Namespace: ns,
		Pod:       pod,
		Container: ctnr,
		Unit:      string(unit),
	}

	query := models.BuildResourceRequestQuery(&params)
	var metric models.MetricType
	if unit == models.PrometheusUnit_CORE {
		metric = models.MetricType_CPU_REQUEST
	} else {
		metric = models.MetricType_MEMORY_REQUEST
	}

	return ExecutePrometheusQuery(query, metric, timeReq)
}

func GetRunningPods(instance string, ns string, ctnr string, timeReq time.Time) []models.PrometheusResult {
	var params = models.PrometheusQueryParams{
		Instance:  instance,
		Namespace: ns,
		Container: ctnr,
	}

	query := models.BuildRunningPodsQuery(&params)
	metric := models.MetricType_RUNNING_POD

	return ExecutePrometheusQuery(query, metric, timeReq)
}

func ProcessPrometheusMetricResult(result model.Value, metric models.MetricType) []models.PrometheusResult {
	var output []models.PrometheusResult
	var err error

	switch v := result.(type) {
	case model.Vector:
		// Vector
		// logger.AniLog.Infof("Result type %T\n", v)
		// Vector
		if len(v) == 0 {
			err = fmt.Errorf("no data found in Prometheus response")
		}

		for _, sample := range v {
			// Extraer el valor del metric map
			// logger.AniLog.Infof("Sample: %s, %s", sample, sample.Value)
			metricMap := sample.Metric
			namespace := string(metricMap["namespace"])
			pod := string(metricMap["pod"])
			container := string(metricMap["container"])
			phase := string(metricMap["phase"])
			uid := string(metricMap["uid"])

			// Crear una instancia de la estructura con los datos extraídos
			prometheusData := models.PrometheusResult{
				Timestamp:  float64(sample.Timestamp),
				Value:      float64(sample.Value),
				MetricType: metric,
				Namespace:  namespace,
				Pod:        pod,
				Container:  container,
				Phase:      phase,
				Uid:        uid,
			}

			// Agregar la estructura al slice de resultados
			output = append(output, prometheusData)
		}

	case model.Matrix:
		// Matrix
		// logger.AniLog.Infof("Result type %T\n", v)
		// Iterar a través de cada serie temporal
		for _, stream := range v {
			metricMap := stream.Metric
			namespace := string(metricMap["namespace"])
			pod := string(metricMap["pod"])
			container := string(metricMap["container"])
			phase := string(metricMap["phase"])
			uid := string(metricMap["uid"])
			for _, sample := range stream.Values {
				// logger.AniLog.Warn("Sample: ", sample)
				// Aquí procesas cada muestra con su timestamp y valor
				// metrics = append(metrics, PrometheusResult{
				//     Timestamp: sample.Timestamp.Time(),
				//     Value:     float64(sample.Value),
				//     Metric:    query,
				// })

				// Extraer el valor del metric map
				// logger.AniLog.Infof("Sample: %s, %s", sample, sample.Value)
				// metricMap := sample.Metric
				// namespace := string(metricMap["namespace"])
				// pod := string(metricMap["pod"])
				// container := string(metricMap["container"])
				// phase := string(metricMap["phase"])
				// uid := string(metricMap["uid"])

				// Crear una instancia de la estructura con los datos extraídos
				prometheusData := models.PrometheusResult{
					Timestamp:  float64(sample.Timestamp),
					Value:      float64(sample.Value),
					MetricType: metric,
					Namespace:  namespace,
					Pod:        pod,
					Container:  container,
					Phase:      phase,
					Uid:        uid,
				}

				// Agregar la estructura al slice de resultados
				output = append(output, prometheusData)
			}
		}

	case *model.Scalar:
		// Scalar
		err = fmt.Errorf(" Result type %T no implemented", v)

	default:
		// Default
		err = fmt.Errorf("unexpected result type: %T", v)
	}

	// Verify errors
	if err != nil {
		logger.PcmLog.Errorf("Error processing Prometheus data (%s): %v", metric, err)
	} else {
		// Convertir el resultado a JSON para imprimirlo
		_, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			logger.PcmLog.Errorf("Error converting Prometheus result to JSON: %v", err)
		}
	}

	return output
}

func BuiildTargetPeriod(num int64) string {
	minutes := math.Round(float64(num) / 60)
	rounded := minutes
	result := strconv.FormatFloat(rounded, 'f', -1, 64) + "m"
	return result
}
