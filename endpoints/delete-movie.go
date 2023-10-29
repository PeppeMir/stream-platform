package endpoints

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"stream-platform/customerrors"
	"stream-platform/repositories"
	"stream-platform/utils"

	"github.com/gorilla/mux"
)

func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	requestVariables := mux.Vars(r)

	id, err := strconv.ParseInt(requestVariables["id"], 10, 64)
	if err != nil {
		slog.Error("Invalid 'id' parameter", err)
		http.Error(w, fmt.Errorf(customerrors.ErrInvalidParameter.Error(), "id").Error(), http.StatusBadRequest)
		return
	}

	loggedInUserId := utils.ExtractUserIdFromHeader(r)

	// Delete the movie in the DB
	rowsAffected, err := repositories.DeleteMovie(id, loggedInUserId)
	if err != nil {
		slog.Error("Error deleting the movie", err)
		http.Error(w, customerrors.ErrMovieUnexpectedError.Error(), http.StatusInternalServerError)
	} else if rowsAffected <= 0 {
		slog.Error("Unable to delete: movie not owned by the user", err)
		http.Error(w, customerrors.ErrMovieDeleteNotOwner.Error(), http.StatusInternalServerError)
	} else {
		slog.Info("Movie successfully deleted.", "id", id)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	}
}
