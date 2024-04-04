package transport

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

func parseToken(authorizationHeader string) (string, error) {
	fields := strings.Fields(authorizationHeader)
	if len(fields) == 2 && strings.EqualFold(fields[0], "bearer") {
		return fields[1], nil
	}

	return "", fmt.Errorf("not a bearer authorization")
}

func WithAuth(jwtSecret string) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/admin/auth" {
				handler.ServeHTTP(w, r)
				return
			}

			token, err := parseToken(r.Header.Get("Authorization"))
			if err != nil {
				_ = WriteError(w, http.StatusUnauthorized, err.Error())
				return
			}

			if _, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			}); err != nil {
				_ = WriteError(w, http.StatusUnauthorized, err.Error())
				return
			}

			handler.ServeHTTP(w, r)
		})
	}
}
