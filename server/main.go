package main

import (
	"golang.org/x/net/websocket"
	"log"
	"os"
	"net/http"
)
type ClientConn struct {
	ws *websocket.Conn
	ip  string
}
type Message struct {
	Timestamp int
	Data string
}
var master *ClientConn = nil

func echoHandler(ws *websocket.Conn) {
	var err error
	var message Message

	client := ClientConn{ws, ws.Request().RemoteAddr}
	log.Println("Client connected:", client.ip)

	for {
		// получили сообщенько
		if err = websocket.JSON.Receive(ws, &message); err != nil {
			log.Println("Disconnected waiting", err.Error())
			return
		}
		// разбираем, назначаем мастером того, кто так представился
		if message.Data == "master" {
			master = &client
			log.Println("Master client:", master.ip)
			continue;
		}
		// если не мастера, то некому слать
		if master == nil {
			continue
		}
		// шлем то, что пришло, мастеру
		if err = websocket.JSON.Send(master.ws, message); err != nil {
			log.Println("Could not send message to ", master.ip, err.Error())
		}
  	}
}

func main() {
	var pwd, _ = os.Getwd()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	http.Handle("/echo", websocket.Handler(echoHandler))
	http.Handle("/", http.FileServer(http.Dir(pwd + "/public")))

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

