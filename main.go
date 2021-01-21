package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var id int
var mainHandler = MainHandler{}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Index")
}

func channelEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("Channel Endpoint")
	conn, err := mainHandler.upgrader.Upgrade(w, r, nil)

	mainHandler.clients[conn] = id
	id++

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
			delete(mainHandler.clients, conn)
			break
		}

		var complexMessage ComplexMsg
		log.Println(string(message))
		err = json.Unmarshal(message, &complexMessage)
		if err == nil {
			complexMessage.client = conn
			mainHandler.broadcast <- complexMessage
		} else {
			log.Println(err)
		}
	}
}

func handleChannelMessages() {
	for {
		msg := <-mainHandler.broadcast

		switch {
		case msg.Comand == "JoinChannel":
			JoinChannel(msg.client, msg.Value)
		case msg.Comand == "CreateChannel":
			CreateChannel(msg.client, msg.Value)
		case msg.Comand == "LeaveChannel":
			LeaveChat(msg.client)
		case msg.Comand == "Message":
			SendMessageToChannel(msg.client, msg.Value)
		default:
			log.Println("Comand not recognized")
		}
	}
}

func setupRoutes() {
	http.HandleFunc("/", index)
	http.HandleFunc("/ws", channelEndpoint)
}

func main() {
	fmt.Println("Start server")
	mainHandler.upgrader = websocket.Upgrader{}
	mainHandler.clients = make(map[*websocket.Conn]int)
	mainHandler.broadcast = make(chan ComplexMsg)

	mainHandler.channelsBroadcasts = make(map[string](chan string))
	mainHandler.channelsClients = make(map[string]map[*websocket.Conn]bool)
	mainHandler.clientInChannel = make(map[*websocket.Conn]string)

	mainHandler.upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	setupRoutes()
	go handleChannelMessages()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
