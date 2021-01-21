package main

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

func SendMessageToChannel(creator *websocket.Conn, msgValue string) {
	fmt.Println("channel exist")
	channelName, ok := mainHandler.clientInChannel[creator]
	if ok {
		fmt.Println("client in channel")
		mainHandler.channelsBroadcasts[channelName] <- msgValue
	} else {
		fmt.Println("client not in channel")
	}
}

func CreateChannel(creator *websocket.Conn, channelName string) {
	_, ok := mainHandler.channelsClients[channelName]
	if ok {
		WriteToClient(creator, fmt.Sprintf("[Name taken %s]", channelName))
	} else {

		mainHandler.channelsClients[channelName] = make(map[*websocket.Conn]bool)
		mainHandler.channelsClients[channelName][creator] = true
		mainHandler.clientInChannel[creator] = channelName
		mainHandler.channelsBroadcasts[channelName] = make(chan string)

		WriteToClient(creator, fmt.Sprintf("[Channel created %s]", channelName))

		go HandleMessagesInChannel(mainHandler.channelsBroadcasts[channelName], channelName)
	}
}

func HandleMessagesInChannel(channelChannel chan string, name string) {
	for {
		msg := <-channelChannel
		for client, _ := range mainHandler.channelsClients[name] {
			err := client.WriteMessage(1, []byte(msg))
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()

				mainHandler.RemoveClient(client) //this isnt safe?
			}
		}
	}
}

//duplicated from above ...
func WriteToClient(client *websocket.Conn, msg string) {
	err := client.WriteMessage(1, []byte(msg))
	if err != nil {
		log.Printf("error: %v", err)
		client.Close()

		mainHandler.RemoveClient(client)
	}
}

func JoinChannel(joiner *websocket.Conn, channelName string) {
	_, checkChannel := mainHandler.channelsClients[channelName]
	_, checkClient := mainHandler.clientInChannel[joiner]

	if checkClient {
		WriteToClient(joiner, fmt.Sprintf("[You are in channel already %s]", channelName))
		return
	}

	if checkChannel {
		fmt.Println("channel exist") //send this to client
		mainHandler.channelsClients[channelName][joiner] = true
		mainHandler.clientInChannel[joiner] = channelName

		fmt.Println(len(mainHandler.channelsClients[channelName]))
		WriteToClient(joiner, fmt.Sprintf("[Succesfuly joined channel %s]", channelName))
		mainHandler.channelsBroadcasts[channelName] <- "Someone Joined"

	} else {
		fmt.Println("channel dont exist")
		WriteToClient(joiner, fmt.Sprintf("Channel doesnt exist %s]", channelName))
	}
}

func LeaveChat(leaver *websocket.Conn) {
	channelName, ok := mainHandler.clientInChannel[leaver]
	if ok {
		fmt.Println("client left")
		delete(mainHandler.clientInChannel, leaver)
		delete(mainHandler.channelsClients[channelName], leaver)
		WriteToClient(leaver, fmt.Sprintf("[You have left %s]", channelName))
		mainHandler.channelsBroadcasts[channelName] <- "Someone Left"
	} else {
		fmt.Println("client wasnt in any channel")
		WriteToClient(leaver, fmt.Sprintf("[You were not in channel %s]", channelName))
	}
}
