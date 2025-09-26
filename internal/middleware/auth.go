package middleware

import (
	"net/http"
	"strings"

	"github.com/navyn13/go-tasks-erp/internal/utils"
)

func Auth(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		endpoint := r.URL.Path
		authTokenString := utils.GetHeader(r, "authTokenString")
		if authTokenString == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		claims, err := utils.ParseJWTClaims(authTokenString)
		if err != nil {
			http.Error(w, "Invalid Map Claims from jwt", http.StatusUnauthorized)
			return
		}
		role, ok := claims["role"].(string)
		if !ok {
			http.Error(w, "Role not found in token", http.StatusForbidden)
			return
		}
		if strings.HasPrefix(endpoint, "/admin") && role != "admin" {
			http.Error(w, "Only admins can access", http.StatusForbidden)
			return
		}
		if strings.HasPrefix(endpoint, "/employee") && role != "employee" {
			http.Error(w, "Only admins can access", http.StatusForbidden)
			return
		}

		r = utils.SetContext(r, claims)

		next.ServeHTTP(w, r)

	})
}
