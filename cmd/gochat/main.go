package main

import (
	"github.com/Chanmachan/GoChat/internal/pkg"
	"log"
	"net/http"
)

func main() {
	// ファイルを指定する
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/template/", fs)
	// クライアントからのWebSocketの接続を処理
	http.HandleFunc("/ws", pkg.HandleConnections)
	go pkg.HandleMessages()

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
		return
	}
}
