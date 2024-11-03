package util

import (
	// "github.com/enable-intelligent-containerized-5g/openapi"
	"github.com/enable-intelligent-containerized-5g/openapi/models"
	"github.com/free5gc/nwdaf/internal/logger"
	"github.com/free5gc/nwdaf/pkg/factory"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type NfRecord struct {
	ID     string `json:"id"`      // ID del registro
	NfType string `json:"nf_type"` // Tipo de NF (Network Function)
}

func InitSqlLiteDB() (err error) {
	sqldb := factory.NwdafConfig.Configuration.SqlLiteDB
	// Nombres de las tablas
	// mlInfoTable := factory.NwdafConfig.Configuration.SqlLiteTableName
	mlInfoTable := string(models.NwdafMLModelDB_ML_MODEL_INFO)
	accuracyTable := string(models.NwdafMLModelDB_ACCURACY_VALUES)
	nfTable := string(models.NwdafMLModelDB_NF_TYPES)
	eventTable := string(models.NwdafMLModelDB_EVENT_ID)
	db, err := sql.Open("sqlite3", sqldb)
	if err != nil {
		logger.UtilLog.Errorf("Error opening the %s database  ", sqldb)
		return err
	}
	defer db.Close()

	// Eliminar tablas si existen
	_, err = db.Exec("DROP TABLE IF EXISTS " + eventTable + ";")

	_, err = db.Exec("DROP TABLE IF EXISTS " + mlInfoTable + ";")

	_, err = db.Exec("DROP TABLE IF EXISTS " + accuracyTable + ";")

	_, err = db.Exec("DROP TABLE IF EXISTS " + nfTable + ";")

	// Crear tabla event_id
	_, err = db.Exec(`
		CREATE TABLE ` + eventTable + ` (
			id TEXT PRIMARY KEY NOT NULL UNIQUE,
			event_id TEXT NOT NULL UNIQUE
		);
	`)

	// Crear tabla nf_types
	_, err = db.Exec(`
		CREATE TABLE ` + nfTable + ` (
			id TEXT PRIMARY KEY NOT NULL UNIQUE,
			nf_type TEXT NOT NULL UNIQUE
		);
	`)

	// Crear tabla accuracy_values
	_, err = db.Exec(`
		CREATE TABLE ` + accuracyTable + ` (
			id TEXT PRIMARY KEY NOT NULL UNIQUE,
			accuracy TEXT NOT NULL UNIQUE
		);
	`)

	// Crear tabla records
	_, err = db.Exec(`
		CREATE TABLE ` + mlInfoTable + ` (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			uri TEXT NOT NULL,
			accuracy TEXT NOT NULL,
			nf_type TEXT NOT NULL,
			event_id TEXT NOT NULL,
			target_period TEXT NOT NULL,
			FOREIGN KEY (accuracy) REFERENCES ` + accuracyTable + ` (id),
			FOREIGN KEY (nf_type) REFERENCES ` + nfTable + ` (id),
			FOREIGN KEY (event_id) REFERENCES ` + eventTable + ` (id)
		);
	`)

	// Habilitar claves foráneas
    _, err = db.Exec("PRAGMA foreign_keys = ON;")

	// Insertar un registro en tablas
	_, err = db.Exec(`
	INSERT INTO `+eventTable+` (id, event_id) 
	VALUES (?, ?);`, "NF_LOAD", "NF_LOAD")


	_, err = db.Exec(
	`INSERT INTO `+accuracyTable+` (id, accuracy) 
	VALUES (?, ?);`, "LOW", "LOW")

	_, err = db.Exec(
		`INSERT INTO `+nfTable+` (id, nf_type) 
	VALUES (?, ?);`,
		"AMF", "AMF")

	_, err = db.Exec(`
    INSERT INTO `+mlInfoTable+` (uri, accuracy, nf_type, event_id, target_period) 
    VALUES (?, ?, ?, ?, ?);`,
		"http://example.com/model1", "LOW", "AMF", "NF_LOAD", "60") // Valores de ejemplo

	_, err = db.Exec(`
    INSERT INTO `+mlInfoTable+` (uri, accuracy, nf_type, event_id, target_period) 
    VALUES (?, ?, ?, ?, ?);`,
	"http://example.com/model2", "LOW", "AMF", "NF_LOAD", "70") // Valores de ejemplo


	// // Consultar
	// rows, err := db.Query(`
    // SELECT id, nf_type 
    // FROM ` + nfTable + `;`)

	// if err != nil {
	// 	logger.UtilLog.Error("Error al consultar los registros: ", err)
	// }

	// // Estructura para almacenar los resultados
	// var nfRecords []NfRecord

	// // Iterar sobre los resultados y mapearlos a la estructura
	// for rows.Next() {
	// 	var record NfRecord
	// 	err := rows.Scan(&record.ID, &record.NfType)
	// 	if err != nil {
	// 		logger.UtilLog.Error("Error al leer los registros: ", err)
	// 		return nil
	// 	}
	// 	nfRecords = append(nfRecords, record)
	// }

	// // Verificar si hubo algún error durante la iteración
	// if err := rows.Err(); err != nil {
	// 	logger.UtilLog.Error("Error en el procesamiento de filas: ", err)
	// 	return nil
	// }
	// logger.UtilLog.Info("Datos de "+nfTable+": ", nfRecords)

	if err != nil {
		logger.UtilLog.Errorf("Error initializing the database: %s", sqldb)
		return err
	}

	logger.UtilLog.Infof("The %s database has been initialized successfully.", sqldb)

	db.Close()

	return nil
}
