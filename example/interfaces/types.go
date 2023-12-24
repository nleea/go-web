package interfaces

import (
	"database/sql"
	"time"
)

type User struct {
	Id         int       `json:"id"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	Created_at time.Time `json:"created_at"`
}

type RenderData struct {
	Title string
	Users []User
}

type UseController struct {
	DB *sql.DB
}
