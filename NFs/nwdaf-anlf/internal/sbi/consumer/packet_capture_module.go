package consumer

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"context"

	"github.com/free5gc/nwdaf/internal/logger"
	"github.com/free5gc/nwdaf/pkg/factory"

	"github.com/prometheus/client_golang/api" // go get github.com/prometheus/client_golang/api
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type PrometheusResult struct {
	Timestamp  float64    `json:"timestamp"`
	Value      float64    `json:"value"`
	MetricType MetricType `json:"metric"`
	Pod        string     `json:"pod"`
	Container  string     `json:"container"`
}

func NewPrometheusResult() PrometheusResult {
	return PrometheusResult{
		Timestamp:  0,
		Value:      0,
		MetricType: "",
		Pod:        "",
		Container:  "",
	}
}

type MetricType string

const (
	MetricType_CPU_USAGE    MetricType = "CpuUsage"
	MetricType_MEMORY_USAGE MetricType = "MemUsage"
	MetricType_CPU_LIMIT    MetricType = "CpuLimit"
	MetricType_MEMORY_LIMIT MetricType = "MenLimit"
)

type PrometheusQuery string

const (
	PrometheusQuery_CPU_USAGE    PrometheusQuery = `sum(rate(container_cpu_usage_seconds_total{container=~"POD",container!~"wait-.*",}[TARGETPERIOD] offset OFFSET)) by (pod, container)`
	PrometheusQuery_MEMORY_USAGE PrometheusQuery = ``
	PrometheusQuery_CPU_LIMIT    PrometheusQuery = ``
	PrometheusQuery_MEMORY_LIMIT PrometheusQuery = ``
)

func CreateClient() (client api.Client, err error) {
	// Get PcmUri
	pcmUri := factory.NwdafConfig.Configuration.PcmUri

	// Create the Prometheus Client
	client, err = api.NewClient(api.Config{
		Address: pcmUri,
	})
	if err != nil {
		return client, fmt.Errorf(" Error creating Prometheus client: %s", err)
	}

	return client, err
}

func ExecutePrometheusQuery(query PrometheusQuery, metric MetricType) []PrometheusResult {
	// Get PcmUri
	client, errClient := CreateClient()
	if errClient != nil {
		logger.PcmLog.Error(errClient)
	}

	apiClient := v1.NewAPI(client)

	// Definir el contexto y el timeout para la consulta
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Realizar una consulta para obtener el uso de CPU en tiempo real
	result, warnings, err := apiClient.Query(ctx, string(query), time.Now())
	if err != nil {
		logger.PcmLog.Errorf("Error in the Request: %v", err)
	}
	if len(warnings) > 0 {
		logger.PcmLog.Warnf("Warnings: %v", warnings)
	}

	resultMap, err := ProcessPrometheusResult(result, metric)
	if err != nil {
		logger.PcmLog.Errorf("Error processing Prometheus data: %v", err)
	} else {
		// Convertir el resultado a JSON para imprimirlo
		_, err := json.MarshalIndent(resultMap, "", "  ")
		if err != nil {
			logger.PcmLog.Errorf("Error converting Prometheus result to JSON: %v", err)
		}
	}

	return resultMap
}

func GetCpuUsage(pod string, tp int64, offSet int64) []PrometheusResult {
	metric := MetricType_CPU_USAGE
	// Reemplazar "POD" y "TargetPeriod" en la consulta
	query := PrometheusQuery_CPU_USAGE
	query = PrometheusQuery(strings.Replace(string(query), "POD", pod, -1))
	query = PrometheusQuery(strings.Replace(string(query), "TARGETPERIOD", BuiildTargetPeriod(tp), -1))
	query = PrometheusQuery(strings.Replace(string(query), "OFFSET", BuiildTargetPeriod(offSet), -1))

	result := ExecutePrometheusQuery(query, metric)

	if result == nil {
		var results []PrometheusResult
		return append(results, NewPrometheusResult())
	}
	return result
}

func ProcessPrometheusResult(result model.Value, metric MetricType) ([]PrometheusResult, error) {
	var output []PrometheusResult

	switch v := result.(type) {
	case model.Vector:
		// Vector
		if len(v) == 0 {
			return nil, fmt.Errorf("no data found in Prometheus response")
		}

		for _, sample := range v {
			// Extraer el valor del metric map
			metricMap := sample.Metric
			pod := string(metricMap["pod"])
			container := string(metricMap["container"])

			// Crear una instancia de la estructura con los datos extraídos
			prometheusData := PrometheusResult{
				Timestamp:  float64(sample.Timestamp),
				Value:      float64(sample.Value),
				MetricType: metric,
				Pod:        pod,
				Container:  container,
			}

			// Agregar la estructura al slice de resultados
			output = append(output, prometheusData)
		}

	case *model.Scalar:
		// Scalar
		return nil, fmt.Errorf(" Result type %T no implemented", v)
		// logger.PcmLog.Info("Scalar")
		// entry := map[string]interface{}{
		// 	"value":     v.Value,
		// 	"timestamp": v.Timestamp,
		// }
		// output = append(output, entry)

	default:
		// Default
		return nil, fmt.Errorf("unexpected result type: %T", v)
	}

	return output, nil
}

func BuiildTargetPeriod(num int64) string {
	return strconv.FormatInt(num, 10) + "m"
}
