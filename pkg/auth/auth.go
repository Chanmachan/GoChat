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
	"net/url"
	"os"
	"sync"
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

// SessionData 認証の最中に一時的に保持するデータの集合
type SessionData struct {
	// OAuthリダイレクトの際にCSRF攻撃を防ぐために使用され、認証要求が開始された同じユーザーによって完了されることを確認する
	State string
	// PKCE（Proof Key for Code Exchange）プロセスで使用され、認証コードをトークンに交換する際の追加の保護
	Verifier string
}

// SessionManager AuthSessionDataをまとめて管理
// mapにアクセスするゴルーチンは一つなのを保証するためにmutex
type SessionManager struct {
	sessions map[string]*SessionData
	mu       sync.Mutex
}

func NewAuthSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*SessionData),
	}
}

// CreateSession AuthSessionDataを作成して、格納されているセッションのIDを返す
func (m *SessionManager) CreateSession() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	temporaryID := random.GenerateRandomString()
	sessionData := &SessionData{
		State:    random.GenerateRandomString(),
		Verifier: oauth2.GenerateVerifier(),
	}

	m.sessions[temporaryID] = sessionData
	return temporaryID
}

// GetSession 一意に与えられたIDのAuthSessionDataを取得する
func (m *SessionManager) GetSession(temporaryID string) (*SessionData, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	sessionData, exists := m.sessions[temporaryID]
	return sessionData, exists
}

func (m *SessionManager) RemoveSession(temporaryID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.sessions, temporaryID)
}

var (
	sessionManager *SessionManager
	store          *sessions.CookieStore
	oauth2Config   *oauth2.Config
)

func SetUp() {
	sessionManager = NewAuthSessionManager()
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

	temporaryID := sessionManager.CreateSession()
	sessionData, _ := sessionManager.GetSession(temporaryID)

	// リクエスト先のURLを作成する
	// AccessTypeOfflineを設定するとリフレッシュトークンの期限が無期限になる
	u, err := url.Parse(oauth2Config.AuthCodeURL(sessionData.State, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(sessionData.Verifier)))
	if err != nil {
		log.Fatal(err)
		return
	}
	q := u.Query()
	q.Set("temporary_id", temporaryID)
	u.RawQuery = q.Encode()
	http.Redirect(w, r, u.String(), http.StatusTemporaryRedirect)
}

// CallBackHandler 認可サーバーからのリダイレクトに対するハンドラー
func CallBackHandler(w http.ResponseWriter, r *http.Request) {
	// urlのクエリに設定しておいたidを取り出し、そのidをもとにdataをsessionManagerから取り出す
	temporaryID := r.URL.Query().Get("temporary_id")
	sessionData, exists := sessionManager.GetSession(temporaryID)
	if !exists {
		http.Error(w, "Session does not exist", http.StatusNotFound)
		return
	}
	// 設定したstateと一致しているかを確かめる
	if r.URL.Query().Get("state") != sessionData.State {
		http.Error(w, "State did not match", http.StatusBadRequest)
		return
	}

	// 設定したstateと一致しているかを確かめる
	// Exchangeでverifierを使えばまとめてできる？
	// 認証コードをトークンに変換
	token, err := oauth2Config.Exchange(context.Background(), r.URL.Query().Get("code"), oauth2.SetAuthURLParam("code_verifier", sessionData.Verifier))
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
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var userInfo UserInfo
	if err := json.Unmarshal(data, &userInfo); err != nil {
		http.Error(w, "Failed to unmarshal user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// ユーザー情報をセッションに保存 -> セッションの状態がSaveで自動的にクライアントに同期
	session, err := store.Get(r, "auth-session")
	session.Values["state"] = sessionData.State
	session.Values["verifier"] = sessionData.Verifier
	session.Values["userInfo"] = userInfo
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func GetStore() *sessions.CookieStore {
	return store
}
