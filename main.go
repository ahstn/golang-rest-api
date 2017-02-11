package main

//go:generate sqlboiler postgres

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/phazyy/golang-rest-api/daos"
	"gopkg.in/inconshreveable/log15.v2"
)

var log = log15.New()

func main() {
	r := gin.Default()

	db := daos.Open(log)
	pilots := daos.GetPilots(db, log)

	r.GET("/pilots", func(c *gin.Context) {
		c.JSON(200, pilots)
	})

	r.Run(":8080")
}
