package endpoints

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"stream-platform/converters"
	"stream-platform/customerrors"
	"stream-platform/models/dto"
	"stream-platform/repositories"
	"stream-platform/utils"
	"strings"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var requestDto dto.CreateUserRequestDto
	err := json.NewDecoder(r.Body).Decode(&requestDto)
	if err != nil {
		slog.Error("Unable to decode body")
		http.Error(w, customerrors.ErrInvalidRequest.Error(), http.StatusBadRequest)
		return
	}

	// Validate DTO
	err = utils.Validate[dto.CreateUserRequestDto](&requestDto)
	if err != nil {
		slog.Error("Validation failed on given user", err)
		http.Error(w, fmt.Sprintf(customerrors.ErrInvalidFields.Error(), err.Error()), http.StatusBadRequest)
		return
	}

	// Hash the password to avoid to store it in clear
	hashedPwd, err := utils.HashPassword(requestDto.Password)
	if err != nil {
		slog.Error("Unable to hash user specified password")
		http.Error(w, customerrors.ErrUserInvalidPassword.Error(), http.StatusBadRequest)
		return
	}

	// Insert the user in the DB
	dbUser := converters.FromFieldsToUserDb(requestDto.Email, hashedPwd)
	insertedId, insertErr := repositories.InsertUser(&dbUser)
	if insertErr != nil {
		slog.Error("Unable to create user", insertErr)

		var errorMsg string
		if strings.Contains(insertErr.Error(), "1062 (23000)") {
			errorMsg = customerrors.ErrUserAlreadyExists.Error()
		} else {
			errorMsg = customerrors.ErrUserUnexpectedError.Error()
		}

		http.Error(w, errorMsg, http.StatusInternalServerError)
		return
	}

	userDto := converters.FromFieldsToUserDto(insertedId, dbUser.Email)
	slog.Info("User successfully created.", "Id", userDto.Id, "Email", userDto.Email)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userDto)
}
