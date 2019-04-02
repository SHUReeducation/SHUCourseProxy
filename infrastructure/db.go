package infrastructure

import (
	"database/sql"
	_ "github.com/lib/pq"
	"os"
)

var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("postgres", os.Getenv("DB_ADDRESS"))
	if err != nil {
		panic(err)
	}
}
