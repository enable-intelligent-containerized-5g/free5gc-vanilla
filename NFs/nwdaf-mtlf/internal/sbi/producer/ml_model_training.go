package producer

import (
	"github.com/enable-intelligent-containerized-5g/openapi/models"
	"github.com/free5gc/nwdaf/internal/logger"
	"github.com/free5gc/nwdaf/internal/util"
	"github.com/free5gc/nwdaf/pkg/factory"
	"github.com/free5gc/util/httpwrapper"
	"net/http"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)


func HandleSaveMlModel(request *httpwrapper.Request) *httpwrapper.Response {
	logger.MlModelTrainingLog.Info("Handle SaveMlModel")

    mlmodeldata, ok := request.Body.(models.MlModelData)
    if !ok {
        return httpwrapper.NewResponse(http.StatusForbidden, nil, "The request body is't type MlModelData")
    }

	logger.MlModelTrainingLog.Warn("Handle ModelData: ", mlmodeldata)

	putData, created, problemDetails := SaveMlModelProcedure(mlmodeldata)
	if created {
		logger.MlModelTrainingLog.Info("SaveMlModel success")
		return httpwrapper.NewResponse(http.StatusCreated, nil, putData)
	} else if problemDetails != nil {
		logger.MlModelTrainingLog.Errorf("SaveMlModel failed: %s", problemDetails.Cause)
		return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}

	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	logger.MlModelTrainingLog.Error("SaveMlModel failed")
	return httpwrapper.NewResponse(http.StatusForbidden, nil, problemDetails)
}

func SaveMlModelProcedure(mldata models.MlModelData) (models.MlModelData, bool, *models.ProblemDetails) {
	logger.MlModelTrainingLog.Info("Procedure SaveMlModel")

	Conectar a la base de datos SQLite
	sqldb := factory.NwdafConfig.Configuration.SqlLiteDB
	db, errCon := util.OpenDatabase(sqldb)
	if errCon != nil {
		ProblemSql := &models.ProblemDetails{
			Status: http.StatusForbidden,
			Cause:  errCon.Error(),
		}
		logger.MlModelTrainingLog.Errorf("Error to open the database %s: %s",sqldb, errCon)
		return models.MlModelData{}, false, ProblemSql
	}
	defer db.Close()

	SaveMLModel
	logger.MlModelTrainingLog.Info("Procedure SaveMlModel: ", mldata)
	ProblemPut := &models.ProblemDetails{}
	putData, err := db.Exec(`
		INSERT INTO `+ string(models.NwdafMLModelDB_ML_MODEL_INFO) + ` (uri, accuracy, nf_type, event_id, target_period, size) 
		VALUES (?, ?, ?, ?, ?, ?);`,
		mldata.URI, mldata.Accuracy, mldata.NfType, mldata.EventId, mldata.TargetPeriod, mldata.Size)
	if err != nil {
		ProblemPut = &models.ProblemDetails{
			Status: http.StatusForbidden,
			Cause:  err.Error(),
		}
		// logger.MlModelTrainingLog.Error("Error to insert the MlModel: ", err)
		return models.MlModelData{}, false, ProblemPut
	}

	lastInsertId, err := putData.LastInsertId()
	if err!=nil {
		ProblemPut = &models.ProblemDetails{
			Status: http.StatusForbidden,
			Cause:  err.Error(),
		}
		// logger.MlModelTrainingLog.Errorf("Error geting the inserted MlModel id: %v", err)
		return models.MlModelData{}, true, ProblemPut
	}

	logger.MlModelTrainingLog.Warn("Before GetByID: ", lastInsertId)
	retrievedModelInfo, err := GetMlModelInfoByID(db, lastInsertId)
	if err != nil {
		ProblemPut = &models.ProblemDetails{
			Status: http.StatusForbidden,
			Cause:  err.Error(),
		}
        // logger.MlModelTrainingLog.Errorf("Error geting the inserted MlModel: %v", err)
		return retrievedModelInfo, true, ProblemPut
    }

	logger.MlModelTrainingLog.Info("Mlmodel saved succesful", retrievedModelInfo)
	db.Close()
	return retrievedModelInfo, true, nil
	return models.MlModelData{}, true, nil


}


// Función para recuperar un modelo ML por ID
func GetMlModelInfoByID(db *sql.DB, id int64) (models.MlModelData, error) {
    var modelInfo models.MlModelData
    selectSQL := `SELECT uri AS uri, size AS size, accuracy AS accuracy, nf_type AS nfType, event_id AS eventId, target_period AS targetPeriod FROM `+ string(models.NwdafMLModelDB_ML_MODEL_INFO) + ` WHERE id = ?;`
    row := db.QueryRow(selectSQL, id)


	if row == nil {
		logger.MlModelTrainingLog.Warn("row nil")
        return modelInfo, http.ErrAbortHandler // Retornar error
    }

	logger.MlModelTrainingLog.Warn("responde select by id: ", row)
    err := row.Scan(&modelInfo.URI, &modelInfo.Size, &modelInfo.Accuracy, &modelInfo.NfType, &modelInfo.EventId, &modelInfo.TargetPeriod)
    if err != nil {
        return models.MlModelData{}, err // Retornar error
    }

    return modelInfo, nil // Retornar el objeto recuperado
}