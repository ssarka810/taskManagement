package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

const secretkey = "my-secret-key"

func GenerateToken(username string) (string, error) {
	claim := &Claims{
			Username: username,
			StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
					IssuedAt:  time.Now().Unix(),
			},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(secretkey))
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(secretkey), nil
	})
	if err != nil {
			return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			return claims, nil
	}
	return nil, errors.New("invalid token")
}

func AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
					http.Error(w, "Authorization Header Required", http.StatusUnauthorized)
					return
			}
			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) != 2 {
					http.Error(w, "Invalid Token Format", http.StatusUnauthorized)
					return
			}
			claims, err := ValidateToken(bearerToken[1])
			if err != nil {
					logrus.Error("Invalid Token ", err)
					http.Error(w, "Invalid Token", http.StatusUnauthorized)
					return
			}
			ctx := context.WithValue(r.Context(), "username", claims.Username)
			next.ServeHTTP(w, r.WithContext(ctx))
	})
}
