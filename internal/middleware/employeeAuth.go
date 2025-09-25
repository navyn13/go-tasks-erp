package middleware

import (
	"fmt"
	"net/http"

	"context"

	"github.com/golang-jwt/jwt/v5"
)

func EmployeeOnly(next http.Handler) http.Handler {
	fmt.Println("EmployeeOnly triggered")
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
			http.Error(w, "Only Employees can access", http.StatusForbidden)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok {
			role, ok := claims["role"].(string) // get role as string
			if !ok {
				fmt.Println("Role not found")
				return
			}
			if role != "employee" {
				http.Error(w, "Only Employees can access", http.StatusForbidden)
				return
			}
		}
		ctx := context.WithValue(r.Context(), "role", claims["role"])
		username, ok := claims["username"].(string)
		if ok {
			ctx = context.WithValue(ctx, "username", username)
		}
		fmt.Println("context has been set with username:", username)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
