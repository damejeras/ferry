package ferry

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func Encode(w http.ResponseWriter, r *http.Request, status int, payload any) error {
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	var out io.Writer = w
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		gzw := gzip.NewWriter(w)
		out = gzw
		defer gzw.Close()
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if _, err := out.Write(bodyBytes); err != nil {
		return fmt.Errorf("write body: %w", err)
	}

	return nil
}

func Decode(r *http.Request, v any) error {
	if r.Header.Get("Content-type") != "application/json" {
		return ClientError{
			Code:    http.StatusUnsupportedMediaType,
			Message: "application/json content-type expected",
		}
	}

	if err := json.NewDecoder(io.LimitReader(r.Body, 1024*1024)).Decode(v); err != nil {
		return fmt.Errorf("decode request body: %w", err)
	}

	return r.Body.Close()
}
