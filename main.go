package main

//go:generate sqlboiler postgres

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/phazyy/golang-rest-api/middleware"
	"github.com/phazyy/golang-rest-api/routes"
	"gopkg.in/inconshreveable/log15.v2"
)

var log = log15.New()

func main() {
	r := gin.Default()
	r.Use(middleware.Logger(log))
	r.Use(middleware.Database(log))

	v1 := r.Group("/v1")
	{
		pilot := new(routes.PilotRoutes)

		v1.GET("/pilots", pilot.GetAll)
		v1.GET("/pilots/:id", pilot.Get)
		v1.POST("/pilots", pilot.Create)
		v1.PUT("/pilots/:id", pilot.Update)
		v1.DELETE("/pilots/:id", pilot.Delete)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"status":  "404",
			"message": "Not Found.",
		})
	})

	r.Run(":8080")
}
