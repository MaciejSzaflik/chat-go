package main

import (
	"github.com/gorilla/websocket"
)

type MainHandler struct {
	upgrader  websocket.Upgrader
	clients   map[*websocket.Conn]int
	broadcast chan ComplexMsg

	channelsBroadcasts map[string](chan string)
	channelsClients    map[string]map[*websocket.Conn]bool
	clientInChannel    map[*websocket.Conn]string
}
