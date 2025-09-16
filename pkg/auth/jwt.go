package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWT claims structure
type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

// Генерация токена
func GenerateToken(userID int64) (string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	if len(secret) == 0 {
		secret = []byte("supersecretkey") // fallback
	}

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// Middleware для проверки токена
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1]
		secret := []byte(os.Getenv("JWT_SECRET"))
		if len(secret) == 0 {
			secret = []byte("supersecretkey") // fallback
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Добавляем userID в контекст запроса
		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Вспомогательная функция для получения userID из контекста
func GetUserIDFromContext(ctx context.Context) (int64, error) {
	uid, ok := ctx.Value("userID").(int64)
	if !ok {
		return 0, fmt.Errorf("userID not found in context")
	}
	return uid, nil
}
