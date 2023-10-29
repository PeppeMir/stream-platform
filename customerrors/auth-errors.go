package customerrors

import "errors"

var ErrMissingAuthHeader = errors.New("missing Authorization Header")
var ErrInvalidJwtToken = errors.New("error validating JWT token: %s")
