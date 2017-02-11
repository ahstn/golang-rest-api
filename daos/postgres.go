package daos

import (
	"database/sql"

	_ "github.com/lib/pq"
	"gopkg.in/inconshreveable/log15.v2"
)

func Open(log log15.Logger) *sql.DB {
	db, err := sql.Open("postgres", `dbname=flight host=localhost user=boiler password=boiler sslmode=disable`)
	if err != nil {
		log.Error("failed to open database", "err", err)
		return nil
	}
	return db
}
