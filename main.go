package main

//go:generate sqlboiler postgres

import (
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/phazyy/golang-rest-api/daos"
	"github.com/phazyy/golang-rest-api/models"
	"gopkg.in/inconshreveable/log15.v2"
)

var log = log15.New()

func main() {
	r := gin.Default()

	db := daos.Open(log)

	r.GET("/pilots", func(c *gin.Context) {
		pilots := daos.GetAll(db, log)
		c.JSON(200, pilots)
	})

	r.POST("/pilots", func(c *gin.Context) {
		var json models.Pilot
		if c.BindJSON(&json) != nil {
			log.Error("gin: error creating creating")
		}

		daos.Create(json.Name, db, log)
	})

	r.PUT("/pilots/:id", func(c *gin.Context) {
		paramID := c.Param("id")
		var json models.Pilot
		if c.BindJSON(&json) != nil {
			log.Error("gin: error updating pilot")
		}

		id, _ := strconv.Atoi(paramID)
		daos.Update(id, json.Name, db, log)
	})

	r.DELETE("/pilots/:id", func(c *gin.Context) {
		paramID := c.Param("id")

		id, _ := strconv.Atoi(paramID)
		daos.Delete(id, db, log)
	})

	r.Run(":8080")
}
