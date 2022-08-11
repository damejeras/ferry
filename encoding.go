package ferry

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func EncodeJSON(w http.ResponseWriter, r *http.Request, status int, payload any) error {
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

func DecodeJSON(r *http.Request, v any) error {
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

func DecodeQuery(r *http.Request, v any) error {
	targetType := reflect.TypeOf(v)
	if targetType.Kind() != reflect.Ptr || reflect.ValueOf(v).Elem().Kind() != reflect.Struct {
		return errors.New("target should be pointer to a struct")
	}

	targetValue := reflect.ValueOf(v)

	for i := 0; i < targetType.Elem().NumField(); i++ {
		field := targetType.Elem().Field(i)
		key := field.Tag.Get("query")
		kind := field.Type.Kind()
		// Get the value from query params with given key
		val := r.URL.Query().Get(key)
		//  Get reference of field value provided to input `d`
		result := targetValue.Elem().Field(i)

		switch kind {
		case reflect.String:
			result.SetString(val)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			num, err := strconv.Atoi(val)
			if err != nil {
				return ClientError{
					Code:    http.StatusBadRequest,
					Message: fmt.Sprintf("can not convert %q to int", val),
				}
			}
			result.SetInt(int64(num))
		case reflect.Float32:
			num, err := strconv.ParseFloat(val, 32)
			if err != nil {
				return ClientError{
					Code:    http.StatusBadRequest,
					Message: fmt.Sprintf("can not convert %q to float32", val),
				}
			}
			result.SetFloat(num)
		case reflect.Float64:
			num, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return ClientError{
					Code:    http.StatusBadRequest,
					Message: fmt.Sprintf("can not convert %q to float64", val),
				}
			}
			result.SetFloat(num)
		case reflect.Bool:
			boolean, err := strconv.ParseBool(val)
			if err != nil {
				return ClientError{
					Code:    http.StatusBadRequest,
					Message: fmt.Sprintf("can not convert %q to bool", val),
				}
			}
			result.SetBool(boolean)
		default:
			return fmt.Errorf("unsupported kind parameter %q in request", kind)
		}
	}
	return nil
}
