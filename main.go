package main

//go:generate sqlboiler postgres

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/phazyy/sqlboiler-basic/models"
	"gopkg.in/inconshreveable/log15.v2"
)

var log = log15.New()

func main() {
	r := gin.Default()

	db, err := db_open()
	if err != nil {
		log.Error("failed to open database", "err", err)
		return
	}
	pilots, err := db_getPilots(db)
	if err != nil {
		log.Error("failed to open database", "err", err)
		return
	}

	r.GET("/pilots", func(c *gin.Context) {
		c.JSON(200, pilots)
	})

	r.Run(":8080")
}
func db_open() (*sql.DB, error) {
	db, err := sql.Open("postgres", `dbname=flight host=localhost user=boiler password=boiler sslmode=disable`)
	if err != nil {
		log.Error("failed to open database", "err", err)
		return nil, err
	}
	return db, nil
}
func db_getPilots(db *sql.DB) (models.PilotSlice, error) {
	pilots, err := models.Pilots(db).All()
	if err != nil {
		log.Error("failed to get pilots", "err", err)
		return nil, err
	}

	for i, j := range pilots {
		log.Info("found pilot", "i", i, "Name", j.Name)
	}

	return pilots, nil
}
