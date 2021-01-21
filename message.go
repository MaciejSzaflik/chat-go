package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

//obsolete
type Msg struct {
	Data    []byte
	MsgType int
}

type ComplexMsg struct {
	/*
		comands:
		JoinChannel
		CreateChannel
		LeaveChannel
		Message
	*/
	Comand string `json:"comand"`
	/*
		Value for channel related is name of channel
		Value for message is message
	*/
	Value string `json:"value"`

	client *websocket.Conn
}

func (cMsg ComplexMsg) String() string {
	return fmt.Sprintf("cmsg %s | %s", cMsg.Comand, cMsg.Value)
}
