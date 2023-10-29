package customerrors

import "errors"

var ErrCannotConnectToDB = errors.New("unable to connect to SQL database")
var ErrInvalidRequest = errors.New("the request is invalid")
var ErrInvalidFields = errors.New("some fields are invalid: %s")
var ErrInvalidParameter = errors.New("invalid '%s' parameter")
