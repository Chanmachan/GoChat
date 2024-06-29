package main

import (
	"github.com/Chanmachan/GoChat/internal/api"
	"github.com/Chanmachan/GoChat/internal/pkg"
	"github.com/Chanmachan/GoChat/pkg/auth"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func main() {
	if godotenv.Load() != nil {
		log.Fatal("Failed to load .env file")
	}
	router := mux.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Reactアプリのオリジンを許可
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		AllowCredentials: true, // 認証情報を許可
	})
	// auth
	auth.SetUp()
	//http.HandleFunc("/auth/", auth.LoginHandler)
	//http.HandleFunc("/auth/callback", auth.CallBackHandler)
	router.HandleFunc("/auth/", auth.LoginHandler).Methods("GET")
	router.HandleFunc("/auth/callback", auth.CallBackHandler).Methods("GET")
	// ファイルを指定する
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/", fs)
	// api
	// ハンドラーの設定
	router.HandleFunc("/api/userinfo", api.UserInfoHandler).Methods("GET")
	router.HandleFunc("/api/login", api.GetUserHandler).Methods("GET")
	// クライアントからのWebSocketの接続を処理
	router.HandleFunc("/ws", pkg.HandleConnections)

	// CORSミドルウェアを使用
	handler := c.Handler(router)

	err := http.ListenAndServe(":9090", handler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
		return
	}
}
