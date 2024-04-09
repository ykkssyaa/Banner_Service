package middleware

import (
	"BannerService/internal/consts"
	"net/http"
)

func CheckToken(token string, allowedTokens ...string) bool {
	for _, allowedToken := range allowedTokens {
		if token == allowedToken {
			return true
		}
	}
	return false
}

func TokenMiddleware(allowedTokens ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("token")

			print(token)
			if token == "" {
				http.Error(w, consts.ErrorUnauthorized, http.StatusUnauthorized)
				return

			} else if !CheckToken(token, allowedTokens...) {
				http.Error(w, consts.ErrorForbidden, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
