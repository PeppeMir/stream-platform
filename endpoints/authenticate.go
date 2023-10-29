package endpoints

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"stream-platform/customerrors"
	"stream-platform/models/dto"
	"stream-platform/repositories"
	"stream-platform/utils"
)

func Authenticate(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	if len(email) == 0 || len(password) == 0 {
		http.Error(w, customerrors.ErrUserMissingEmailPassword.Error(), http.StatusBadRequest)
		return
	}

	// Find the user with the given email in the DB
	user, err := repositories.FindUserByEmail(email)
	if err != nil {
		http.Error(w, customerrors.ErrUserNotFound.Error(), http.StatusNotFound)
		return
	}

	// Check if the given password matches the hash stored in the DB
	if !utils.PasswordMatches(password, user.Password) {
		http.Error(w, customerrors.ErrUserWrongEmailPassword.Error(), http.StatusUnauthorized)
		return
	}

	// Generate new token and return it as request response
	token, err := utils.GenerateToken(user.Id, user.Email)

	if err != nil {
		slog.Error("Error generating JWT token", err)
		http.Error(w, customerrors.ErrUserUnexpectedError.Error(), http.StatusInternalServerError)
	} else {
		slog.Info("Token successfully generated for user", "id", user.Id, "email", user.Email)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(dto.AuthResponseDto{Token: token})
	}

}
