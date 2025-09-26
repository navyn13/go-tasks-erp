package utils

import (
	"context"
	"net/http"
)

type contextKey string

func SetContext(r *http.Request, claims map[string]interface{}) *http.Request {
	ctx := r.Context()

	if idVal, ok := claims["id"].(float64); ok {
		ctx = context.WithValue(ctx, "id", int(idVal))
	}

	if username, ok := claims["username"].(string); ok {
		ctx = context.WithValue(ctx, "username", username)
	}

	if role, ok := claims["role"].(string); ok {
		ctx = context.WithValue(ctx, "role", role)
	}

	return r.WithContext(ctx)
}
