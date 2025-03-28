package auth

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var store = sessions.NewCookieStore([]byte(getSessionKey()))

func getSessionName() string {
	name := os.Getenv("SESSION_NAME")
	if name == "" {
		fmt.Println("Warning: SESSION_NAME environment variable not set")
		return ""
	}
	return name
}

func getSessionKey() string {
	key := os.Getenv("SESSION_KEY")
	if key == "" {
		fmt.Println("Warning: SESSION_KEY environment variable not set")
		return ""
	}
	return key
}

func GetEmail() string {
	email := os.Getenv("AUTH_EMAIL")
	if email == "" {
		fmt.Println("Warning: AUTH_EMAIL environment variable not set")
		return ""
	}
	return email
}

func GetPasswordHash() string {
	hash := os.Getenv("AUTH_PASSWORD_HASH")
	if hash == "" {
		fmt.Println("Warning: AUTH_PASSWORD_HASH environment variable not set")
		return ""
	}
	return hash
}

func Login(w http.ResponseWriter, r *http.Request, email, password string) bool {
	storedEmail := GetEmail()
	if storedEmail == "" {
		fmt.Println("Login failed: AUTH_EMAIL not configured")
		return false
	}

	if subtle.ConstantTimeCompare([]byte(email), []byte(storedEmail)) != 1 {
		fmt.Println("email mismatch")
		return false
	}

	storedHash := GetPasswordHash()
	if storedHash == "" {
		fmt.Println("Login failed: AUTH_PASSWORD_HASH not configured")
		return false
	}

	hashBytes, err := base64.StdEncoding.DecodeString(storedHash)
	if err != nil {
		fmt.Printf("Error decoding hash: %v\n", err)
		return false
	}

	err = bcrypt.CompareHashAndPassword(hashBytes, []byte(password))
	if err != nil {
		fmt.Printf("Password comparison failed: %v\n", err)
		return false
	}

	session, _ := store.Get(r, getSessionName())
	session.Values["authenticated"] = true
	session.Values["login_time"] = time.Now().Format(time.RFC3339)
	err = session.Save(r, w)
	if err != nil {
		fmt.Printf("Error saving session: %v\n", err)
		return false
	}

	return true
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, getSessionName())
	session.Values["authenticated"] = false
	session.Save(r, w)
}

func IsAuthenticated(r *http.Request) bool {
	session, err := store.Get(r, getSessionName())
	if err != nil {
		return false
	}

	authenticated, ok := session.Values["authenticated"].(bool)
	return ok && authenticated
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/login" || r.URL.Path == "/favicon.ico" ||
			(len(r.URL.Path) >= 8 && r.URL.Path[:8] == "/static/") {
			next.ServeHTTP(w, r)
			return
		}

		if r.URL.Path == "/" && r.Method == "GET" {
			next.ServeHTTP(w, r)
			return
		}

		if !IsAuthenticated(r) {
			if strings.HasPrefix(r.URL.Path, "/api/") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
