package api

import (
	"encoding/json"
	"github.com/Chanmachan/GoChat/pkg/auth"
	"log"
	"net/http"
)

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("start GetUserHandler")
	// ここでユーザーの認証状態をチェックするロジックを実装
	// 例: sessionからユーザー情報を取得し、存在しない場合はリダイレクト
	session, err := auth.Store.Get(r, "auth-session")
	if err != nil || session.Values["userInfo"] == nil {
		// 認証されていないユーザーの場合、OAuth認証プロセスを開始
		log.Println("trying redirect ... ")
		authURL := "http://localhost:9090/auth/" // OAuth開始のためのURL
		http.Redirect(w, r, authURL, http.StatusSeeOther)
		return
	}

	// 認証済みのユーザーの場合、ユーザー情報を返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session.Values["userInfo"])

	log.Println("end GetUserHandler")
}

// UserInfoHandler - 認証されたユーザーの情報をフロントエンドに送信するためのハンドラー
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	session, err := auth.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, "Failed get session error: "+err.Error(), http.StatusInternalServerError)
		log.Fatal("Failed get session")
		return
	}
	for key, value := range session.Values {
		log.Printf("Session key: %s, value: %v, type: %T", key, value, value)
	}
	name, ok := session.Values["name"].(string)
	if !ok || name == "" {
		http.Error(w, "Failed get session name", http.StatusInternalServerError)
		log.Fatal("Failed get session name")
	}
	respData := struct {
		Name string `json:"name"`
	}{
		Name: name,
	}
	//userinfo, ok := session.Values["userInfo"].(auth.UserInfo)
	//if !ok || userinfo.Name == "" {
	//	if !ok {
	//		log.Println("not ok")
	//	}
	//	http.Error(w, "No user data available", http.StatusNotFound)
	//	log.Fatal("No user data available")
	//	return
	//}
	//respData, err := json.Marshal(userinfo)
	respDataJson, err := json.Marshal(respData)
	if err != nil {
		http.Error(w, "Failed marshal user data: "+err.Error(), http.StatusInternalServerError)
		log.Fatal("Failed marshal user data")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(respDataJson)
}
