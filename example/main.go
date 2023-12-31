package main

import (
	"context"
	C "example/controller"
	DB "example/db"
	"flag"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	var dir string

	flag.StringVar(&dir, "dir", "static/", "File server")
	flag.Parse()

	fs := http.FileServer(http.Dir(dir))
	r := mux.NewRouter().StrictSlash(true)

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	r.HandleFunc("/", C.Home).Methods("GET")
	r.HandleFunc("/hello", C.Hello).Methods("GET")
	r.HandleFunc("/greet/{name}/", C.Greeting).Methods("GET")

	db := DB.DBCONNECT()
	c := C.Controller(db)

	r.HandleFunc("/users/", c.GetUsers).Methods("GET")
	r.HandleFunc("/users/register/", c.CreateTemplate).Methods("GET")
	r.HandleFunc("/users/create/", c.CreateUser).Methods("POST")

	svr := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:4001",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := svr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("ListenAndServe(): s%", err)
		}
	}()

	<-stop

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := svr.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	} else {
		log.Println("Server stopped gracefully")
	}
}
