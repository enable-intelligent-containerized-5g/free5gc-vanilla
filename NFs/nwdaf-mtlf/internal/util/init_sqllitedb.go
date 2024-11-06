package util

import (
	// "github.com/enable-intelligent-containerized-5g/openapi"
	"github.com/enable-intelligent-containerized-5g/openapi/models"
	"github.com/free5gc/nwdaf/internal/logger"
	"github.com/free5gc/nwdaf/pkg/factory"

	// "database/sql"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type EventTable struct {
	ID    int64 `gorm:"primaryKey"`
	Event models.EventId
}

type AccuracyTable struct {
	ID       int64 `gorm:"primaryKey"`
	Accuracy models.NwdafMlModelAccuracy
}

type NFTypeTable struct {
	ID     int64 `gorm:"primaryKey"`
	NfType models.NfType
}

type MlModelDataTable struct {
	ID           int64 `gorm:"primaryKey"`
	URI          string
	Size         int64
	TargetPeriod int64
	NfTypeID     int64         `gorm:"foreignKey:ID"`
	AccuracyID   int64         `gorm:"foreignKey:ID"`
	EventID      int64         `gorm:"foreignKey:ID"`
	NfType       NFTypeTable   // property name in: models.NwdafMLModelDB_NF_TYPE_KEY
	Accuracy     AccuracyTable // property name in: models.NwdafMLModelDB_ACCURACY_KEY
	Event        EventTable    // property name in: models.NwdafMLModelDB_EVENT_ID_KEY
}

func InitSqlLiteDB() (err error) {
	// Database
	sqldb := factory.NwdafConfig.Configuration.SqlLiteDB

	db, err := OpenDatabase(sqldb)
	if err != nil {
		return err
	}

	err = db.Migrator().DropTable(&MlModelDataTable{}, &NFTypeTable{}, &EventTable{}, &AccuracyTable{})
	if err != nil {
		return err
	}

	// Migrate the database
	err = db.AutoMigrate(&NFTypeTable{}, &MlModelDataTable{}, &EventTable{}, &AccuracyTable{})
	if err != nil {
		return err
	}

	err = insertData(db)
	if err != nil {
		return err
	}

	logger.UtilLog.Infof("The database %s was started successfully", sqldb)

	return nil
}

func SetupDatabase(db *gorm.DB) error {
	// Habilitar claves foráneas
	err := db.Exec("PRAGMA foreign_keys = ON;").Error

	return err
}

func OpenDatabase(dataSourceName string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dataSourceName), &gorm.Config{})
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
	events := []EventTable{
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
	accuracies := []AccuracyTable{
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
	nfTypes := []NFTypeTable{
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
		logger.UtilLog.Errorf("Error al insertar NfTypes: %v", errNfTypes)
		return errNfTypes
	}

	// Crear modelos ML utilizando las claves foráneas
	mlModels := []MlModelDataTable{
		{
			URI:          "http://example.com/model1", // URI
			Size:         1024,                        // in Bytes
			TargetPeriod: 60,                          // In seconds
			AccuracyID:   accuracies[0].ID,            // ID from Accuracy
			EventID:      events[0].ID,                // ID from EVent
			NfTypeID:     nfTypes[0].ID,               // ID from NfType
			Accuracy:     accuracies[0],               // Accuracy
			Event:        events[0],                   // Event
			NfType:       nfTypes[0],                  // NfType
		},
		{
			URI:          "http://example.com/model2",
			Size:         3096,
			TargetPeriod: 120,
			AccuracyID:   accuracies[1].ID,
			EventID:      events[1].ID,
			NfTypeID:     nfTypes[1].ID,
			Accuracy:     accuracies[1],
			Event:        events[1],
			NfType:       nfTypes[1],
		},
		{
			URI:          "http://example.com/model3",
			Size:         2048,
			TargetPeriod: 120,
			AccuracyID:   accuracies[2].ID,
			EventID:      events[2].ID,
			NfTypeID:     nfTypes[2].ID,
			Accuracy:     accuracies[2],
			Event:        events[2],
			NfType:       nfTypes[2],
		},
	}
	errModels := db.Create(&mlModels).Error
	if errModels != nil {
		return errModels
	}

	return nil
}
