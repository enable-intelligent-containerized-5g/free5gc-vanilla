package producer

import (
	// "embed"
	"net/http"

	"github.com/enable-intelligent-containerized-5g/openapi/models"
	"github.com/free5gc/nwdaf/internal/logger"
	"github.com/free5gc/nwdaf/internal/util"
	"github.com/free5gc/nwdaf/pkg/factory"
	"github.com/free5gc/util/httpwrapper"

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
		return models.MlModelDataResponse{}, errCon
	}

	// Variable para almacenar los resultados
	var mlModels []util.MlModelDataTable
	// Consultar todos los registros de la tabla MlModelDataTable
	result := db.
		Preload(string(models.NwdafMLModelDB_EVENT_ID_KEY)).
		Preload(string(models.NwdafMLModelDB_NF_TYPE_KEY)).
		Preload(string(models.NwdafMLModelDB_ACCURACY_KEY)).
		Find(&mlModels)
	if result.Error != nil {
		return models.MlModelDataResponse{}, result.Error
	}

	var mlModelList []models.MlModelData
	logger.MlModelInfoLog.Info("Registros obtenidos:")
	for _, model := range mlModels {
		mlModelData := models.MlModelData{
			URI:          model.URI,
			Size:         model.Size,
			TargetPeriod: model.TargetPeriod,
			NfType:       model.NfType.NfType,
			EventId:      model.Event.Event,
			Accuracy:     model.Accuracy.Accuracy,
		}
		mlModelList = append(mlModelList, mlModelData)
	}

	logger.MlModelInfoLog.Info("Registros obtenidos con éxito")

	return models.MlModelDataResponse{MlModels: mlModelList}, nil
}
