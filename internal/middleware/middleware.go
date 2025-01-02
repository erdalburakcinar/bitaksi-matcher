package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type JWTMiddleware struct {
	SecretKey string
}

func NewJWTMiddleware(secretKey string) *JWTMiddleware {
	return &JWTMiddleware{SecretKey: secretKey}
}

func (m *JWTMiddleware) ValidateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error": "Missing Authorization header"}`, http.StatusUnauthorized)
			return
		}

		// Extract token from "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		fmt.Println("tokenString: ", tokenString)

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.NewValidationError("Unexpected signing method", jwt.ValidationErrorSignatureInvalid)
			}
			return []byte(m.SecretKey), nil
		})

		fmt.Println("token: ", token)
		if err != nil || !token.Valid {
			http.Error(w, `{"error": "Invalid token"}`, http.StatusUnauthorized)
			return
		}

		_, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, `{"error": "Invalid claims"}`, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}