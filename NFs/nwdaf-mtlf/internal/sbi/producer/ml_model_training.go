package producer

import (
	"github.com/enable-intelligent-containerized-5g/openapi/models"
	"github.com/free5gc/nwdaf/internal/logger"
	"github.com/free5gc/nwdaf/pkg/factory"
	"github.com/free5gc/util/httpwrapper"
	"net/http"
	// "encoding/json"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// type MlModelInfo struct {
// 	URI          string `json:"uri,omitempty" yaml:"uri" bson:"uri" mapstructure:"uri"`
// 	Accuracy     string `json:"accuracy,omitempty" yaml:"accuracy" bson:"accuracy" mapstructure:"accuracy"`
// 	NF           string `json:"nf,omitempty" yaml:"nf" bson:"nf" mapstructure:"nf"`
// 	TargetPeriod string `json:"targetPeriod,omitempty" yaml:"targetPeriod" bson:"targetPeriod" mapstructure:"targetPeriod"`
// 	EventId      string `json:"eventId,omitempty" yaml:"eventId" bson:"eventId" mapstructure:"eventId"`
// }

func HandleSaveMlModel(request *httpwrapper.Request) *httpwrapper.Response {
	logger.MlModelTrainingLog.Warn("request.Body in HandleSaveMlModel: ", request.Body)

	mlmodeldata, ok := request.Body.(models.MlModelInfoData)
	if !ok {
		httpwrapper.NewResponse(http.StatusBadRequest, nil, nil)
	}

	response, problemDetails := SaveMlModelProcedure(mlmodeldata)
	if response != nil {
		logger.MlModelTrainingLog.Warn("CreateSubscription success")
		return httpwrapper.NewResponse(http.StatusCreated, nil, response)
	} else if problemDetails != nil {
		logger.MlModelTrainingLog.Warn("CreateSubscription failed")
		return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}

	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	logger.MlModelTrainingLog.Error("CreateSubscription failed")
	return httpwrapper.NewResponse(http.StatusForbidden, nil, problemDetails)
}

func SaveMlModelProcedure(mldata models.MlModelInfoData) (sql.Result, *models.ProblemDetails) {
	// Conectar a la base de datos SQLite
	sqldb := factory.NwdafConfig.Configuration.SqlLiteDB
	db, errCon := sql.Open("sqlite3", sqldb)
	if errCon != nil {
		ProblemSql := &models.ProblemDetails{
			Status: http.StatusForbidden,
			Cause:  errCon.Error(),
		}
		logger.MlModelTrainingLog.Error("Error al abrir la base de datos: ", errCon)
		return nil, ProblemSql
	}
	defer db.Close()

	// Insertar el registro en la tabla 'records'
	ProblemPut := &models.ProblemDetails{}
	putData, err := db.Exec(`
		INSERT INTO `+ string(models.NwdafMLModelDB_ML_MODEL_INFO) + ` (uri, accuracy, nf, event, target_period) 
		VALUES (?, ?, ?, ?, ?);`,
		mldata.URI, mldata.Accuracy, mldata.NfType, mldata.TargetPeriod, mldata.EventId)
	if err != nil {
		ProblemPut = &models.ProblemDetails{
			Status: http.StatusForbidden,
			Cause:  err.Error(),
		}
		logger.MlModelTrainingLog.Error("Error al insertar el registro: ", err)
		return putData, ProblemPut
	}

	logger.MlModelTrainingLog.Info("Registro insertado con éxito")
	db.Close()
	return putData, nil
}
