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

func Update(id int, name string, db *sql.DB, log log15.Logger) int {
	pilot, _ := models.FindPilot(db, id)
	pilot.Name = name
	if err := pilot.Update(db); err != nil {
		log.Error("db: failed to update pilot", "id", id, "err", err)
	}

	log.Info("db: updated pilot", "id", id)
	return id
}

func Delete(id int, db *sql.DB, log log15.Logger) int {
	pilot, _ := models.FindPilot(db, id)
	if err := pilot.Delete(db); err != nil {
		log.Error("db: failed to delete pilot", "id", id, "err", err)
		return 0
	}

	log.Info("db: deleted pilot", "id", id)
	return id
}
