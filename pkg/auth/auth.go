package auth

import (
	"context"
	"encoding/json"
	"github.com/Chanmachan/GoChat/pkg/random"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"log"
	"net/http"
	"os"
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
	store        *sessions.CookieStore
	oauth2Config *oauth2.Config
)

func GetStore() *sessions.CookieStore {
	return store
}

func SetUp() {
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
	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "auth-session")
	// OAuthリダイレクトの際にCSRF攻撃を防ぐために使用され、認証要求が開始された同じユーザーによって完了されることを確認する
	state := random.GenerateRandomString()
	// PKCE（Proof Key for Code Exchange）プロセスで使用され、認証コードをトークンに交換する際の追加の保護
	verifier := oauth2.GenerateVerifier()
	session.Values["state"] = state
	session.Values["verifier"] = verifier
	session.Save(r, w)
	// リクエスト先のURLを作成する
	// AccessTypeOfflineを設定するとリフレッシュトークンの期限が無期限になる
	url := oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// CallBackHandler 認可サーバーからのリダイレクトに対するハンドラー
func CallBackHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, "Failed to get session: "+err.Error(), http.StatusInternalServerError)
		log.Println("Failed to get session: " + err.Error())
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
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		log.Println("Failed to exchange token: " + err.Error())
		return
	}
	// トークンを使い、HTTPクライアントを取得
	client := oauth2Config.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		log.Println("Failed to get user info: " + err.Error())
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
		http.Error(w, "Failed to read response body: "+err.Error(), http.StatusInternalServerError)
		log.Println("Failed to read response body: " + err.Error())
		return
	}
	var userInfo UserInfo
	if err := json.Unmarshal(data, &userInfo); err != nil {
		http.Error(w, "Failed to unmarshal user info: "+err.Error(), http.StatusInternalServerError)
		log.Println("Failed to unmarshal user info: " + err.Error())
		return
	}
	// ユーザー情報をセッションに保存 -> セッションの状態がSaveで自動的にクライアントに同期
	session.Values["access_token"] = token.AccessToken
	session.Values["refresh_token"] = token.RefreshToken
	session.Values["expiry"] = token.Expiry
	session.Values["userInfo"] = userInfo
	log.Printf("Retrieved UserInfo: %+v", userInfo)
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
