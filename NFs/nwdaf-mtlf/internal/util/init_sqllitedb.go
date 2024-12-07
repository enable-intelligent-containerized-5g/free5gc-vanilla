package util

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/enable-intelligent-containerized-5g/openapi/models"
	"github.com/free5gc/nwdaf/internal/logger"
	"github.com/free5gc/nwdaf/pkg/factory"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	g_logger "gorm.io/gorm/logger"
)

func InitSqlLiteDB() (err error) {
	logger.UtilLog.Info("Initializing SQL DB")
	// Database
	dbPath := factory.NwdafConfig.Configuration.SqlLiteDB

	// Open the database
	db, err := OpenDatabase(dbPath)
	if err != nil {
		return fmt.Errorf(" Fail opening the database %s: %+v", dbPath, err)
	}

	// Verify the conecction
	errConnDB := db.Raw("SELECT 1").Error
	if errConnDB != nil {
		return fmt.Errorf(" Fail to connect to database %+v: %v", dbPath, errConnDB)
	}

	// Delete tables
	err = db.Migrator().DropTable(&models.MlModelDataTable{}, &models.NFTypeTable{}, &models.EventTable{}, &models.AccuracyTable{})
	if err != nil {
		return fmt.Errorf(" Error dropping tables: %v", err)
	}

	// Migrate the database
	err = db.AutoMigrate(&models.NFTypeTable{}, &models.MlModelDataTable{}, &models.EventTable{}, &models.AccuracyTable{})
	if err != nil {
		return fmt.Errorf(" Error during AutoMigrate: %v", err)
	}

	// Insert initial data
	err = insertData(db)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	logger.InitLog.Infof("Database %s initialized successfully", dbPath)

	return nil
}

func SetupDatabase(db *gorm.DB) error {
	// Habilitar claves foráneas
	err := db.Exec("PRAGMA foreign_keys = ON;").Error

	return err
}

func OpenDatabase(dataSourceName string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dataSourceName), &gorm.Config{
		Logger: g_logger.Default.LogMode(g_logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	// Configuración inicial
	if err := SetupDatabase(db); err != nil {
		return nil, err
	}

	return db, nil
}

func insertData(db *gorm.DB) error {
	// Insertar eventos
	events := []models.EventTable{
		{
			Event: models.EventId_LOAD_LEVEL_INFORMATION, // Event Id value
		},
		{
			Event: models.EventId_NETWORK_PERFORMANCE,
		},
		{
			Event: models.EventId_NF_LOAD,
		},
		{
			Event: models.EventId_SERVICE_EXPERIENCE,
		},
		{
			Event: models.EventId_UE_MOBILITY,
		},
		{
			Event: models.EventId_UE_COMMUNICATION,
		},
		{
			Event: models.EventId_QOS_SUSTAINABILITY,
		},
		{
			Event: models.EventId_ABNORMAL_BEHAVIOUR,
		},
		{
			Event: models.EventId_USER_DATA_CONGESTION,
		},
		{
			Event: models.EventId_NSI_LOAD_LEVEL,
		},
		{
			Event: models.EventId_SM_CONGESTION,
		},
		{
			Event: models.EventId_DISPERSION,
		},
		{
			Event: models.EventId_RED_TRANS_EXP,
		},
		{
			Event: models.EventId_WLAN_PERFORMANCE,
		},
		{
			Event: models.EventId_DN_PERFORMANCE,
		},
	}
	errEvents := db.Create(&events).Error
	if errEvents != nil {
		return errEvents
	}

	// Insertar precisiones
	accuracies := []models.AccuracyTable{
		{
			Accuracy: models.NwdafMlModelAccuracy_LOW, // Accuracy value
		},
		{
			Accuracy: models.NwdafMlModelAccuracy_MEDIUM,
		},
		{
			Accuracy: models.NwdafMlModelAccuracy_HIGH,
		},
	}
	errAccuracies := db.Create(&accuracies).Error
	if errAccuracies != nil {
		return errAccuracies
	}

	// Insertar tipos de NF
	nfTypes := []models.NFTypeTable{
		{
			NfType: models.NfType_NRF, // NfType value
		},
		{
			NfType: models.NfType_UDM,
		},
		{
			NfType: models.NfType_AMF,
		},
		{
			NfType: models.NfType_SMF,
		},
		{
			NfType: models.NfType_AUSF,
		},
		{
			NfType: models.NfType_NEF,
		},
		{
			NfType: models.NfType_PCF,
		},
		{
			NfType: models.NfType_SMSF,
		},
		{
			NfType: models.NfType_NSSF,
		},
		{
			NfType: models.NfType_UDR,
		},
		{
			NfType: models.NfType_LMF,
		},
		{
			NfType: models.NfType_GMLC,
		},
		{
			NfType: models.NfType__5_G_EIR,
		},
		{
			NfType: models.NfType_SEPP,
		},
		{
			NfType: models.NfType_UPF,
		},
		{
			NfType: models.NfType_N3_IWF,
		},
		{
			NfType: models.NfType_AF,
		},
		{
			NfType: models.NfType_UDSF,
		},
		{
			NfType: models.NfType_BSF,
		},
		{
			NfType: models.NfType_CHF,
		},
		{
			NfType: models.NfType_NWDAF,
		},
	}
	errNfTypes := db.Create(&nfTypes).Error
	if errNfTypes != nil {
		logger.UtilLog.Errorf("Error to insert NfTypes: %v", errNfTypes)
		return errNfTypes
	}

	// Crear modelos ML utilizando las claves foráneas
	// mlModels := []models.MlModelDataTable{
	// 	{
	// 		Name:         "RFR_2024-11-16_16-09-17_820566",                // Name
	// 		FigureURI:    "models/figures/RFR_2024-11-16_16-09-17_820566", // figure URI
	// 		URI:          "models/RFR_2024-11-16_16-09-17_820566.pkl",     // URI
	// 		Size:         1024,                                            // in Bytes
	// 		TargetPeriod: 60,                                              // In seconds
	// 		R2:           -0.4728247842753558,                             // R2
	// 		MSE:          0.0012533946616273815,                           // MSE
	// 		R2Cpu:        -0.48400333562381537,                            // R2 cpu
	// 		R2Mem:        -0.4616462329268962,                             // R2 memory
	// 		MSECpu:       0.001124055571933576,                            // MSE cpu
	// 		MSEMem:       0.001382733751321187,                            // MSE memory
	// 		AccuracyID:   accuracies[0].ID,                                // ID from Accuracy
	// 		EventID:      events[0].ID,                                    // ID from EVent
	// 		NfTypeID:     nfTypes[0].ID,                                   // ID from NfType
	// 		Accuracy:     accuracies[0],                                   // Accuracy
	// 		Event:        events[0],                                       // Event
	// 		NfType:       nfTypes[0],                                      // NfType
	// 	},
	// 	{
	// 		Name:         "RFR_2024-11-16_16-09-17_820566",
	// 		FigureURI:    "models/figures/RFR_2024-11-16_16-09-17_820566",
	// 		URI:          "models/RFR_2024-11-16_16-09-17_820566.pkl",
	// 		Size:         3096,
	// 		TargetPeriod: 60,
	// 		R2:           -0.4728247842753558,
	// 		MSE:          0.0012533946616273815,
	// 		R2Cpu:        -0.48400333562381537,
	// 		R2Mem:        -0.4616462329268962,
	// 		MSECpu:       0.001124055571933576,
	// 		MSEMem:       0.001382733751321187,
	// 		AccuracyID:   accuracies[1].ID,
	// 		EventID:      events[2].ID,
	// 		NfTypeID:     nfTypes[2].ID,
	// 		Accuracy:     accuracies[1],
	// 		Event:        events[2],
	// 		NfType:       nfTypes[2],
	// 	},
	// 	{
	// 		Name:         "RFR_2024-11-16_16-09-17_820566",
	// 		FigureURI:    "models/figures/RFR_2024-11-16_16-09-17_820566",
	// 		URI:          "models/RFR_2024-11-16_16-09-17_820566.pkl",
	// 		Size:         2048,
	// 		TargetPeriod: 60,
	// 		R2:           -0.4728247842753558,
	// 		MSE:          0.0012533946616273815,
	// 		R2Cpu:        -0.48400333562381537,
	// 		R2Mem:        -0.4616462329268962,
	// 		MSECpu:       0.001124055571933576,
	// 		MSEMem:       0.001382733751321187,
	// 		AccuracyID:   accuracies[1].ID,
	// 		EventID:      events[2].ID,
	// 		NfTypeID:     nfTypes[2].ID,
	// 		Accuracy:     accuracies[1],
	// 		Event:        events[2],
	// 		NfType:       nfTypes[2],
	// 	},
	// 	{
	// 		Name:         "RFR_2024-11-16_16-09-17_820566",
	// 		FigureURI:    "models/figures/RFR_2024-11-16_16-09-17_820566",
	// 		URI:          "models/RFR_2024-11-16_16-09-17_820566.pkl",
	// 		Size:         2048,
	// 		TargetPeriod: 60,
	// 		R2:           -0.4728247842753558,
	// 		MSE:          0.0012533946616273815,
	// 		R2Cpu:        -0.48400333562381537,
	// 		R2Mem:        -0.4616462329268962,
	// 		MSECpu:       0.001124055571933576,
	// 		MSEMem:       0.001382733751321187,
	// 		AccuracyID:   accuracies[1].ID,
	// 		EventID:      events[2].ID,
	// 		NfTypeID:     nfTypes[2].ID,
	// 		Accuracy:     accuracies[1],
	// 		Event:        events[2],
	// 		NfType:       nfTypes[2],
	// 	},
	// 	{
	// 		Name:         "RFR_2024-11-16_16-09-17_820566",
	// 		FigureURI:    "models/figures/RFR_2024-11-16_16-09-17_820566",
	// 		URI:          "models/RFR_2024-11-16_16-09-17_820566.pkl",
	// 		Size:         2048,
	// 		TargetPeriod: 60,
	// 		R2:           -0.4728247842753558,
	// 		MSE:          0.0012533946616273815,
	// 		R2Cpu:        -0.48400333562381537,
	// 		R2Mem:        -0.4616462329268962,
	// 		MSECpu:       0.001124055571933576,
	// 		MSEMem:       0.001382733751321187,
	// 		AccuracyID:   accuracies[2].ID,
	// 		EventID:      events[2].ID,
	// 		NfTypeID:     nfTypes[2].ID,
	// 		Accuracy:     accuracies[2],
	// 		Event:        events[2],
	// 		NfType:       nfTypes[2],
	// 	},
	// 	{
	// 		Name:         "RFR_2024-11-16_16-09-17_820566",
	// 		FigureURI:    "models/figures/RFR_2024-11-16_16-09-17_820566",
	// 		URI:          "models/RFR_2024-11-16_16-09-17_820566.pkl",
	// 		Size:         2048,
	// 		TargetPeriod: 60,
	// 		R2:           -0.4728247842753558,
	// 		MSE:          0.0012533946616273815,
	// 		R2Cpu:        -0.48400333562381537,
	// 		R2Mem:        -0.4616462329268962,
	// 		MSECpu:       0.001124055571933576,
	// 		MSEMem:       0.001382733751321187,
	// 		AccuracyID:   accuracies[0].ID,
	// 		EventID:      events[2].ID,
	// 		NfTypeID:     nfTypes[2].ID,
	// 		Accuracy:     accuracies[0],
	// 		Event:        events[2],
	// 		NfType:       nfTypes[2],
	// 	},
	// 	{
	// 		Name:         "RFR_2024-11-16_16-09-17_820566",
	// 		FigureURI:    "models/figures/RFR_2024-11-16_16-09-17_820566",
	// 		URI:          "models/RFR_2024-11-16_16-09-17_820566.pkl",
	// 		Size:         2048,
	// 		TargetPeriod: 60,
	// 		R2:           -0.4728247842753558,
	// 		MSE:          0.0012533946616273815,
	// 		R2Cpu:        -0.48400333562381537,
	// 		R2Mem:        -0.4616462329268962,
	// 		MSECpu:       0.001124055571933576,
	// 		MSEMem:       0.001382733751321187,
	// 		AccuracyID:   accuracies[0].ID,
	// 		EventID:      events[2].ID,
	// 		NfTypeID:     nfTypes[2].ID,
	// 		Accuracy:     accuracies[0],
	// 		Event:        events[2],
	// 		NfType:       nfTypes[2],
	// 	},
	// 	{
	// 		Name:         "RFR_2024-11-16_16-09-17_820566",
	// 		FigureURI:    "models/figures/RFR_2024-11-16_16-09-17_820566",
	// 		URI:          "models/RFR_2024-11-16_16-09-17_820566.pkl",
	// 		Size:         3048,
	// 		TargetPeriod: 60,
	// 		R2:           -0.4728247842753558,
	// 		MSE:          0.0012533946616273815,
	// 		R2Cpu:        -0.48400333562381537,
	// 		R2Mem:        -0.4616462329268962,
	// 		MSECpu:       0.001124055571933576,
	// 		MSEMem:       0.001382733751321187,
	// 		AccuracyID:   accuracies[1].ID,
	// 		EventID:      events[2].ID,
	// 		NfTypeID:     nfTypes[2].ID,
	// 		Accuracy:     accuracies[1],
	// 		Event:        events[2],
	// 		NfType:       nfTypes[2],
	// 	},
	// 	{
	// 		Name:         "RFR_2024-11-16_16-09-17_820566",
	// 		FigureURI:    "models/figures/RFR_2024-11-16_16-09-17_820566",
	// 		URI:          "models/RFR_2024-11-16_16-09-17_820566.pkl",
	// 		Size:         3048,
	// 		TargetPeriod: 60,
	// 		R2:           -0.4728247842753558,
	// 		MSE:          0.0012533946616273815,
	// 		R2Cpu:        -0.48400333562381537,
	// 		R2Mem:        -0.4616462329268962,
	// 		MSECpu:       0.001124055571933576,
	// 		MSEMem:       0.001382733751321187,
	// 		AccuracyID:   accuracies[1].ID,
	// 		EventID:      events[2].ID,
	// 		NfTypeID:     nfTypes[2].ID,
	// 		Accuracy:     accuracies[1],
	// 		Event:        events[2],
	// 		NfType:       nfTypes[2],
	// 	},
	// }

	modelsFilePath := NwdafDefaultModelsPath + NwdafDefaultInitModelInfoListFile
	mlModels, errGetModels := loadDefaultMlModels(modelsFilePath)
	if errGetModels != nil {
		logger.UtilLog.Warnf("Failed to load default model list: %s", errGetModels.Error())
	} else {
		var modelsErrorList []string
		var modelsSavedList []string
		var errSaveList []string
		for index, model := range mlModels {
			modelNum := index + 1
			modelConfidence := models.MlModelDataConfidence{
				R2:           model.R2,
				MSE:          model.MSE,
				R2Cpu:        model.R2,
				R2Mem:        model.R2Mem,
				R2Troughput:  model.R2Thrpt,
				MSECpu:       model.MSECpu,
				MSEMem:       model.MSEMem,
				MSETroughput: model.MSEThrpt,
			}

			mlModelInfo := models.MlModelData{
				EventId:      model.Event.Event,
				Name:         model.Name,
				Size:         model.Size,
				FigureURI:    model.FigureURI,
				TargetPeriod: model.TargetPeriod,
				Confidence:   modelConfidence,
				URI:          model.URI,
				Accuracy:     model.Accuracy.Accuracy,
				NfType:       model.NfType.NfType,
			}

			mlModelSaveResponse, saved, errSave := SaveMlModelProcedure(mlModelInfo)

			if errSave != nil {
				errSaveList = append(errSaveList, fmt.Sprintf("Model #%d: %s", modelNum, errSave.Cause))
			}

			if !saved {
				modelsErrorList = append(modelsErrorList, fmt.Sprintf("Model #%d: %s", modelNum, mlModelInfo.Name))
			} else {
				modelsSavedList = append(modelsSavedList, fmt.Sprintf("Model #%d: %s", modelNum, mlModelSaveResponse.MlModels[0].Name))
			}
		}

		logger.UtilLog.Infof("Saved models: %v", modelsSavedList)

		if len(errSaveList) > 0 {
			logger.UtilLog.Warnf("Some models were not saved %v: %v", modelsErrorList, errSaveList)
		}

		// errModels := db.Create(&mlModels).Error
		// if errModels != nil {
		// 	return errModels
		// }
	}

	return nil
}

func loadDefaultMlModels(filePath string) ([]models.MlModelDataTable, error) {
	// Verify if the file exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("error: no found the file '%s'", filePath)
	}

	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading the file: %v", err)
	}

	// Parse the data
	var models []models.MlModelDataTable
	err = json.Unmarshal(data, &models)
	if err != nil {
		return nil, fmt.Errorf("error decoing the JSON file: %v", err)
	}

	return models, nil
}
