package middleware

import (
	"BannerService/internal/consts"
	"context"
	"net/http"
)

func ParseRole(token string) string {
	if token == consts.UserToken {
		return consts.UserRole
	} else if token == consts.AdminToken {
		return consts.AdminRole
	}
	return ""
}

func CheckRoles(token string, allowedTokens ...string) bool {
	for _, allowedToken := range allowedTokens {
		if token == allowedToken {
			return true
		}
	}
	return false
}

func TokenMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("token")

			if token == "" {
				http.Error(w, "", http.StatusUnauthorized)
				return

			} else {
				role := ParseRole(token)

				if !CheckRoles(role, allowedRoles...) {
					http.Error(w, "", http.StatusForbidden)
					return
				}

				ctx := context.WithValue(r.Context(), "role", role)
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		})
	}
}
