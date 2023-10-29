package customerrors

import "errors"

var ErrMovieNotFound = errors.New("movie not found")
var ErrMovieAlreadyExists = errors.New("a movie already exists with the given title")
var ErrMovieInvalidCastMembers = errors.New("some of the specified cast members do not exist")
var ErrMovieUpdateNotOwner = errors.New("cannot update the movie because the logged-in user is not the owner")
var ErrMovieDeleteNotOwner = errors.New("cannot delete the movie because the logged-in user is not the owner")
var ErrMovieUnexpectedError = errors.New("an unexpected error occurred while creating the movie")
