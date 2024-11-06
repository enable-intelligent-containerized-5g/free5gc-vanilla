package producer

import (
	// "embed"
	"net/http"

	"github.com/enable-intelligent-containerized-5g/openapi/models"
	"github.com/free5gc/nwdaf/internal/logger"
	"github.com/free5gc/nwdaf/internal/util"
	"github.com/free5gc/nwdaf/pkg/factory"
	"github.com/free5gc/util/httpwrapper"

	// "database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func HandleNwdafMlModelInfoRequest(request *httpwrapper.Request) *httpwrapper.Response {
	logger.MlModelInfoLog.Infoln("Handle MlModelInfoRequest")
	// nfInstanceId := request.Params["nfInstanceID"]

	response, err := NwdafMlModelInfoRequestProcedure()

	// logger.MlModelInfoLog.Warn("Response from NwdafMlModelInfoRequestProcedure: ", response)

	if err != nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  err.Error(),
		}
		return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	} else {
		return httpwrapper.NewResponse(http.StatusOK, nil, response)

	}
}

func NwdafMlModelInfoRequestProcedure() (models.MlModelDataResponse, error) {
	logger.MlModelInfoLog.Infoln("Procedure MlModelInfoRequest")

	// Conectar a la base de datos SQLite
	sqldb := factory.NwdafConfig.Configuration.SqlLiteDB
	db, errCon := util.OpenDatabase(sqldb)
	if errCon != nil {
		// logger.MlModelInfoLog.Errorf("Error al conectar a la base de datos: %v", errCon)
		return models.MlModelDataResponse{}, errCon
	}

	// Variable para almacenar los resultados
	var mlModels []util.MlModelDataTable
	// Consultar todos los registros de la tabla MlModelDataTable
	result := db.
		Preload("Event").
		Preload("NfType").
		Preload("Accuracy").
		// Preload(string(models.NwdafMLModelDB_NF_TYPE_TABLE_NAME)).
		// Preload(string(models.NwdafMLModelDB_ACCURACY_TABLE_NAME)).
		// Preload(string(models.NwdafMLModelDB_EVENT_ID_TABLE_NAME)).
		Find(&mlModels)
	if result.Error != nil {
		// logger.MlModelInfoLog.Errorf("Error al consultar la tabla: %v", result.Error)
		return models.MlModelDataResponse{}, result.Error
	}

	var mlModelList []models.MlModelData
	logger.MlModelInfoLog.Info("Registros obtenidos:")
	for _, model := range mlModels {
		mlModelData := models.MlModelData{
			URI:          model.URI,
			Size:         model.Size,
			TargetPeriod: model.TargetPeriod,
			NfType:       models.NfType(model.NfType.NfType),
			EventId:      models.EventId(model.Event.Event),
			Accuracy:     models.NwdafMlModelAccuracy(model.Accuracy.Accuracy),
		}

		// logger.MlModelInfoLog.Infof("ID: %d, URI: %s, Size: %d, TargetPeriod: %d, NfTypeID: %s, AccuracyID: %s, EventID: %s\n",
		// 	model.ID, model.URI, model.Size, model.TargetPeriod, model.NfTypeID, model.AccuracyID, model.EventID)

		mlModelList = append(mlModelList, mlModelData)
	}

	logger.MlModelInfoLog.Info("Registros obtenidos con éxito")

	return models.MlModelDataResponse{MlModels: mlModelList}, nil
}
