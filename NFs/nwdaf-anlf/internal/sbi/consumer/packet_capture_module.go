package consumer

import (
	"encoding/json"
	"fmt"
	"math"
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
	Namespace  string     `json:"namespace"`
	Pod        string     `json:"pod"`
	Container  string     `json:"container"`
	Phase      string     `json:"phase"`
	Uid        string     `json:"uid"`
}

func NewPrometheusResult() PrometheusResult {
	return PrometheusResult{
		Timestamp:  0,
		Value:      0,
		MetricType: "",
		Namespace:  "",
		Pod:        "",
		Container:  "",
		Phase:      "",
		Uid:        "",
	}
}

type MetricType string

const (
	MetricType_CPU_USAGE            MetricType = "cpu-usage"
	MetricType_MEMORY_USAGE         MetricType = "mem-usage"
	MetricType_CPU_USAGE_AVERAGE    MetricType = "cpu-average"
	MetricType_MEMORY_USAGE_AVERAGE MetricType = "mem-average"
	MetricType_CPU_LIMIT            MetricType = "cpu-limit"
	MetricType_MEMORY_LIMIT         MetricType = "men-limit"
	MetricType_CPU_REQUEST          MetricType = "cpu-request"
	MetricType_MEMORY_REQUEST       MetricType = "men-request"
	MetricType_RUNNING_POD          MetricType = "pod-status"
)

type KubernetesPhase string

const (
	KubernetesPhase_PENDING   KubernetesPhase = "Pending"
	KubernetesPhase_RUNNING   KubernetesPhase = "Running"
	KubernetesPhase_SUCCEEDED KubernetesPhase = "Succeeded"
	KubernetesPhase_FAILED    KubernetesPhase = "Failed"
	KubernetesPhase_UNKNOWN   KubernetesPhase = "Unknown"
)

type PrometheusUnit string

const (
	PrometheusUnit_CORE PrometheusUnit = "core"
	PrometheusUnit_BYTE PrometheusUnit = "byte"
)

type PrometheusQueryParams struct {
	Namespace    string `json:"namespace"`
	Pod          string `json:"pod"`
	Container    string `json:"container"`
	TargetPeriod string `json:"targetPeriod"`
	Offset       string `json:"offset"`
	Unit         string `json:"unit"`
	Instance     string `json:"instance"`
	Phase        string `json:"phase"`
}

// Memory rate (OK)
func BuildMemRateQuery(p *PrometheusQueryParams) string {
	return fmt.Sprintf(`sum(rate(container_memory_usage_bytes{namespace="%s", pod="%s", container="%s", container!~".*wait-.*"}[%s] offset %s)) by (pod, container)`,
		p.Namespace, p.Pod, p.Container, p.TargetPeriod, p.Offset)
}

// Memory Usage (OK)
func BuildMemUsageQuery(p *PrometheusQueryParams) string {
	return fmt.Sprintf(`avg(container_memory_usage_bytes{namespace="%s", pod="%s", container="%s", container!~"wait-.*"}) by (pod, container)`,
		p.Namespace, p.Pod, p.Container)
}

// CPU Usage Average (OK)
func BuildCpuUsageAverageQuery(p *PrometheusQueryParams) string {
	var offsetQuery string
	offSet := p.Offset
	if strings.TrimSpace(offSet) != "" {
		offsetQuery = fmt.Sprintf(` offset %s`, p.Offset)
	}
	return fmt.Sprintf(`avg(rate(container_cpu_usage_seconds_total{namespace="%s", pod="%s", container="%s", container!~".*wait-.*"}[%s]%s)) by (pod, container)`,
		p.Namespace, p.Pod, p.Container, p.TargetPeriod, offsetQuery)
}

// Memory Usage average (OK)
func BuildMemUsageAverageQuery(p *PrometheusQueryParams) string {
	var offsetQuery string
	offSet := p.Offset
	if strings.TrimSpace(offSet) != "" {
		offsetQuery = fmt.Sprintf(` offset %s`, p.Offset)
	}
	return fmt.Sprintf(`avg(avg_over_time(container_memory_usage_bytes{namespace="%s", pod="%s", container="%s", container!~"wait-.*"}[%s]%s)) by (pod, container)`,
		p.Namespace, p.Pod, p.Container, p.TargetPeriod, offsetQuery)
}

// CPU and Memory resources request (OK)
func BuildResourceRequestQuery(p *PrometheusQueryParams) string {
	return fmt.Sprintf(`avg(kube_pod_container_resource_requests{namespace="%s", pod="%s", container="%s",container!~"wait-.*", unit="%s"}) by (pod, container)`,
		p.Namespace, p.Pod, p.Container, p.Unit)
}

// CPU and Momory resources limit (OK)
func BuildResourceLimitQuery(p *PrometheusQueryParams) string {
	return fmt.Sprintf(`avg(kube_pod_container_resource_limits{namespace="%s", pod="%s", container="%s",container!~"wait-.*", unit="%s"}) by (pod, container)`,
		p.Namespace, p.Pod, p.Container, p.Unit)
}

// Pods running
func BuildRunningPodsQuery(p *PrometheusQueryParams) string {
	var ctnrQuery string
	ctnr := p.Container
	if strings.TrimSpace(ctnr) != "" {
		ctnrQuery = fmt.Sprintf(`, container="%s"`, p.Container)
	}
	return fmt.Sprintf(`kube_pod_container_status_running{instance="%s", namespace="%s"%s}`,
		p.Instance, p.Namespace, ctnrQuery)
}

// Pods by Phase
func BuildPodsByStatusQuery(p *PrometheusQueryParams) string {
	var phaseQuery string
	ctnr := p.Phase
	if strings.TrimSpace(ctnr) != "" {
		phaseQuery = fmt.Sprintf(`, phase="%s"`, p.Phase)
	}
	return fmt.Sprintf(`kube_pod_status_phase{instance="%s", namespace="%s"%s}`,
		p.Instance, p.Namespace, phaseQuery)
}

func CreateClient() (client api.Client, err error) {
	// Get PcmUri
	pcmUri := factory.NwdafConfig.Configuration.OamUri

	// Create the Prometheus Client
	client, err = api.NewClient(api.Config{
		Address: pcmUri,
	})
	if err != nil {
		return client, fmt.Errorf(" Error creating Prometheus client: %s", err)
	}

	return client, err
}

func ExecutePrometheusQuery(query string, metric MetricType) []PrometheusResult {
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

	metrics := ProcessPrometheusMetricResult(result, metric)

	if metrics == nil {
		var value PrometheusResult
		metrics = append(metrics, value)
	}

	return metrics
}

func GetPodsByPhase(instance string, ns string, phase KubernetesPhase) []PrometheusResult {
	var params = PrometheusQueryParams{
		Instance:  instance,
		Namespace: ns,
		Phase:     string(phase),
	}

	query := BuildPodsByStatusQuery(&params)
	metric := MetricType_RUNNING_POD

	return ExecutePrometheusQuery(query, metric)
}

func GetCpuUsageAverage(ns string, pod string, ctnr string, tp int64, offSet int64) []PrometheusResult {
	var params = PrometheusQueryParams{
		Namespace:    ns,
		Pod:          pod,
		Container:    ctnr,
		TargetPeriod: BuiildTargetPeriod(tp),
		Offset:       BuiildTargetPeriod(offSet),
	}

	query := BuildCpuUsageAverageQuery(&params)
	metric := MetricType_CPU_USAGE_AVERAGE

	return ExecutePrometheusQuery(query, metric)
}

func GetMemUsageAverage(ns string, pod string, ctnr string, tp int64, offSet int64) []PrometheusResult {
	var params = PrometheusQueryParams{
		Namespace:    ns,
		Pod:          pod,
		Container:    ctnr,
		TargetPeriod: BuiildTargetPeriod(tp),
		Offset:       BuiildTargetPeriod(offSet),
	}

	query := BuildMemUsageAverageQuery(&params)
	metric := MetricType_MEMORY_USAGE_AVERAGE

	return ExecutePrometheusQuery(query, metric)
}

func GetResourceLimit(ns string, pod string, ctnr string, unit PrometheusUnit) []PrometheusResult {
	var params = PrometheusQueryParams{
		Namespace: ns,
		Pod:       pod,
		Container: ctnr,
		Unit:      string(unit),
	}

	query := BuildResourceLimitQuery(&params)
	var metric MetricType
	if unit == PrometheusUnit_CORE {
		metric = MetricType_CPU_LIMIT
	} else {
		metric = MetricType_MEMORY_LIMIT
	}

	return ExecutePrometheusQuery(query, metric)
}

func GetResourceRequest(ns string, pod string, ctnr string, unit PrometheusUnit) []PrometheusResult {
	var params = PrometheusQueryParams{
		Namespace: ns,
		Pod:       pod,
		Container: ctnr,
		Unit:      string(unit),
	}

	query := BuildResourceRequestQuery(&params)
	var metric MetricType
	if unit == PrometheusUnit_CORE {
		metric = MetricType_CPU_REQUEST
	} else {
		metric = MetricType_MEMORY_REQUEST
	}

	return ExecutePrometheusQuery(query, metric)
}

func GetRunningPods(instance string, ns string, ctnr string) []PrometheusResult {
	var params = PrometheusQueryParams{
		Instance:  instance,
		Namespace: ns,
		Container: ctnr,
	}

	query := BuildRunningPodsQuery(&params)
	metric := MetricType_RUNNING_POD

	return ExecutePrometheusQuery(query, metric)
}

func ProcessPrometheusMetricResult(result model.Value, metric MetricType) []PrometheusResult {
	var output []PrometheusResult
	var err error

	switch v := result.(type) {
	case model.Vector:
		logger.AniLog.Infof("Result type %T", v)
		// Vector
		if len(v) == 0 {
			err = fmt.Errorf("no data found in Prometheus response")
		}

		for _, sample := range v {
			// Extraer el valor del metric map
			logger.AniLog.Infof("Sample: %s, %s", sample, sample.Value)
			metricMap := sample.Metric
			namespace := string(metricMap["namespace"])
			pod := string(metricMap["pod"])
			container := string(metricMap["container"])
			phase := string(metricMap["phase"])
			uid := string(metricMap["uid"])

			// Crear una instancia de la estructura con los datos extraídos
			prometheusData := PrometheusResult{
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

	case *model.Scalar:
		// Scalar
		err = fmt.Errorf(" Result type %T no implemented", v)

	default:
		// Default
		err = fmt.Errorf("unexpected result type: %T", v)
	}

	// Verify errors
	if err != nil {
		logger.PcmLog.Errorf("Error processing Prometheus data: %v", err)
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
