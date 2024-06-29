package api

import (
	"encoding/json"
	"github.com/Chanmachan/GoChat/pkg/auth"
	"log"
	"net/http"
)

// UserInfoHandler - 認証されたユーザーの情報をフロントエンドに送信するためのハンドラー
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	session, err := auth.GetStore().Get(r, "auth-session")
	if err != nil {
		http.Error(w, "Failed get session error: "+err.Error(), http.StatusInternalServerError)
		log.Fatal("Failed get session")
		return
	}
	userinfo, ok := session.Values["userInfo"].(auth.UserInfo)
	if !ok || userinfo.Name == "" {
		http.Error(w, "No user data available", http.StatusNotFound)
		log.Fatal("No user data available")
		return
	}
	respData, err := json.Marshal(userinfo)
	if err != nil {
		http.Error(w, "Failed marshal user data: "+err.Error(), http.StatusInternalServerError)
		log.Fatal("Failed marshal user data")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(respData)
}
