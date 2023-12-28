package db

import (
	"database/sql"
	IN "example/interfaces"
	"fmt"
)

const (
	host     = "db.seoqxlxkudnimkoezefb.supabase.co"
	port     = 5432
	user     = "postgres"
	password = "NBunX9sBoovmVOZQ"
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
