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

	logger.MlModelInfoLog.Warn("Response from NwdafMlModelInfoRequestProcedure: ", response)

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

func NwdafMlModelInfoRequestProcedure() []models.MlModelInfoData {
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
	rows, err := db.Query(`
		SELECT uri, accuracy, nf_type, event_id, target_period 
		FROM ` + string(models.NwdafMLModelDB_ML_MODEL_INFO) + `;`)
	if err != nil {
		logger.MlModelInfoLog.Error("Error al consultar los registros: ", err)
		return nil
	}
	defer rows.Close()

	// Iterar sobre los resultados y mapearlos a una estructura
	var records []models.MlModelInfoData
	for rows.Next() {
		var mlmodels models.MlModelInfoData
		err := rows.Scan(&mlmodels.URI, &mlmodels.Accuracy, &mlmodels.NfType, &mlmodels.EventId, &mlmodels.TargetPeriod)
		if err != nil {
			logger.MlModelInfoLog.Error("Error al leer los registros: ", err)
			return nil
		}
		records = append(records, mlmodels)
	}

	if err := rows.Err(); err != nil {
		logger.MlModelInfoLog.Error("Error en el procesamiento de filas: ", err)
		return nil
	}

	logger.MlModelInfoLog.Info("Registros obtenidos con éxito")
	return records
}
