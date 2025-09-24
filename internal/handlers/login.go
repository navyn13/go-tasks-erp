package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("my_secret_key")

type LoginResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var users = map[string]struct {
	Role string
}{
	"naveen": {Role: "admin"},
	"mayank": {Role: "worker"},
}

func login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "username and password are required", http.StatusBadRequest)
		return
	}
	// Connect to DB
	db, err := getDB()
	if err != nil {
		http.Error(w, "DB connection error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	//check for username and password check in user table in db
	var storedPassword, role string
	err = db.QueryRow("SELECT password, role FROM users WHERE username = ? LIMIT 1", req.Username).Scan(&storedPassword, &role)
	if err != nil {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}

	if req.Password != storedPassword {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": req.Username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // expires in 24h
	})
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "could not generate token", http.StatusInternalServerError)
		return
	}

	resp := LoginResponse{
		Token: tokenString,
		Role:  users[req.Username].Role,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
