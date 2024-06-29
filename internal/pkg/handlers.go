package pkg

import (
	"github.com/Chanmachan/GoChat/internal/models"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: 後で変える
		return true
	},
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	// urlからルーム番号を取得
	// ルーム番号のエラーチェック
	roomNumberStr := r.URL.Query().Get("room")
	if roomNumberStr == "" {
		log.Println("room number is empty")
		return
	}
	// ルーム番号をint型に変換
	roomNumber, err := strconv.Atoi(roomNumberStr)
	if err != nil {
		log.Println("room number is not number")
		return
	} else if roomNumber > 100 {
		log.Println("room number is over 100")
		return
	}
	// ルーム番号が存在しない場合はルーム作成
	// 入力された番号ごとにルームを作成
	room, exists := rooms[roomNumber]
	if !exists {
		room = CreateNewRoom(roomNumber)
	}
	// クライアントからのWebSocketの接続を処理
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// ルームにクライアントを追加
	room.Clients[conn] = true
	log.Printf("room: %d, clients: %d", room.RoomNumber, len(room.Clients))
	// 関数が終わるときにクライアントが削除される
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
	}(conn)
	// メッセージの受信
	go handleRoomMessages(room)
	for {
		var msg models.Message
		// 新しいメッセージをJSONとして読み込み、Message構造体にマッピング
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println(err)
			delete(room.Clients, conn)
			break
		}
		room.Broadcast <- msg
	}
}

func handleRoomMessages(room *Room) {
	for msg := range room.Broadcast {
		for client := range room.Clients {
			if err := client.WriteJSON(msg); err != nil {
				log.Printf("Error sending message: %v", err)
				err := client.Close()
				if err != nil {
					return
				}
				delete(room.Clients, client)
			}
		}
	}
}
