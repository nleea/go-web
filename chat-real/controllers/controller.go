package controllers

import (
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to chat room"))
}
