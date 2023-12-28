package controllers

import (
	IN "chat/interfaces"
	SK "chat/socket"
	"fmt"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to chat room"))
}

func HandleConnetions(w http.ResponseWriter, r *http.Request) {
	conn, err := SK.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	SK.Clients[conn] = true

	for {
		var msg IN.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println(err)
			delete(SK.Clients, conn)
			return
		}

		SK.Broadcast <- msg
	}
}
