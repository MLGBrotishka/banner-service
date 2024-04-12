package router

import (
	"context"
	"net/http"
)

type AccessCheckFunc func(*http.Request) bool

func AuthMiddleware(check AccessCheckFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("token")
			if token == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			if !check(r) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

type Role string
type ContextKey string

const (
	RoleKey   ContextKey = "role"
	AdminRole Role       = "admin"
	UserRole  Role       = "user"
)

//todo: nice validation

func userAccessCheck(r *http.Request) bool {
	token := r.Header.Get("token")
	if token != "user_token" {
		return false
	}
	ctx := r.Context()
	*r = *r.WithContext(context.WithValue(ctx, RoleKey, UserRole))
	role, ok := r.Context().Value(RoleKey).(Role)
	if !ok || role != UserRole {
		return false
	}
	return true
}

func adminAccessCheck(r *http.Request) bool {
	token := r.Header.Get("token")
	if token != "admin_token" {
		return false
	}
	ctx := r.Context()
	*r = *r.WithContext(context.WithValue(ctx, RoleKey, AdminRole))
	role, ok := r.Context().Value(RoleKey).(Role)
	if !ok || role != AdminRole {
		return false
	}
	return true
}

func userOrAdminAccessCheck(r *http.Request) bool {
	return userAccessCheck(r) || adminAccessCheck(r)
}
