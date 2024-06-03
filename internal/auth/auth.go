package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

const (
	headerAuthorization = "Authorization"
	headerPrefixBearer  = "Bearer "

	tokenExpirationDuration = 24 * time.Hour

	ctxKeyUsername = "username"

	errMissingToken = "error missing token"
	errInvalidToken = "error invalid token"
)

// Config holds the configuration for JWT
type Config struct {
	Secret string
}

// Claims defines the structure of the JWT claims
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var config Config

// Initialize sets the JWT configuration
func Initialize(secret string) {
	config = Config{Secret: secret}
}

// GenerateToken generates a new JWT token
func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(tokenExpirationDuration)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Secret))
}

// ValidateToken validates a JWT token
func ValidateToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	return claims, nil
}

// Middleware is the authentication middleware
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get(headerAuthorization)
		if tokenStr == "" {
			http.Error(w, errMissingToken, http.StatusUnauthorized)
			return
		}

		tokenStr = strings.TrimPrefix(tokenStr, headerPrefixBearer)

		claims, err := ValidateToken(tokenStr)
		if err != nil {
			http.Error(w, errInvalidToken, http.StatusUnauthorized)
			return
		}

		context.Set(r, ctxKeyUsername, claims.Username)
		next.ServeHTTP(w, r)
	})
}
