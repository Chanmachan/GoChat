package pkg

import (
	"github.com/Chanmachan/GoChat/internal/models"
	"github.com/gorilla/websocket"
)

type Room struct {
	RoomNumber int
	// ルームごとにクライアントが存在
	Clients map[*websocket.Conn]bool
	// それぞれのルームにチャンネルを一つ持たせる
	Broadcast chan models.Message
}

var rooms = make(map[int]*Room)

func CreateNewRoom(roomNumber int) *Room {
	r := &Room{
		RoomNumber: roomNumber,
		Clients:    make(map[*websocket.Conn]bool),
		Broadcast:  make(chan models.Message),
	}
	rooms[roomNumber] = r
	return r
}
