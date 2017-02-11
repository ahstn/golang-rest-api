package daos

import (
	"database/sql"

	"github.com/phazyy/golang-rest-api/models"
	"gopkg.in/inconshreveable/log15.v2"
)

func GetAll(db *sql.DB, log log15.Logger) models.PilotSlice {
	pilots, err := models.Pilots(db).All()
	if err != nil {
		log.Error("db: failed to get pilots", "err", err)
		return nil
	}

	log.Info("db: fetched pilots", "count", len(pilots))
	return pilots
}

func Create(name string, db *sql.DB, log log15.Logger) int {
	pilot := &models.Pilot{
		Name: name,
	}
	if err := pilot.Insert(db); err != nil {
		log.Error("db: failed to insert pilot", "err", err)
		return 0
	}

	log.Info("db: insterted pilot", "id", pilot.ID)
	return pilot.ID
}
