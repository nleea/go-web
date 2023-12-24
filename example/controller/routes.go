package controller

import (
	"database/sql"
	"encoding/json"
	IN "example/interfaces"
	RP "example/repository"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Greet struct {
	name string
}

type C IN.UseController

func Home(w http.ResponseWriter, r *http.Request) {
	tmpt := template.Must(template.ParseFiles("./layout.html"))

	err := tmpt.Execute(w, "")

	if err != nil {
		fmt.Println("Error:", err)
	}

}

func Hello(w http.ResponseWriter, r *http.Request) {
	tmpt := template.Must(template.ParseFiles("./hello.html"))

	err := tmpt.Execute(w, "")

	if err != nil {
		fmt.Println("Error:", err)
	}

}

func Greeting(w http.ResponseWriter, r *http.Request) {
	tmpt := template.Must(template.ParseFiles("./greet.html"))

	vars := mux.Vars(r)

	name := vars["name"]

	data := Greet{
		name: name,
	}

	err := tmpt.Execute(w, data)

	if err != nil {
		fmt.Println("Error:", err)
	}

}

func Controller(db *sql.DB) *C {
	return &C{
		DB: db,
	}
}

func (uc *C) CreateUser(w http.ResponseWriter, req *http.Request) {

	fmt.Println(req.Method)

	if req.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	contentType := req.Header.Get("Content-Type")

	fmt.Println(contentType)

	if contentType != "application/json; charset=utf-8" {
		http.Error(w, "Tipo de contenido incorrecto, se esperaba application/json", http.StatusUnsupportedMediaType)
		return
	}

	var user IN.User

	dataBy, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	unerr := json.Unmarshal(dataBy, &user)
	if unerr != nil {
		http.Error(w, unerr.Error(), http.StatusBadRequest)
		return
	}

	_, err = json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	RP.InsertDB(user, uc.DB)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Datos guardados con éxito"))

}

func (uc *C) GetUsers(w http.ResponseWriter, req *http.Request) {

	rows := RP.GetAll(uc.DB)

	var users []IN.User

	for rows.Next() {
		var u IN.User
		err := rows.Scan(&u.Id, &u.Username, &u.Password, &u.Created_at)
		IN.CheckError(err)
		users = append(users, u)

	}

	template := template.Must(template.ParseFiles("./users.html"))

	data := IN.RenderData{
		Title: "Users",
		Users: users,
	}

	templateError := template.Execute(w, data)

	IN.CheckError(templateError)

	err2 := rows.Err()

	IN.CheckError(err2)
}
