package ferry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"

	"github.com/fatih/structtag"
	"github.com/mitchellh/reflectwalk"
)

type bodyWalk map[string]string
type queryWalk map[string]string

func (s bodyWalk) Struct(value reflect.Value) error  { return nil }
func (s queryWalk) Struct(value reflect.Value) error { return nil }

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
			return fmt.Errorf("unsupported kind parameter %q in request", kind)
		}
	}
	return nil
}

func (s bodyWalk) StructField(field reflect.StructField, value reflect.Value) error {
	tags, err := structtag.Parse(string(field.Tag))
	if err != nil {
		return err
	}

	jsonTag, err := tags.Get("json")
	if err != nil {
		// skip fields without json tag
		return nil
	}

	kind := field.Type.Kind()

	if kind == reflect.Invalid {
		return nil
	}

	if kind == reflect.String {
		s[jsonTag.Name] = "string"
	}

	if kind == reflect.Bool {
		s[jsonTag.Name] = "boolean"
	}

	if kind >= reflect.Int && kind <= reflect.Uint64 {
		s[jsonTag.Name] = "integer"
	}

	if kind == reflect.Float32 || kind == reflect.Float64 {
		s[jsonTag.Name] = "double"
	}

	if kind == reflect.Array || kind == reflect.Slice {
		s[jsonTag.Name] = "array"
	}

	if kind == reflect.Pointer || kind == reflect.Struct || kind == reflect.Map {
		s[jsonTag.Name] = "object"
	}

	return nil
}

func (s queryWalk) StructField(field reflect.StructField, value reflect.Value) error {
	tags, err := structtag.Parse(string(field.Tag))
	if err != nil {
		return err
	}

	queryTag, err := tags.Get("query")
	if err != nil {
		// skip fields without json tag
		return nil
	}

	kind := field.Type.Kind()

	if kind == reflect.Invalid {
		return nil
	}

	if kind == reflect.String {
		s[queryTag.Name] = "string"
	}

	if kind == reflect.Bool {
		s[queryTag.Name] = "boolean"
	}

	if kind >= reflect.Int && kind <= reflect.Uint64 {
		s[queryTag.Name] = "integer"
	}

	if kind == reflect.Float32 || kind == reflect.Float64 {
		s[queryTag.Name] = "double"
	}

	if kind == reflect.Array || kind == reflect.Slice {
		s[queryTag.Name] = "array"
	}

	if kind == reflect.Pointer || kind == reflect.Struct || kind == reflect.Map {
		s[queryTag.Name] = "object"
	}

	return nil
}

func jsonMapping(v interface{}) (map[string]string, error) {
	mapping := make(bodyWalk)

	if err := reflectwalk.Walk(v, mapping); err != nil {
		return nil, err
	}

	return mapping, nil
}

func queryMapping(v interface{}) (map[string]string, error) {
	mapping := make(queryWalk)

	if err := reflectwalk.Walk(v, mapping); err != nil {
		return nil, err
	}

	return mapping, nil
}
