package producer

import (
	"net/http"

	"github.com/enable-intelligent-containerized-5g/openapi/models"
	"github.com/free5gc/nwdaf/internal/logger"
	"github.com/free5gc/nwdaf/internal/util"
	"github.com/free5gc/nwdaf/pkg/factory"
	"github.com/free5gc/util/httpwrapper"
	"gorm.io/gorm"

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

func SaveMlModelProcedure(mldata models.MlModelData) (models.MlModelDataResponse, bool, *models.ProblemDetails) {
	logger.MlModelTrainingLog.Info("Procedure SaveMlModel")

	// Conect to database
	sqldb := factory.NwdafConfig.Configuration.SqlLiteDB
	db, errCon := util.OpenDatabase(sqldb)
	if errCon != nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  errCon.Error(),
		}
		return models.MlModelDataResponse{}, false, problemDetails
	}

	// Validate EventId
	eventFound := util.EventTable{Event: mldata.EventId}
	errGetEvent := GetEventByName(&eventFound, db)
	if errGetEvent != nil {
		// logger.MlModelTrainingLog.Errorf("Event %s not found: %s", mldata.EventId, errGetEvent)
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  errGetEvent.Error(),
		}
		return models.MlModelDataResponse{}, false, problemDetails
		
	}

	// Validate Accuracy
	accuFound := util.AccuracyTable{Accuracy: mldata.Accuracy}
	errGetAccu := GetAccuracyByName(&accuFound, db)
	if errGetAccu != nil {
		// logger.MlModelTrainingLog.Errorf("Accuracy %s not found: %s", mldata.Accuracy, errGetAccu)
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  errGetAccu.Error(),
		}
		return models.MlModelDataResponse{}, false, problemDetails
	}

	// Validate NfType
	nfTypeFound := util.NFTypeTable{NfType: mldata.NfType}
	errGetNfType := GetNfTypeByName(&nfTypeFound, db)
	if errGetNfType != nil {
		// logger.MlModelTrainingLog.Errorf("NfType %s not found: %s", mldata.NfType, errGetNfType)
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  errGetNfType.Error(),
		}
		return models.MlModelDataResponse{}, false, problemDetails
	}

	// Create the struct
	mlModelTableRequest := util.MlModelDataTable{
		URI:          mldata.URI,
		Size:         mldata.Size,
		TargetPeriod: mldata.TargetPeriod,
		AccuracyID:   accuFound.ID,
		EventID:      eventFound.ID,
		NfTypeID:     nfTypeFound.ID,
		Accuracy:     accuFound,
		Event:        eventFound,
		NfType:       nfTypeFound,
	}
	errSaves := SaveMlmodel(&mlModelTableRequest, db)
	if errSaves != nil {
		// logger.MlModelTrainingLog.Errorf("MlModel not saved: %s", errSaves)
		problemDetails := &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  errSaves.Error(),
		}
		return models.MlModelDataResponse{}, false, problemDetails
	}

	mlmodelSaved := models.MlModelData{
		URI:          mlModelTableRequest.URI,
		Size:         mlModelTableRequest.Size,
		TargetPeriod: mlModelTableRequest.TargetPeriod,
		Accuracy:     mlModelTableRequest.Accuracy.Accuracy,
		NfType:       mlModelTableRequest.NfType.NfType,
		EventId:      mlModelTableRequest.Event.Event,
	}

	return models.MlModelDataResponse{MlModels: append([]models.MlModelData{}, mlmodelSaved)}, true, nil
}

// Función para recuperar un modelo ML por ID
func GetMlModelInfoByID(db *sql.DB, id int64) (models.MlModelData, error) {
	var modelInfo models.MlModelData
	// selectSQL := `SELECT uri AS uri, size AS size, accuracy AS accuracy, nf_type AS nfType, event_id AS eventId, target_period AS targetPeriod FROM `+ string(models.NwdafMLModelDB_ML_MODEL_INFO) + ` WHERE id = ?;`
	// row := db.QueryRow(selectSQL, id)

	// if row == nil {
	// 	logger.MlModelTrainingLog.Warn("row nil")
	//     return modelInfo, http.ErrAbortHandler // Retornar error
	// }

	// logger.MlModelTrainingLog.Warn("responde select by id: ", row)
	// err := row.Scan(&modelInfo.URI, &modelInfo.Size, &modelInfo.Accuracy, &modelInfo.NfType, &modelInfo.EventId, &modelInfo.TargetPeriod)
	// if err != nil {
	//     return models.MlModelData{}, err // Retornar error
	// }

	return modelInfo, nil // Retornar el objeto recuperado
}

func GetEventByName(event *util.EventTable, db *gorm.DB) error {
	err := db.Where(&event).First(&event).Error
	if err != nil {
		return err
	}

	return nil
}

func GetAccuracyByName(accuracy *util.AccuracyTable, db *gorm.DB) error {
	err := db.Where(&accuracy).First(&accuracy).Error
	if err != nil {
		return err
	}

	return nil
}

func GetNfTypeByName(nf *util.NFTypeTable, db *gorm.DB) error {
	err := db.Where(&nf).First(&nf).Error
	if err != nil {
		return err
	}

	return nil
}

func SaveMlmodel(mlModel *util.MlModelDataTable, db *gorm.DB) error {
	result := db.Create(&mlModel)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
