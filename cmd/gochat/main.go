package main

import (
	"github.com/Chanmachan/GoChat/internal/pkg"
	"github.com/Chanmachan/GoChat/pkg/auth"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {
	if godotenv.Load() != nil {
		log.Fatal("Failed to load .env file")
	}
	// auth
	auth.SetUp()
	http.HandleFunc("/auth/", auth.LoginHandler)
	http.HandleFunc("/auth/callback", auth.CallBackHandler)
	// ファイルを指定する
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/", fs)
	// クライアントからのWebSocketの接続を処理
	http.HandleFunc("/ws", pkg.HandleConnections)

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
		return
	}
}
