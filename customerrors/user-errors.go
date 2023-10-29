package customerrors

import "errors"

var ErrUserNotFound = errors.New("user not found")
var ErrUserMissingEmailPassword = errors.New("please provide both email and password to obtain a token")
var ErrUserWrongEmailPassword = errors.New("wrong email and/or password")
var ErrUserInvalidPassword = errors.New("the given password is invalid")
var ErrUserAlreadyExists = errors.New("a user already exists with the given email")
var ErrUserUnexpectedError = errors.New("an unexpected error occurred while creating the user")
