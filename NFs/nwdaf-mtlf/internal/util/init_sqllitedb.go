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

// Definir los modelos

// func (EventTable) TableName() string {
//     return models.
// }

// func (AccuracyTable) TableName() string {
// 	return string(models.NwdafMLModelDB_ACCURACY_TABLE_NAME) // Si deseas un nombre específico para la tabla
// }

// func (EventTable) TableName() string {
// 	return string(models.NwdafMLModelDB_EVENT_ID_TABLE_NAME) // Si deseas un nombre específico para la tabla
// }

// func (NFTypeTable) TableName() string {
// 	return string(models.NwdafMLModelDB_NF_TYPE_TABLE_NAME) // Si deseas un nombre específico para la tabla
// }

// func (MlModelDataTable) TableName() string {
// 	return string(models.NwdafMLModelDB_ML_MODEL_INFO_TABLE_NAME) // Si deseas un nombre específico para la tabla
// }

type EventTable struct {
	ID      int64 `gorm:"primaryKey"`
	Event   string
	MlModel []MlModelDataTable `gorm:"foreignKey:EventID"`
}

type AccuracyTable struct {
	ID       int64 `gorm:"primaryKey"`
	Accuracy string
	MlModel  []MlModelDataTable `gorm:"foreignKey:AccuracyID"`
}

type NFTypeTable struct {
	ID      int64 `gorm:"primaryKey"`
	NfType  string
	MlModel []MlModelDataTable `gorm:"foreignKey:NfTypeID"`
}

type MlModelDataTable struct {
	ID           int64 `gorm:"primaryKey"`
	URI          string
	Size         int64
	TargetPeriod int64
	NfTypeID     int64 `gorm:"foreignKey:ID"`
	AccuracyID   int64 `gorm:"foreignKey:ID"`
	EventID      int64 `gorm:"foreignKey:ID"`
	NfType       NFTypeTable
	Accuracy     AccuracyTable
	Event        EventTable
	// NFType       NFTypeTable `gorm:"foreignKey:NfTypeID"`
	// Accuracy     AccuracyTable `gorm:"foreignKey:NfTypeID"`
	// Event        EventTable `gorm:"foreignKey:NfTypeID"`
	// 	NfTypeID   models.NfType               `gorm:"not null"`
	// 	AccuracyID models.NwdafMlModelAccuracy `gorm:"not null"`
	// 	EventID    models.EventId              `gorm:"not null"`
}

func InitSqlLiteDB() (err error) {
	// Database
	sqldb := factory.NwdafConfig.Configuration.SqlLiteDB
	// Table names
	// mlInfoTable := string(models.NwdafMLModelDB_ML_MODEL_INFO)
	// accuracyTable := string(models.NwdafMLModelDB_ACCURACY_VALUES)
	// nfTable := string(models.NwdafMLModelDB_NF_TYPES)
	// eventTable := string(models.NwdafMLModelDB_EVENT_ID)

	db, err := OpenDatabase(sqldb)
	// db, err := sql.Open("sqlite3", sqldb)
	if err != nil {
		logger.UtilLog.Errorf("Error opening the %s database  ", sqldb)
		return err
	}

	err = db.Migrator().DropTable(&MlModelDataTable{}, &NFTypeTable{}, &EventTable{}, &AccuracyTable{})
	if err != nil {
		logger.UtilLog.Errorf("Error al eliminar las tablas: %s", err)
	}

	logger.UtilLog.Info("¡Tablas eliminadas exitosamente!")

	// Migrar las estructuras a la base de datos
	err = db.AutoMigrate(&NFTypeTable{}, &MlModelDataTable{}, &EventTable{}, &AccuracyTable{})
	if err != nil {
		logger.UtilLog.Errorf("Error al migrar la base de datos: %s", err)
		return err
	}

	logger.UtilLog.Info("¡Migración exitosa!")

	err = insertData(db)
	if err != nil {
		logger.UtilLog.Errorf("Error al insertar datos: %s", err)
		return err
	}

	logger.UtilLog.Info("¡Inicio de DB exitoso!")

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
			// ID:      string(models.EventId_LOAD_LEVEL_INFORMATION),
			Event:   string(models.EventId_LOAD_LEVEL_INFORMATION),
			MlModel: []MlModelDataTable{},
		},
		// {
		// 	ID:      models.EventId_NETWORK_PERFORMANCE,
		// 	Event:   models.EventId_NETWORK_PERFORMANCE,
		// 	MlModel: []MlModelDataTable{},
		// },
		// {
		// 	ID:      models.EventId_NF_LOAD,
		// 	Event:   models.EventId_NF_LOAD,
		// 	MlModel: []MlModelDataTable{},
		// },
		// {
		// 	ID:      models.EventId_SERVICE_EXPERIENCE,
		// 	Event:   models.EventId_SERVICE_EXPERIENCE,
		// 	MlModel: []MlModelDataTable{},
		// },
		// {
		// 	ID:      models.EventId_UE_MOBILITY,
		// 	Event:   models.EventId_UE_MOBILITY,
		// 	MlModel: []MlModelDataTable{},
		// },
		// {
		// 	ID:      models.EventId_UE_COMMUNICATION,
		// 	Event:   models.EventId_UE_COMMUNICATION,
		// 	MlModel: []MlModelDataTable{},
		// },
		// {
		// 	ID:      models.EventId_QOS_SUSTAINABILITY,
		// 	Event:   models.EventId_QOS_SUSTAINABILITY,
		// 	MlModel: []MlModelDataTable{},
		// },
		// {
		// 	ID:      models.EventId_ABNORMAL_BEHAVIOUR,
		// 	Event:   models.EventId_ABNORMAL_BEHAVIOUR,
		// 	MlModel: []MlModelDataTable{},
		// },
		// {
		// 	ID:      models.EventId_USER_DATA_CONGESTION,
		// 	Event:   models.EventId_USER_DATA_CONGESTION,
		// 	MlModel: []MlModelDataTable{},
		// },
		// {
		// 	ID:      models.EventId_NSI_LOAD_LEVEL,
		// 	Event:   models.EventId_NSI_LOAD_LEVEL,
		// 	MlModel: []MlModelDataTable{},
		// },
		// {
		// 	ID:      models.EventId_SM_CONGESTION,
		// 	Event:   models.EventId_SM_CONGESTION,
		// 	MlModel: []MlModelDataTable{},
		// },
		// {
		// 	ID:      models.EventId_DISPERSION,
		// 	Event:   models.EventId_DISPERSION,
		// 	MlModel: []MlModelDataTable{},
		// },
		// {
		// 	ID:      models.EventId_RED_TRANS_EXP,
		// 	Event:   models.EventId_RED_TRANS_EXP,
		// 	MlModel: []MlModelDataTable{},
		// },
		// {
		// 	ID:      models.EventId_WLAN_PERFORMANCE,
		// 	Event:   models.EventId_WLAN_PERFORMANCE,
		// 	MlModel: []MlModelDataTable{},
		// },
		// {
		// 	ID:      models.EventId_DN_PERFORMANCE,
		// 	Event:   models.EventId_DN_PERFORMANCE,
		// 	MlModel: []MlModelDataTable{},
		// },
	}
	errEvents := db.Create(&events).Error
	if errEvents != nil {
		logger.UtilLog.Errorf("Error al insertar Events: %v", errEvents)
		return errEvents
	}

	// Insertar precisiones
	accuracies := []AccuracyTable{
		{
			// ID:       models.NwdafMlModelAccuracy_LOW,
			Accuracy: string(models.NwdafMlModelAccuracy_LOW),
			MlModel:  []MlModelDataTable{},
		},
		{
			// ID:       models.NwdafMlModelAccuracy_MEDIUM,
			Accuracy: string(models.NwdafMlModelAccuracy_MEDIUM),
			MlModel:  []MlModelDataTable{},
		},
		// {
		// 	ID:       models.NwdafMlModelAccuracy_HIGH,
		// 	Accuracy: models.NwdafMlModelAccuracy_HIGH,
		// 	// MlModel:  []MlModelDataTable{},
		// },
	}
	errAccuracy := db.Create(&accuracies).Error
	if errAccuracy != nil {
		logger.UtilLog.Errorf("Error al insertar Accuracies: %v", errAccuracy)
		return errAccuracy
	}

	// Insertar tipos de NF
	nfTypes := []NFTypeTable{
		{
			// ID:      models.NfType_NRF,
			NfType:  string(models.NfType_NRF),
			MlModel: []MlModelDataTable{},
		},
		{
			// ID:      models.NfType_UDM,
			NfType:  string(models.NfType_UDM),
			MlModel: []MlModelDataTable{},
		},
		// {
		// 	ID:      models.NfType_AMF,
		// 	NfType:  models.NfType_AMF,
		// 	MlModel: []MlModelDataTable{},
		// },
		// {
		// 	ID:      models.NfType_SMF,
		// 	NfType:  models.NfType_SMF,
		// 	MlModel: []MlModelDataTable{},
		// },
		// {
		// 	ID:      models.NfType_AUSF,
		// 	NfType:  models.NfType_AUSF,
		// 	MlModel: []MlModelDataTable{},
		// },
		// {
		// 	ID:      models.NfType_NEF,
		// 	NfType:  models.NfType_NEF,
		// 	MlModel: []MlModelDataTable{},
		// },
		// {
		// 	ID:      models.NfType_PCF,
		// 	NfType:  models.NfType_PCF,
		// 	MlModel: []MlModelDataTable{},
		// },
		// {
		// 	ID:      models.NfType_SMSF,
		// 	NfType:  models.NfType_SMSF,
		// 	MlModel: []MlModelDataTable{},
		// },
		// {
		// 	ID:      models.NfType_NSSF,
		// 	NfType:  models.NfType_NSSF,
		// 	MlModel: []MlModelDataTable{},
		// },
		// {
		// 	ID:      models.NfType_UDR,
		// 	NfType:  models.NfType_UDR,
		// 	MlModel: []MlModelDataTable{},
		// },
		// {
		// 	ID:      models.NfType_LMF,
		// 	NfType:  models.NfType_LMF,
		// 	MlModel: []MlModelDataTable{},
		// },
		// {
		// 	ID:      models.NfType_GMLC,
		// 	NfType:  models.NfType_GMLC,
		// 	MlModel: []MlModelDataTable{},
		// },
		// {
		// 	ID:      models.NfType__5_G_EIR,
		// 	NfType:  models.NfType__5_G_EIR,
		// 	MlModel: []MlModelDataTable{},
		// },
		// {
		// 	ID:      models.NfType_SEPP,
		// 	NfType:  models.NfType_SEPP,
		// 	MlModel: []MlModelDataTable{},
		// },
		// {
		// 	ID:      models.NfType_UPF,
		// 	NfType:  models.NfType_UPF,
		// 	MlModel: []MlModelDataTable{},
		// },
		// {
		// 	ID:      models.NfType_N3_IWF,
		// 	NfType:  models.NfType_N3_IWF,
		// 	MlModel: []MlModelDataTable{},
		// },
		// {
		// 	ID:      models.NfType_AF,
		// 	NfType:  models.NfType_AF,
		// 	MlModel: []MlModelDataTable{},
		// },
		// {
		// 	ID:      models.NfType_UDSF,
		// 	NfType:  models.NfType_UDSF,
		// 	MlModel: []MlModelDataTable{},
		// },
		// {
		// 	ID:      models.NfType_BSF,
		// 	NfType:  models.NfType_BSF,
		// 	MlModel: []MlModelDataTable{},
		// },
		// {
		// 	ID:      models.NfType_CHF,
		// 	NfType:  models.NfType_CHF,
		// 	MlModel: []MlModelDataTable{},
		// },
		// {
		// 	ID:      models.NfType_NWDAF,
		// 	NfType:  models.NfType_NWDAF,
		// 	MlModel: []MlModelDataTable{},
		// },
	}
	errNfTypes := db.Create(&nfTypes).Error
	if errNfTypes != nil {
		logger.UtilLog.Errorf("Error al insertar NfTypes: %v", errNfTypes)
		return errNfTypes
	}

	// Crear modelos ML utilizando las claves foráneas
	mlModels := []MlModelDataTable{
		{
			URI:          "http://example.com/model1",
			Size:         1024,
			TargetPeriod: 60,
			AccuracyID:   accuracies[0].ID,
			EventID:      events[0].ID,  // Usar ID del evento
			NfTypeID:     nfTypes[0].ID, // Usar ID del tipo de NF
			Accuracy:     accuracies[0], // Usar ID de la precisión
			Event:        events[0],     // Usar ID del evento
			NfType:       nfTypes[0],    // Usar ID del tipo de NF
		},
		// {
		// 	URI:          "http://example.com/model2",
		// 	Size:         2048,
		// 	TargetPeriod: 120,
		// 	NfTypeID:     models.NfType_AMF,                  // Usar ID del tipo de NF
		// 	AccuracyID:   models.NwdafMlModelAccuracy_MEDIUM, // Usar ID de la precisión
		// 	EventID:      models.EventId_NF_LOAD,             // Usar ID del evento
		// },
		// {
		// 	URI:          "http://example.com/model3",
		// 	Size:         3096,
		// 	TargetPeriod: 180,
		// 	NfTypeID:     models.NfType_AMF,                  // Usar ID del tipo de NF
		// 	AccuracyID:   models.NwdafMlModelAccuracy_MEDIUM, // Usar ID de la precisión
		// 	EventID:      models.EventId_NF_LOAD,             // Usar ID del evento
		// },
	}
	errModels := db.Create(&mlModels).Error
	if errModels != nil {
		logger.UtilLog.Errorf("Error al insertar MlModels: %v", errModels)
		return errModels
	}

	return nil
}

// func printData(db *gorm.DB) {
// 	// Modelo que se va a consultar
// 	var nfTypes []NFTypeTable

// 	// Seleccionar todos los registros de la tabla
// 	result := db.Preload("MlModel").Find(&nfTypes)
// 	if result.Error != nil {
// 		logger.UtilLog.Errorf("Error al consultar la tabla: %v", result.Error)
// 	}

// 	// Mostrar los resultados en los logs
// 	logger.UtilLog.Info("Registros obtenidos:")
// 	for _, nfType := range nfTypes {
// 		logger.UtilLog.Warnf("%+v\n", nfType)
// 	}
// }
