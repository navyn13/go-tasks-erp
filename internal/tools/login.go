package tools

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"github.com/navyn13/go-tasks-erp/internal/db"
	"github.com/navyn13/go-tasks-erp/internal/models/usersSchema"
	"github.com/navyn13/go-tasks-erp/internal/utils"
)

var jwtKey = []byte(os.Getenv("JWTKEY"))

func Login(w http.ResponseWriter, r *http.Request) {
	var req usersSchema.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "username and password are required", http.StatusBadRequest)
		return
	}

	db, err := db.GetDB()
	if err != nil {
		http.Error(w, "DB connection error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	//check for username and password check in user table in db
	var storedPassword, role string
	var userID int
	err = db.QueryRow("SELECT password, role, id FROM users WHERE username = ? LIMIT 1", req.Username).Scan(&storedPassword, &role, &userID)
	if err != nil {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}

	if isPasswordMatched := utils.CompareHashPassword(storedPassword, req.Password); !isPasswordMatched {
		http.Error(w, "Password did not match", http.StatusForbidden)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": req.Username,
		"role":     role,
		"id":       userID,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // expires in 24h
	})
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "could not generate token", http.StatusInternalServerError)
		return
	}

	response := usersSchema.LoginResponse{
		Token: tokenString,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
