package utils

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	UserId    int64
	UserEmail string
	jwt.RegisteredClaims
}

func GenerateToken(id int64, email string) (string, error) {
	secret := secret()

	slog.Info("Generating token for", "id", id, "email", email)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JwtClaims{
		id,
		email,
		jwt.RegisteredClaims{
			Issuer:    "streamplatform",
			ExpiresAt: &jwt.NumericDate{time.Now().Add(time.Hour * 24)},
		},
	})

	tokenStr, err := token.SignedString(secret)
	if err != nil {
		slog.Info("Generated", "token", tokenStr)
	}

	return tokenStr, err
}

func ValidateToken(tokenStr string) (*JwtClaims, error) {
	slog.Info("Verifying", "token", tokenStr)

	token, _ := jwt.ParseWithClaims(tokenStr, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secret(), nil
	})

	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		slog.Info("Token is valid for", "id", claims.UserId, "email", claims.UserEmail)
		return claims, nil
	} else {
		slog.Info("Token is invalid")
		return nil, errors.New("token is invalid")
	}
}

func ExtractUserIdFromHeader(r *http.Request) int64 {
	id, _ := strconv.ParseInt(r.Header.Get("userId"), 10, 64)
	return id
}

func ExtractUserEmailFromHeader(r *http.Request) string {
	return r.Header.Get("userEmail")
}

func secret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}
