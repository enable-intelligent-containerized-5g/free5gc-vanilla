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

    mlmodelinfo, ok := request.Body.(models.MlModelData)
    if !ok {
        return httpwrapper.NewResponse(http.StatusForbidden, nil, "The request body is't type MlModelData")
    }

	putData, created, problemDetails := SaveMlModelProcedure(mlmodelinfo)
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

	// Conectar a la base de datos SQLite
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

	// SaveMLModel
	ProblemPut := &models.ProblemDetails{}
	putData, err := db.Exec(`
		INSERT INTO `+ string(models.NwdafMLModelDB_ML_MODEL_INFO) + ` (uri, accuracy, nf_type, event_id, target_period) 
		VALUES (?, ?, ?, ?, ?);`,
		string(mldata.URI), mldata.Accuracy, string(mldata.NfType), string(mldata.EventId), string(mldata.TargetPeriod))
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
}


// Función para recuperar un modelo ML por ID
func GetMlModelInfoByID(db *sql.DB, id int64) (models.MlModelData, error) {
    var modelInfo models.MlModelData
    selectSQL := `SELECT uri, accuracy, nf_type, event_id, target_period FROM `+ string(models.NwdafMLModelDB_ML_MODEL_INFO) + ` WHERE id = ?;`
    row := db.QueryRow(selectSQL, id)

    err := row.Scan(&modelInfo.URI, &modelInfo.Accuracy, &modelInfo.NfType, &modelInfo.EventId, &modelInfo.TargetPeriod)
    if err != nil {
        return models.MlModelData{}, err // Retornar error
    }

    return modelInfo, nil // Retornar el objeto recuperado
}