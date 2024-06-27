package pkg

import "github.com/gorilla/websocket"

type Room struct {
	RoomNumber int
	// ルームごとにクライアントが存在
	Clients map[*websocket.Conn]bool
	// それぞれのルームにチャンネルを一つ持たせる
	Broadcast chan Message
}

type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

var rooms = make(map[int]*Room)

func CreateNewRoom(roomNumber int) *Room {
	r := &Room{
		RoomNumber: roomNumber,
		Clients:    make(map[*websocket.Conn]bool),
		Broadcast:  make(chan Message),
	}
	rooms[roomNumber] = r
	return r
}
