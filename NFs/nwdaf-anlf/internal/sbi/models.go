package sbi

import (
	"time"

	"github.com/enable-intelligent-containerized-5g/openapi/models"
)

type NwdafMlModelTrainingRequest struct {
	TargetPeriod int64          `json:"targetPeriod,omitempty" yaml:"targetPeriod,omitempty" bson:"targetPeriod" validate:"required"`
	EventId      models.EventId `json:"eventId,omitempty" yaml:"eventId,omitempty" bson:"eventId,omitempty" validate:"required"`
	NfType       models.NfType  `json:"nfType,omitempty" yaml:"nfType,omitempty" bson:"nfType,omitempty" validate:"required"`
	StartTime    time.Time      `json:"startTime,omitempty" yaml:"startTime,omitempty" bson:"startTime,omitempty" validate:"required"`
}

// Estructura para almacenar los datos del CSV
type CsvData struct {
	Pod        string
	Container  string
	Timestamp1 int64
	CpuUsage1  float64
	Timestamp2 int64
	CpuUsage2  float64
	Timestamp3 int64
	CpuUsage3  float64
	Timestamp4 int64
	CpuUsage4  float64
}
