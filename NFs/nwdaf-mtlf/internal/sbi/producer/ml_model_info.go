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

	response, err := NwdafMlModelInfoRequestProcedure()

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

	// connect to database
	sqldb := factory.NwdafConfig.Configuration.SqlLiteDB
	db, errCon := util.OpenDatabase(sqldb)
	if errCon != nil {
		return models.MlModelDataResponse{}, errCon
	}

	// Results
	var mlModels []models.MlModelDataTable
	// Get results
	result := db.
		Preload(string(models.NwdafMLModelDB_EVENT_ID_KEY)).
		Preload(string(models.NwdafMLModelDB_NF_TYPE_KEY)).
		Preload(string(models.NwdafMLModelDB_ACCURACY_KEY)).
		Find(&mlModels)
	if result.Error != nil {
		return models.MlModelDataResponse{}, result.Error
	}

	var mlModelList []models.MlModelData
	for _, model := range mlModels {
		mlModelData := models.MlModelData{
			URI:          model.URI,
			Size:         model.Size,
			TargetPeriod: model.TargetPeriod,
			Confidence:   model.Confidence,
			NfType:       model.NfType.NfType,
			EventId:      model.Event.Event,
			Accuracy:     model.Accuracy.Accuracy,
		}
		mlModelList = append(mlModelList, mlModelData)
	}

	// logger.MlModelInfoLog.Info("Registros obtenidos con éxito")

	return models.MlModelDataResponse{MlModels: mlModelList}, nil
}
