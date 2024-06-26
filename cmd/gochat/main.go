package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var clients = make(map[*websocket.Conn]bool) // 接続されているクライアント
var broadcast = make(chan Message)           // ブロードキャストメッセージ

type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	clients[conn] = true
	for {
		var message Message
		// 新しいメッセージをJSONとして読み込み、Message構造体にマッピング
		err := conn.ReadJSON(&message)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, conn)
			break
		}
		// 受け取ったメッセージをbroadcastチャネルに送る
		broadcast <- message
	}
}

func handleMessages() {
	for {
		// メッセージを受信
		msg := <-broadcast
		// クライアントにメッセージを送信
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	// ファイルを指定する
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/template", fs)

	http.HandleFunc("/ws", handleConnections)
	go handleMessages()

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
		return
	}
}
