package producer

import (
	"net/http"

	"github.com/enable-intelligent-containerized-5g/openapi/models"
	"github.com/free5gc/nwdaf/internal/logger"
	"github.com/free5gc/nwdaf/pkg/factory"
	"github.com/free5gc/util/httpwrapper"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func HandleNwdafMlModelInfoRequest(request *httpwrapper.Request) *httpwrapper.Response {
	logger.MlModelInfoLog.Infoln("Handle MlModelInfoRequest")
	// nfInstanceId := request.Params["nfInstanceID"]

	response := NwdafMlModelInfoRequestProcedure()

	// logger.MlModelInfoLog.Warn("Response from NwdafMlModelInfoRequestProcedure: ", response)

	if response != nil {
		return httpwrapper.NewResponse(http.StatusOK, nil, response)
	} else {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "UNSPECIFIED",
		}
		return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
}

func NwdafMlModelInfoRequestProcedure() []models.MlModelData {
	logger.MlModelInfoLog.Infoln("Procedure MlModelInfoRequest")
	// Conectar a la base de datos SQLite
	sqldb := factory.NwdafConfig.Configuration.SqlLiteDB
	db, errCon := sql.Open("sqlite3", sqldb)
	if errCon != nil {
		logger.MlModelInfoLog.Error("Error al abrir la base de datos: ", errCon)
		return nil
	}
	defer db.Close()

	// Consultar todos los registros de la tabla 'records'
	selectSQL := `SELECT uri AS uri, accuracy AS accuracy, size AS size, nf_type AS nfType, event_id AS eventId, target_period AS targetPeriod FROM `+ string(models.NwdafMLModelDB_ML_MODEL_INFO) +`;`
	rows, err := db.Query(selectSQL)
	if err != nil {
		logger.MlModelInfoLog.Error("Error querying mlmodels data: ", err)
		return nil
	}
	defer rows.Close()

	// logger.MlModelInfoLog.Error(rows)

	// Iterar sobre los resultados y mapearlos a una estructura
	var mlmodels []models.MlModelData
	for rows.Next() {
		var mlmodel models.MlModelData
		err := rows.Scan(&mlmodel)
		if err != nil {
			logger.MlModelInfoLog.Error("Error reading records: ", err)
			return nil
		}
		mlmodels = append(mlmodels, mlmodel)
	}

	if err := rows.Err(); err != nil {
		logger.MlModelInfoLog.Error("Error proccesing rows: ", err)
		return nil
	}

	logger.MlModelInfoLog.Error("models: ",mlmodels)

	logger.MlModelInfoLog.Info("Registros obtenidos con éxito")
	return mlmodels
}
