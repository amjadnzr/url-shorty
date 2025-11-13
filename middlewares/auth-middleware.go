package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/amjadnzr/url-shortly/helpers"
)

func RequiresAuth(nextHanlder http.Handler, tokenHelper *helpers.TokenHelper) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]

		claims, err := tokenHelper.ValidateJWTToken(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}

		ctx := context.WithValue(r.Context(), "claims", claims)
		nextHanlder.ServeHTTP(w, r.WithContext(ctx))
	})
}
