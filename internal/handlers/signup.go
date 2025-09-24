package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
)

type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// DB connection function
func getDB() (*sql.DB, error) {
	// Update user, password, host, dbname accordingly
	dsn := "root:baf75918@tcp(127.0.0.1:3306)/go_tasks_erp"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func signup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Signup endpoint hit")
	var req SignupRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" || req.Role == "" {
		http.Error(w, "username, password, and role are required", http.StatusBadRequest)
		return
	}
	password_hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	// Connect to DB
	db, err := getDB()
	if err != nil {
		http.Error(w, "DB connection error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Insert user into users table
	_, err = db.Exec(
		"INSERT INTO users (username, password, role) VALUES (?, ?, ?)",
		req.Username, password_hash, req.Role,
	)
	if err != nil {
		http.Error(w, "DB insert error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := "signup got triggered and user saved to DB"
	fmt.Println(response)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
