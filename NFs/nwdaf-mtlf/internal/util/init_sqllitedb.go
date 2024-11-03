package util

import (
	// "github.com/enable-intelligent-containerized-5g/openapi"
	"github.com/enable-intelligent-containerized-5g/openapi/models"
	"github.com/free5gc/nwdaf/internal/logger"
	"github.com/free5gc/nwdaf/pkg/factory"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func InitSqlLiteDB() error {
	sqldb := factory.NwdafConfig.Configuration.SqlLiteDB
	// Nombres de las tablas
	// mlInfoTable := factory.NwdafConfig.Configuration.SqlLiteTableName
	mlInfoTable := string(models.NwdafMLModelDB_ML_MODEL_INFO)
	accuracyTable := string(models.NwdafMLModelDB_ACCURACY_VALUES)
	nfTable := string(models.NwdafMLModelDB_NF_TYPES)
	eventTable := string(models.NwdafMLModelDB_EVENT_ID)
	db, err := sql.Open("sqlite3", sqldb)
	if err != nil {
		return err
	}
	defer db.Close()

	// Eliminar tablas si existen
	_, err = db.Exec("DROP TABLE IF EXISTS " + eventTable + ";")
	if err != nil {
		return err
	}

	_, err = db.Exec("DROP TABLE IF EXISTS " + mlInfoTable + ";")
	if err != nil {
		return err
	}

	_, err = db.Exec("DROP TABLE IF EXISTS " + accuracyTable + ";")
	if err != nil {
		return err
	}

	_, err = db.Exec("DROP TABLE IF EXISTS " + nfTable + ";")
	if err != nil {
		return err
	}

	// Crear tabla event_id
	_, err = db.Exec(`
		CREATE TABLE ` + eventTable + ` (
			id TEXT PRIMARY KEY NOT NULL UNIQUE,
			event_id TEXT NOT NULL UNIQUE
		);
	`)
	if err != nil {
		return err
	}

	// Crear tabla nf_types
	_, err = db.Exec(`
		CREATE TABLE ` + nfTable + ` (
			id TEXT PRIMARY KEY NOT NULL UNIQUE,
			nf_type TEXT NOT NULL UNIQUE
		);
	`)
	if err != nil {
		return err
	}

	// Crear tabla accuracy_values
	_, err = db.Exec(`
		CREATE TABLE ` + accuracyTable + ` (
			id TEXT PRIMARY KEY NOT NULL UNIQUE,
			accuracy_value TEXT NOT NULL UNIQUE
		);
	`)
	if err != nil {
		return err
	}

	// Crear tabla records
	_, err = db.Exec(`
		CREATE TABLE ` + mlInfoTable + ` (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			uri TEXT NOT NULL,
			accuracy TEXT NOT NULL,
			nf TEXT NOT NULL,
			event TEXT NOT NULL,
			target_period TEXT NOT NULL,
			FOREIGN KEY (accuracy) REFERENCES ` + accuracyTable + ` (accuracy_value),
			FOREIGN KEY (nf) REFERENCES ` + nfTable + ` (nf_type),
			FOREIGN KEY (event) REFERENCES ` + eventTable + ` (event_id)
		);
	`)
	if err != nil {
		return err
	}

	logger.InitLog.Infof("Tables for SQLite DB created")
	return nil
}
