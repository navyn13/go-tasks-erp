package middleware

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("my_secret_key")

func AdminOnly(next http.Handler) http.Handler {
	fmt.Println("AdminOnly triggered")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authTokenString := r.Header.Get("authTokenString")
		if authTokenString == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}
		token, err := jwt.Parse(authTokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			fmt.Println("Invalid token:", err)
			http.Error(w, "Only Admins can access", http.StatusForbidden)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok {
			role, ok := claims["role"].(string) // get role as string
			if !ok {
				fmt.Println("Role not found")
				return
			}
			fmt.Println("User role:", role)
			if role != "admin" {
				http.Error(w, "Only Admins can access", http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)

	})
}
