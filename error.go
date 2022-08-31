package ferry

import "net/http"

// ErrorHandler handles errors in Handler methods. Can be assigned to router with WithErrorHandler option.
type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)

// ClientError encapsulates error with HTTP status code. Can be used to return error to client.
type ClientError struct {
	Code    int    `json:"-"`
	Message string `json:"error"`
}

func (e ClientError) Error() string { return e.Message }

// DefaultErrorHandler knows how to encode ClientError.
// In case of unexpected error it returns 500 Internal Server Error.
var DefaultErrorHandler ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
	switch tErr := err.(type) {
	case ClientError:
		Encode(w, r, tErr.Code, tErr)
	default:
		Encode(w, r, http.StatusInternalServerError, serverError{"internal server error"})
	}
}

type serverError struct {
	Message string `json:"error"`
}
