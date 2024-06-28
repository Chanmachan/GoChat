package auth

import (
	"context"
	"fmt"
	"github.com/Chanmachan/GoChat/pkg/random"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"log"
	"net/http"
	"os"
)

//var (
//	state        string
//	verifier     string
//	oauth2Config *oauth2.Config // グローバル変数として設定を保持
//)
//

var (
	store        *sessions.CookieStore
	oauth2Config *oauth2.Config
)

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
	session, _ := store.Get(r, "session-name")
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
	session, _ := store.Get(r, "session-name")
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
		return
	}
	// トークンを使い、HTTPクライアントを取得
	client := oauth2Config.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
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
		return
	}
	fmt.Println(string(data))
}