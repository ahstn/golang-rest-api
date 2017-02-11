package daos

import (
	"database/sql"

	"github.com/phazyy/sqlboiler-gin/models"
	"gopkg.in/inconshreveable/log15.v2"
)

func GetPilots(db *sql.DB, log log15.Logger) models.PilotSlice {
	pilots, err := models.Pilots(db).All()
	if err != nil {
		log.Error("failed to get pilots", "err", err)
		return nil
	}

	log.Info("db: fetched pilots")
	return pilots
}
