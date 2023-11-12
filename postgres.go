package main

import (
	"database/sql"
	_ "github.com/lib/pq"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type postgresClient struct {
	PostgresUri string
	Logger      log.Logger
}

func newPostgresClient(postgresUri string, logger log.Logger) postgresClient {
	return postgresClient{PostgresUri: postgresUri, Logger: logger}
}

func (c postgresClient) Put(
	pac float64, pacUnit string,
	day float64, dayUnit string,
	year float64, yearUnit string,
	total float64, totalUnit string,
) {

	db, err := sql.Open("postgres", c.PostgresUri)
	defer db.Close()

	if err != nil {
		level.Error(c.Logger).Log("msg", "postgresClient connection", "error", err)
	}

	if err = db.Ping(); err != nil {
		level.Error(c.Logger).Log("msg", "postgresClient ping", "error", err)
	}

	// dynamic
	query := `INSERT INTO "fronius" 
    	("pac", "pac_unit", "day", "day_unit", "year", "year_unit", "total", "total_unit") 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err = db.Exec(query, pac, pacUnit, day, dayUnit, year, yearUnit, total, totalUnit)
}
