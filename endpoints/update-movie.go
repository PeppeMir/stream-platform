package endpoints

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"stream-platform/converters"
	"stream-platform/customerrors"
	"stream-platform/models/database"
	"stream-platform/models/dto"
	"stream-platform/repositories"
	"stream-platform/utils"
	"strings"
)

func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	updateMovieDto, errorMsg, httpStatus := parseAndValidateUpdateRequest(r)
	if errorMsg != "" {
		http.Error(w, errorMsg, httpStatus)
		return
	}

	oldDbMovie, oldDbUser, oldDbCastMembers, err := repositories.GetMovie(updateMovieDto.Id)
	if err != nil {
		if errors.Is(err, customerrors.ErrMovieNotFound) {
			slog.Error("No movie found with", "id", updateMovieDto.Id, err)
			http.Error(w, customerrors.ErrMovieNotFound.Error(), http.StatusNotFound)
		} else {
			slog.Error("Unexpected error while retrieving movie with ", "id", updateMovieDto.Id, err)
			http.Error(w, customerrors.ErrMovieUnexpectedError.Error(), http.StatusInternalServerError)
		}

		return
	}

	loggedInUserId := utils.ExtractUserIdFromHeader(r)
	if loggedInUserId != oldDbMovie.CreateUser_Id {
		http.Error(w, customerrors.ErrMovieUpdateNotOwner.Error(), http.StatusConflict)
		return
	}

	dbMovie := database.Movie{Id: updateMovieDto.Id, Title: updateMovieDto.Title, Release_Date: updateMovieDto.Release_Date,
		Genre: updateMovieDto.Genre, Synopsis: updateMovieDto.Synopsis, CreateUser_Id: oldDbMovie.CreateUser_Id}

	// Compute which cast members relationships we should create and delete w.r.t the version present in the DB
	castMemberIdToCreate, castMemberIdToDelete := computeCastMembersToCreateAndDelete(oldDbCastMembers, &updateMovieDto.CastMemberIds)

	// Update the movie in the DB
	updateErr := repositories.UpdateMovie(&dbMovie, &castMemberIdToCreate, &castMemberIdToDelete)
	if updateErr != nil {
		slog.Error("Unable to update the movie", updateErr)

		var errorMsg string
		if strings.Contains(updateErr.Error(), "1062 (23000)") {
			errorMsg = customerrors.ErrMovieAlreadyExists.Error()
		} else {
			errorMsg = customerrors.ErrMovieUnexpectedError.Error()
		}

		http.Error(w, errorMsg, http.StatusInternalServerError)
		return
	}

	dbCastMembers, err := repositories.FindAllCastMembers(updateMovieDto.CastMemberIds)
	if err != nil {
		slog.Error("Error while fetching cast members for updated movie with", "id", dbMovie.Id)
	}

	createUserDto := converters.FromUserDbToUserDto(oldDbUser)
	castMembersDtos := converters.FromCastMemberDbsToCastMemberDtos(dbCastMembers)
	movieDto := dto.MovieDto{Id: dbMovie.Id, Title: dbMovie.Title, Release_Date: dbMovie.Release_Date,
		Genre: dbMovie.Genre, Synopsis: dbMovie.Synopsis, CreateUser: createUserDto,
		CastMembers: castMembersDtos}

	slog.Info("Movie successfully updated.", "Id", movieDto.Id, "Title", movieDto.Title)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(movieDto)
}

func computeCastMembersToCreateAndDelete(oldDbCastMembers *[]database.CastMember, newCastMemberIds *[]int64) ([]int64, []int64) {
	oldCastMemberIds := make([]int64, 0)
	castMemberIdToDelete := make([]int64, 0)
	for _, oldCastMember := range *oldDbCastMembers {
		oldCastMemberIds = append(oldCastMemberIds, oldCastMember.Id)
		if !slices.Contains(*newCastMemberIds, oldCastMember.Id) {
			castMemberIdToDelete = append(castMemberIdToDelete, oldCastMember.Id)
		}
	}

	castMemberIdToCreate := make([]int64, 0)
	for _, newCastMemberId := range *newCastMemberIds {
		if !slices.Contains(oldCastMemberIds, newCastMemberId) {
			castMemberIdToCreate = append(castMemberIdToCreate, newCastMemberId)
		}
	}

	return castMemberIdToCreate, castMemberIdToDelete
}

func parseAndValidateUpdateRequest(r *http.Request) (*dto.UpdateMovieRequestDto, string, int) {
	var updateMovieDto dto.UpdateMovieRequestDto
	var errorMsg string
	var httpStatus int

	// Parse request body
	err := json.NewDecoder(r.Body).Decode(&updateMovieDto)
	if err != nil {
		slog.Error("Unable to decode body", err)
		errorMsg = customerrors.ErrInvalidRequest.Error()
		httpStatus = http.StatusBadRequest
	} else {
		// Validate main DTO properties
		err := utils.Validate[dto.UpdateMovieRequestDto](&updateMovieDto)
		if err != nil {
			slog.Error("Validation failed on given movie", err)
			errorMsg = fmt.Sprintf(customerrors.ErrInvalidFields.Error(), err.Error())
			httpStatus = http.StatusBadRequest
		} else {
			// Validate cast members
			numDbCastMembers, err := repositories.CountCastMembersByIds(updateMovieDto.CastMemberIds)

			if err != nil {
				slog.Error("Error while retreving cast members", err)
				errorMsg = customerrors.ErrMovieUnexpectedError.Error()
				httpStatus = http.StatusInternalServerError
			} else if numDbCastMembers != len(updateMovieDto.CastMemberIds) {
				errorMsg = customerrors.ErrMovieInvalidCastMembers.Error()
				httpStatus = http.StatusInternalServerError
			}
		}
	}

	return &updateMovieDto, errorMsg, httpStatus
}
