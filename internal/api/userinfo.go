package api

import (
	"github.com/Chanmachan/GoChat/pkg/auth"
	"net/http"
)

// UserInfoHandler - 認証されたユーザーの情報をフロントエンドに送信するためのハンドラー
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	session, err := auth.GetStore().Get(r, "auth-session")
	if err != nil {
		http.Error(w, "Failed get session error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if data, ok := session.Values["userInfo"].(string); ok && data != "" {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(data))
	} else {
		http.Error(w, "No user data available", http.StatusNotFound)
	}
}
