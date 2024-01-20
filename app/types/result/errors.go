package result

import "errors"

var (
	ErrInvalidInputParam    = errors.New("invalid input parameters")
	ErrInvalidAuthorization = errors.New("invalid authorization parameters")
	ErrDatabaseOperation    = errors.New("operate database error")
	ErrNoAuthorization      = errors.New("no authorization params in the header")
)
