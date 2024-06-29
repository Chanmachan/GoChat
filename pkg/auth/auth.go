package auth

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"github.com/Chanmachan/GoChat/pkg/random"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type UserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

var (
	Store        *sessions.CookieStore
	oauth2Config *oauth2.Config
)

func SetUp() {
	gob.Register(time.Time{})
	gob.Register(UserInfo{})
	oauth2Config = &oauth2.Config{
		ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:9090/auth/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	Store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 24, // 24 hours
		HttpOnly: true,
		//Secure:   true,                 // HTTPS環境下でのみ設定
		SameSite: http.SameSiteLaxMode, // ブラウザによるクロスサイトリクエストの扱い
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "auth-session")
	state := random.GenerateRandomString()
	verifier := oauth2.GenerateVerifier()
	session.Values["state"] = state
	session.Values["verifier"] = verifier
	session.Save(r, w)
	url := oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]string{"url": url})
	if err != nil {
		handleError(w, "Failed to encode response", err, http.StatusInternalServerError)
		return
	}
}

// CallBackHandler 認可サーバーからのリダイレクトに対するハンドラー
func CallBackHandler(w http.ResponseWriter, r *http.Request) {
	session, err := Store.Get(r, "auth-session")
	if err != nil {
		handleError(w, "Failed to get session: ", err, http.StatusInternalServerError)
		return
	}
	if r.URL.Query().Get("state") != session.Values["state"] {
		http.Error(w, "State did not match", http.StatusBadRequest)
		return
	}
	// 設定したstateと一致しているかを確かめる
	// Exchangeでverifierを使えばまとめてできる？
	// 認証コードをトークンに変換
	token, err := oauth2Config.Exchange(context.Background(), r.URL.Query().Get("code"), oauth2.SetAuthURLParam("code_verifier", session.Values["verifier"].(string)))
	if err != nil {
		handleError(w, "Failed to exchange token: ", err, http.StatusInternalServerError)
		return
	}
	// トークンを使い、HTTPクライアントを取得
	client := oauth2Config.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		handleError(w, "Failed to get user info: ", err, http.StatusInternalServerError)
		return
	}
	// bodyをcloseするのは呼び出し側の責任
	// fdの枯渇を防ぐため
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		handleError(w, "Failed to read response body: ", err, http.StatusInternalServerError)
		return
	}
	var userInfo UserInfo
	if err := json.Unmarshal(data, &userInfo); err != nil {
		handleError(w, "Failed to unmarshal user info: ", err, http.StatusInternalServerError)
		return
	}
	// ユーザー情報をセッションに保存 -> セッションの状態がSaveで自動的にクライアントに同期
	session.Values["access_token"] = token.AccessToken
	session.Values["refresh_token"] = token.RefreshToken
	session.Values["expiry"] = token.Expiry
	session.Values["userInfo"] = userInfo
	if err := session.Save(r, w); err != nil {
		handleError(w, "Failed to save session: ", err, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleError(w http.ResponseWriter, msg string, err error, statusCode int) {
	log.Printf("%s: %s\n", msg, err.Error())
	http.Error(w, msg+err.Error(), statusCode)
}
