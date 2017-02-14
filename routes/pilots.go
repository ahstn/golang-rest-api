package routes

import (
	"database/sql"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phazyy/golang-rest-api/models"
	"gopkg.in/inconshreveable/log15.v2"
)

// PilotRoutes :
type PilotRoutes struct{}

// Get : Attempts to fetch a single pilot matching passed ID
func (route PilotRoutes) Get(c *gin.Context) {
	db := c.MustGet("DB").(*sql.DB)
	log := c.MustGet("logger").(log15.Logger)

	paramID := c.Param("id")
	id, _ := strconv.Atoi(paramID)

	pilot, err := models.FindPilot(db, id)
	if err != nil {
		log.Error("db: failed to get pilot", "id", id)
	}

	log.Info("db: fetched pilot", "id", id)
	c.JSON(200, pilot)
}

// GetAll : Get all pilots
func (route PilotRoutes) GetAll(c *gin.Context) {
	db := c.MustGet("DB").(*sql.DB)
	log := c.MustGet("logger").(log15.Logger)

	pilots, err := models.Pilots(db).All()
	if err != nil {
		log.Error("db: failed to get pilots", "err", err)
	}

	log.Info("db: fetched pilots", "count", len(pilots))
	c.JSON(200, pilots)
}

// Create : Create pilot with the passed name string
func (route PilotRoutes) Create(c *gin.Context) {
	db := c.MustGet("DB").(*sql.DB)
	log := c.MustGet("logger").(log15.Logger)

	var json models.Pilot
	if c.BindJSON(&json) != nil {
		log.Error("gin: error creating pilot")
	}

	pilot := &models.Pilot{
		Name: json.Name,
	}
	if err := pilot.Insert(db); err != nil {
		log.Error("db: failed to insert pilot", "err", err)
	}

	log.Info("db: insterted pilot", "id", pilot.ID)
}

// Update : Attempts to update the pilot matching the passed id
func (route PilotRoutes) Update(c *gin.Context) {
	db := c.MustGet("DB").(*sql.DB)
	log := c.MustGet("logger").(log15.Logger)

	id, _ := strconv.Atoi(c.Param("id"))
	var json models.Pilot
	if c.BindJSON(&json) != nil {
		log.Error("gin: error creating pilot")
	}

	pilot, _ := models.FindPilot(db, id)
	pilot.Name = json.Name
	if err := pilot.Update(db); err != nil {
		log.Error("db: failed to update pilot", "id", id, "err", err)
	}

	log.Info("db: updated pilot", "id", id)
}

// Delete : Attempts to delete the pilot matching the passed id
func (route PilotRoutes) Delete(c *gin.Context) {
	db := c.MustGet("DB").(*sql.DB)
	log := c.MustGet("logger").(log15.Logger)

	id, _ := strconv.Atoi(c.Param("id"))

	pilot, _ := models.FindPilot(db, id)
	if err := pilot.Delete(db); err != nil {
		log.Error("db: failed to delete pilot", "id", id, "err", err)
	}

	log.Info("db: deleted pilot", "id", id)
}