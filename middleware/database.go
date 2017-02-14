package middleware

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // import postgres driver
	"gopkg.in/inconshreveable/log15.v2"
)

// Database : Middleware that opens db connection
func Database(log log15.Logger) gin.HandlerFunc {
	db, err := sql.Open("postgres", connString())
	if err != nil {
		log.Error("failed to open database", "err", err)
		return nil
	}

	return func(c *gin.Context) {
		c.Set("DB", db)
		c.Next()
	}
}

// TODO: Use environment variables here or a config file in the future
func connString() string {
	return fmt.Sprintf("dbname=%s host=%s user=%s password=%s sslmode=disable",
		"flight", "localhost", "boiler", "boiler")
}
