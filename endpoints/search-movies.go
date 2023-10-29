package endpoints

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"stream-platform/converters"
	"stream-platform/customerrors"
	"stream-platform/models/dto"
	"stream-platform/models/filter"
	"stream-platform/repositories"
	"time"
)

func SearchMovies(w http.ResponseWriter, r *http.Request) {
	filters := parseFilters(r)

	res, err := repositories.GetMovies(filters)
	if err != nil {
		slog.Error("Unexpected error while retrieving movies", err)
		http.Error(w, customerrors.ErrMovieUnexpectedError.Error(), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		dtos := make([]dto.MovieDto, 0)
		for _, tuple := range *res {
			dto := converters.FromMovieDbToMovieDto(&tuple.Movie, &tuple.User, &tuple.CastMembers)
			dtos = append(dtos, *dto)
		}

		json.NewEncoder(w).Encode(dtos)
	}
}

func parseFilters(r *http.Request) *filter.SearchMoviesFilter {
	queryParams := r.URL.Query()

	// Extract filters from URL query parameters
	filters := filter.SearchMoviesFilter{
		Title: queryParams.Get("title"),
		Genre: queryParams.Get("genre"),
	}

	releaseDateFilterValue := queryParams.Get("releaseDate")
	if releaseDateFilterValue != "" {
		parsedValue, err := time.Parse("2006-01-02T15:04:05.000Z", releaseDateFilterValue)
		if err != nil {
			slog.Error("Error while parsing release date filter. Ignoring releaseDate value.", err)
		} else {
			filters.ReleaseDate = parsedValue
		}
	}

	return &filters
}
