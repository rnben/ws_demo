package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	dialer := websocket.Dialer{}
	connect, _, err := dialer.Dial("ws://127.0.0.1:8000", nil)
	if nil != err {
		log.Println(err)
		return
	}
	defer connect.Close()

	go tickWriter(connect)

	for {
		messageType, messageData, err := connect.ReadMessage()
		if nil != err {
			log.Println(err)
			break
		}
		switch messageType {
		case websocket.TextMessage: //文本数据
			fmt.Print(string(messageData))
		case websocket.BinaryMessage: //二进制数据
			fmt.Println(messageData)
		case websocket.CloseMessage: //关闭
		case websocket.PingMessage: //Ping
		case websocket.PongMessage: //Pong
		default:

		}
	}
}

func tickWriter(connect *websocket.Conn) {
	type Message struct {
		Author  string `json:"author"`
		Content string `json:"content"`
	}

	msg := Message{Author: os.Args[1], Content: "register"}
	b, _ := json.Marshal(msg)
	err := connect.WriteMessage(websocket.TextMessage, b)
	if nil != err {
		log.Println(err)
		return
	}

	msg = Message{Author: os.Args[1], Content: "ping"}
	b, _ = json.Marshal(msg)

	for {
		err := connect.WriteMessage(websocket.TextMessage, b)
		if nil != err {
			log.Println(err)
			break
		}
		time.Sleep(time.Second)
	}
}
