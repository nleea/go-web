package main

import (
	C "chat/controllers"
	_ "chat/socket"
	"context"
	"flag"
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
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Serving at localhost:8000...")

	http.HandleFunc("/", C.Home)

	svr := &http.Server{
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
