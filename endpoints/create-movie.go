package endpoints

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"stream-platform/converters"
	"stream-platform/customerrors"
	"stream-platform/models/database"
	"stream-platform/models/dto"
	"stream-platform/repositories"
	"stream-platform/utils"
	"strings"
)

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	createMovieDto, errorMsg, httpStatus := parseAndValidateCreateRequest(r)
	if errorMsg != "" {
		http.Error(w, errorMsg, httpStatus)
		return
	}

	loggedInUserId := utils.ExtractUserIdFromHeader(r)
	loggedInUserEmail := utils.ExtractUserEmailFromHeader(r)

	// Insert the movie in the DB
	dbMovie := database.Movie{Title: createMovieDto.Title, Release_Date: createMovieDto.Release_Date,
		Genre: createMovieDto.Genre, Synopsis: createMovieDto.Synopsis, CreateUser_Id: loggedInUserId}
	insertedId, insertErr := repositories.InsertMovie(&dbMovie, &createMovieDto.CastMemberIds)
	if insertErr != nil {
		slog.Error("Unable to create the movie", insertErr)

		var errorMsg string
		if strings.Contains(insertErr.Error(), "1062 (23000)") {
			errorMsg = customerrors.ErrMovieAlreadyExists.Error()
		} else {
			errorMsg = customerrors.ErrMovieUnexpectedError.Error()
		}

		http.Error(w, errorMsg, http.StatusInternalServerError)
		return
	}

	dbCastMembers, err := repositories.FindAllCastMembers(createMovieDto.CastMemberIds)
	if err != nil {
		slog.Error("Error while fetching cast members for created movie with", "id", insertedId)
	}

	createUserDto := dto.UserDto{Id: loggedInUserId, Email: loggedInUserEmail}
	castMembersDtos := converters.FromCastMemberDbsToCastMemberDtos(dbCastMembers)
	movieDto := dto.MovieDto{Id: insertedId, Title: dbMovie.Title, Release_Date: dbMovie.Release_Date,
		Genre: dbMovie.Genre, Synopsis: dbMovie.Synopsis, CreateUser: createUserDto,
		CastMembers: castMembersDtos}

	slog.Info("Movie successfully created.", "Id", movieDto.Id, "Title", movieDto.Title)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movieDto)
}

func parseAndValidateCreateRequest(r *http.Request) (*dto.CreateMovieRequestDto, string, int) {
	var createMovieDto dto.CreateMovieRequestDto
	var errorMsg string
	var httpStatus int

	// Parse request body
	err := json.NewDecoder(r.Body).Decode(&createMovieDto)
	if err != nil {
		slog.Error("Unable to decode body", err)
		errorMsg = customerrors.ErrInvalidRequest.Error()
		httpStatus = http.StatusBadRequest
	} else {
		// Validate main DTO properties
		err := utils.Validate[dto.CreateMovieRequestDto](&createMovieDto)
		if err != nil {
			slog.Error("Validation failed on given movie", err)
			errorMsg = fmt.Sprintf(customerrors.ErrInvalidFields.Error(), err.Error())
			httpStatus = http.StatusBadRequest
		} else {
			// Validate cast members
			numDbCastMembers, err := repositories.CountCastMembersByIds(createMovieDto.CastMemberIds)

			if err != nil {
				slog.Error("Error while retreving cast members", err)
				errorMsg = customerrors.ErrMovieUnexpectedError.Error()
				httpStatus = http.StatusInternalServerError
			} else if numDbCastMembers != len(createMovieDto.CastMemberIds) {
				errorMsg = customerrors.ErrMovieInvalidCastMembers.Error()
				httpStatus = http.StatusInternalServerError
			}
		}
	}

	return &createMovieDto, errorMsg, httpStatus
}
