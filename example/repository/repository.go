package repository

import (
	"database/sql"
	IN "example/interfaces"
	_ "github.com/lib/pq"
)

func GetAll(db *sql.DB) *sql.Rows {
	getDynStmt := `SELECT id, username, password, created_at FROM users`

	rows, err := db.Query(getDynStmt)

	IN.CheckError(err)

	return rows
}

func InsertDB(user IN.User, db *sql.DB) {
	insertDynStmt := `insert into "users"("username", "password") values($1, $2)`
	_, err := db.Exec(insertDynStmt, user.Username, user.Password)
	IN.CheckError(err)
}
