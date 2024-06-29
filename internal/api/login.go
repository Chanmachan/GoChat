package api

import (
	"encoding/json"
	"github.com/Chanmachan/GoChat/pkg/auth"
	"log"
	"net/http"
)

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GetUserHandler")
	session, err := auth.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}
	// Log all session key-value pairs
	log.Println("Session Values:")
	for key, value := range session.Values {
		log.Printf("Key: %v, Value: %v\n", key, value)
	}
	// Check if user info is present in session
	if userInfo, ok := session.Values["userInfo"].(auth.UserInfo); ok {
		// User is authenticated, return user info
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(userInfo)
	} else {
		// User is not authenticated, return unauthorized
		log.Println("Failed")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}
