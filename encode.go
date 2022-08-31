package ferry

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Encode encodes payload to JSON and writes it to http.ResponseWriter along with all required headers.
func Encode(w http.ResponseWriter, r *http.Request, status int, payload any) error {
	body, err := json.Marshal(payload)
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

	if _, err := out.Write(body); err != nil {
		return fmt.Errorf("write body: %w", err)
	}

	return nil
}
