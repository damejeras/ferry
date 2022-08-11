package ferry

import "net/http"

type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)

type ClientError struct {
	Code    int    `json:"-"`
	Message string `json:"error"`
}

func (e ClientError) Error() string { return e.Message }

var DefaultErrorHandler ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
	switch err.(type) {
	case ClientError:
		_ = Encode(w, r, err.(ClientError).Code, err.(ClientError))
	default:
		_ = Encode(w, r, http.StatusInternalServerError, map[string]string{"message": "internal server error"})
	}
}
