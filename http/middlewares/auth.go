package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"safa-went/database/models"
	"safa-went/internal/responses"
)

type contextKey string

const UserContextKey contextKey = "user"

func Auth(jwtSecret string, db *gorm.DB) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            publicPaths := map[string]bool{
                "/auth/register": true,
                "/auth/login":    true,
                "/ping":          true,
            }
            if publicPaths[r.URL.Path] || strings.HasPrefix(r.URL.Path, "/swagger") {
                next.ServeHTTP(w, r)
                return
            }
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
                responses.JSONError(w, r, http.StatusUnauthorized, "missing or invalid authorization header")
                return
            }

            tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

            token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
                if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
                    return nil, jwt.ErrSignatureInvalid
                }
                return []byte(jwtSecret), nil
            })
            if err != nil || !token.Valid {
                responses.JSONError(w, r, http.StatusUnauthorized, "invalid or expired token")
                return
            }

            claims, ok := token.Claims.(jwt.MapClaims)
            if !ok {
                responses.JSONError(w, r, http.StatusUnauthorized, "invalid token claims")
                return
            }

            userIDFloat, ok := claims["sub"].(float64)
            if !ok {
                responses.JSONError(w, r, http.StatusUnauthorized, "invalid token subject")
                return
            }

            var user models.User
            if err := db.First(&user, uint(userIDFloat)).Error; err != nil {
                responses.JSONError(w, r, http.StatusUnauthorized, "user not found")
                return
            }

            ctx := context.WithValue(r.Context(), UserContextKey, &user)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}