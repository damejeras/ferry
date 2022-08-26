package ferry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
)

// decodeJSON decodes *http.Request into target struct.
// Request must have "Content-Type" header set to "application/json".
func decodeJSON[T any](r *http.Request, v *T) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return ClientError{
			Code:    http.StatusUnsupportedMediaType,
			Message: "application/json content-type expected",
		}
	}

	if err := json.NewDecoder(io.LimitReader(r.Body, 1024*1024)).Decode(v); err != nil {
		if err == io.EOF {
			return ClientError{
				Code:    http.StatusBadRequest,
				Message: "empty request body",
			}
		}
		return fmt.Errorf("decode request body: %w", err)
	}

	return r.Body.Close()
}

// decodeQuery decodes query values from http.Request into target struct.
// This function maps r.URL.Query values to struct properties by `query` tag.
func decodeQuery[T any](r *http.Request, v *T) error {
	targetType := reflect.TypeOf(v)
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
			return fmt.Errorf("unsupported kind %q in request query", kind)
		}
	}
	return nil
}
