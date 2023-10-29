package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"stream-platform/customerrors"
	"stream-platform/utils"
	"strings"
)

func AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		if len(tokenString) == 0 {
			http.Error(w, customerrors.ErrMissingAuthHeader.Error(), http.StatusUnauthorized)
			return
		}

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, fmt.Sprintf(customerrors.ErrInvalidJwtToken.Error(), err.Error()), http.StatusUnauthorized)
			return
		}

		slog.Info("Claims read from token", "id", claims.UserId, "email", claims.UserEmail)

		r.Header.Set("userId", strconv.FormatInt(claims.UserId, 10))
		r.Header.Set("userEmail", claims.UserEmail)

		next.ServeHTTP(w, r)
	})
}
