package main

import (
	"fmt"
	"github.com/googollee/go-socket.io"
	"log"
	"net/http"
)

type response map[string]interface{}

func main() {

	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/", "msg", func(s socketio.Conn, resp response) {
		fmt.Println(resp["msg"])
		// log.Println("message:", resp.msg)
		// s.Emit("reply", "Message received: "+resp.msg)
		s.Emit("rooms", s.Rooms())
	})

	// server.OnEvent("/", "rooms", func(s socketio.Conn) {
	// 	s.Emit("rooms", s.Rooms())
	// })

	// server.OnEvent("/", "reply", func(s socketio.Conn, msg string) {
	// 	s.Emit("reply", msg)
	// })

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("disconnected:", s.ID())
	})

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))

}
