package endpoints

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"stream-platform/converters"
	"stream-platform/customerrors"
	"stream-platform/repositories"

	"github.com/gorilla/mux"
)

func GetMovie(w http.ResponseWriter, r *http.Request) {
	requestVariables := mux.Vars(r)

	id, err := strconv.ParseInt(requestVariables["id"], 10, 64)
	if err != nil {
		slog.Error("Invalid 'id' parameter", err)
		http.Error(w, fmt.Errorf(customerrors.ErrInvalidParameter.Error(), "id").Error(), http.StatusBadRequest)
		return
	}

	dbMovie, dbUser, dbCastMembers, err := repositories.GetMovie(id)
	if err != nil {
		if errors.Is(err, customerrors.ErrMovieNotFound) {
			slog.Error("No movie found with", "id", id, err)
			http.Error(w, customerrors.ErrMovieNotFound.Error(), http.StatusNotFound)
		} else {
			slog.Error("Unexpected error while retrieving movie with ", "id", id, err)
			http.Error(w, customerrors.ErrMovieUnexpectedError.Error(), http.StatusInternalServerError)
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(converters.FromMovieDbToMovieDto(dbMovie, dbUser, dbCastMembers))
	}
}
