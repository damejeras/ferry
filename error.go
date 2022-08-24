package ferry

import "net/http"

type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)

type ServerError struct {
	Message string `json:"error"`
}

func (e ServerError) Error() string { return e.Message }

type ClientError struct {
	Code    int    `json:"-"`
	Message string `json:"error"`
}

func (e ClientError) Error() string { return e.Message }

var DefaultErrorHandler ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
	clientErr, ok := err.(ClientError)
	if ok {
		_ = EncodeJSON(w, r, clientErr.Code, clientErr)
	} else {
		_ = EncodeJSON(w, r, http.StatusInternalServerError, ServerError{"internal server error"})
	}
}
