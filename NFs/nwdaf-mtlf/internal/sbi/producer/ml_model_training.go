package producer

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/enable-intelligent-containerized-5g/openapi/models"
	"github.com/free5gc/nwdaf/internal/logger"
	"github.com/free5gc/nwdaf/internal/util"
	"github.com/free5gc/nwdaf/pkg/factory"
	"github.com/free5gc/util/httpwrapper"
	"gorm.io/gorm"

	_ "github.com/mattn/go-sqlite3"
)

func HandleSaveMlModel(request *httpwrapper.Request) *httpwrapper.Response {
	logger.MlModelTrainingLog.Info("Handle SaveMlModel")

	mlmodeldata, ok := request.Body.(models.MlModelData)
	if !ok {
		return httpwrapper.NewResponse(http.StatusForbidden, nil, "The request body is't type MlModelData")
	}

	putData, created, problemDetails := SaveMlModelProcedure(mlmodeldata)
	if created {
		// logger.MlModelTrainingLog.Info("SaveMlModel success")
		return httpwrapper.NewResponse(http.StatusCreated, nil, putData)
	} else if problemDetails != nil {
		// logger.MlModelTrainingLog.Errorf("SaveMlModel failed: %s", problemDetails.Cause)
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

	// Validate Size
	if mldata.Size <= 0 {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  fmt.Sprintf("The Size must be greater than 0, but got %d", mldata.Size),
		}
		return models.MlModelDataResponse{}, false, problemDetails
	}

	// Validate TargetPeriod
	if mldata.TargetPeriod <= 0 {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  fmt.Sprintf("The TargetPeriod must be greater than 0, but got %d", mldata.TargetPeriod),
		}
		return models.MlModelDataResponse{}, false, problemDetails
	}

	// Validate URI
	if strings.TrimSpace(mldata.URI) == "" {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  "The URI cannot be empty",
		}
		return models.MlModelDataResponse{}, false, problemDetails
	}

	// Validate EventId
	eventFound := models.EventTable{Event: mldata.EventId}
	errGetEvent := GetEventByName(&eventFound, db)
	if errGetEvent != nil {
		// logger.MlModelTrainingLog.Errorf("Event %s not found: %s", mldata.EventId, errGetEvent)
		return models.MlModelDataResponse{}, false, errGetEvent

	}

	// Validate Accuracy
	accuFound := models.AccuracyTable{Accuracy: mldata.Accuracy}
	errGetAccu := GetAccuracyByName(&accuFound, db)
	if errGetAccu != nil {
		// logger.MlModelTrainingLog.Errorf("Accuracy %s not found: %s", mldata.Accuracy, errGetAccu)
		return models.MlModelDataResponse{}, false, errGetAccu
	}

	// Validate NfType
	nfTypeFound := models.NFTypeTable{NfType: mldata.NfType}
	errGetNfType := GetNfTypeByName(&nfTypeFound, db)
	if errGetNfType != nil {
		// logger.MlModelTrainingLog.Errorf("NfType %s not found: %s", mldata.NfType, errGetNfType)
		return models.MlModelDataResponse{}, false, errGetNfType
	}

	// Create the struct
	mlModelTableRequest := models.MlModelDataTable{
		URI:          mldata.URI,
		Size:         mldata.Size,
		TargetPeriod: mldata.TargetPeriod,
		Confidence:   mldata.Confidence,
		AccuracyID:   accuFound.ID,
		EventID:      eventFound.ID,
		NfTypeID:     nfTypeFound.ID,
		Accuracy:     accuFound,
		Event:        eventFound,
		NfType:       nfTypeFound,
	}
	errSaving := SaveMlmodel(&mlModelTableRequest, db)
	if errSaving != nil {
		// logger.MlModelTrainingLog.Errorf("MlModel not saved: %s", errSaves)
		return models.MlModelDataResponse{}, false, errSaving
	}

	var model2 models.MlModelDataTable
	errGetMlModel := GetMlModelById(&model2, db, mlModelTableRequest.ID)
	if errGetMlModel != nil {
		// logger.MlModelTrainingLog.Errorf("MlModel not found: %s", errSaves)
		return models.MlModelDataResponse{}, false, errGetMlModel
	}

	mlmodelSaved := models.MlModelData{
		URI:          mlModelTableRequest.URI,
		Size:         mlModelTableRequest.Size,
		TargetPeriod: mlModelTableRequest.TargetPeriod,
		Confidence:   mlModelTableRequest.Confidence,
		Accuracy:     mlModelTableRequest.Accuracy.Accuracy,
		NfType:       mlModelTableRequest.NfType.NfType,
		EventId:      mlModelTableRequest.Event.Event,
	}

	return models.MlModelDataResponse{MlModels: append([]models.MlModelData{}, mlmodelSaved)}, true, nil
}

// Get MlModel by ID
func GetMlModelById(mlModel *models.MlModelDataTable, db *gorm.DB, id int64) *models.ProblemDetails {
	result := db.First(&mlModel, id) // Search by ID = 1
	if result.Error != nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  fmt.Sprintf("MlModel with id %d not found", id),
		}
		return problemDetails
	}

	return nil
}

func GetEventByName(event *models.EventTable, db *gorm.DB) *models.ProblemDetails {
	err := db.Where(&event).First(&event).Error
	if err != nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  fmt.Sprintf("EventId %s not found", event.Event),
		}
		return problemDetails
	}

	return nil
}

func GetAccuracyByName(accuracy *models.AccuracyTable, db *gorm.DB) *models.ProblemDetails {
	err := db.Where(&accuracy).First(&accuracy).Error
	if err != nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  fmt.Sprintf("Accuracy %s not found", accuracy.Accuracy),
		}
		return problemDetails
	}

	return nil
}

func GetNfTypeByName(nf *models.NFTypeTable, db *gorm.DB) *models.ProblemDetails {
	err := db.Where(&nf).First(&nf).Error
	if err != nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  fmt.Sprintf("NfType %s not found", nf.NfType),
		}
		return problemDetails
	}

	return nil
}

func SaveMlmodel(mlModel *models.MlModelDataTable, db *gorm.DB) *models.ProblemDetails {
	result := db.Create(&mlModel)
	if result.Error != nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "Could not save the model to the database",
		}
		return problemDetails
	}

	return nil
}
