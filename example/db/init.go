package db

import (
	"database/sql"
	IN "example/interfaces"
	"fmt"
)

const (
	host     = ""
	port     = 5432
	user     = "postgres"
	password = ""
	dbname   = "postgres"
)

func DBCONNECT() *sql.DB {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)

	IN.CheckError(err)

	// check db
	err = db.Ping()

	IN.CheckError(err)

	return db

}
