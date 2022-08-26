package ferry

import "net/http"

// DefaultErrorHandler knows how to work with ClientError and converts other types of error to ServerError
var DefaultErrorHandler ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
	clientErr, ok := err.(ClientError)
	if ok {
		_ = Encode(w, r, clientErr.Code, clientErr)
	} else {
		_ = Encode(w, r, http.StatusInternalServerError, ServerError{"internal server error"})
	}
}

// ErrorHandler can be used to handle custom errors. Default is DefaultErrorHandler.
type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)

// ServerError is used for server's error responses.
type ServerError struct {
	Message string `json:"error"`
}

// ClientError can be used for client's error responses.
type ClientError struct {
	Code    int    `json:"-"`
	Message string `json:"error"`
}

func (e ServerError) Error() string { return e.Message }

func (e ClientError) Error() string { return e.Message }
