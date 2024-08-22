package middleware

import (
	"context"
	"github.com/antibomberman/mego-api/pkg/response"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

var jwtKey = []byte("secret")

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Fail(w, "Invalid authorization header")
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			response.Fail(w, "Token not provided")
			return
		}

		tokenString := bearerToken[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			response.Fail(w, "Invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Fail(w, "Invalid token claims")
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", claims["user_id"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
