package ferry

import (
	"fmt"
	"reflect"

	"github.com/fatih/structtag"
	"github.com/mitchellh/reflectwalk"
)

// jsonMapping walks over the target struct and returns a map of json tags to their corresponding go types.
func jsonMapping(v interface{}) (map[string]string, error) {
	mapping := make(bodyMap)

	if err := reflectwalk.Walk(v, mapping); err != nil {
		return nil, err
	}

	return mapping, nil
}

// queryMapping walks over the target struct and returns a map of query tags to their corresponding go types.
func queryMapping(v interface{}) (map[string]string, error) {
	mapping := make(queryMap)

	if err := reflectwalk.Walk(v, mapping); err != nil {
		return nil, err
	}

	return mapping, nil
}

type bodyMap map[string]string
type queryMap map[string]string

func (s bodyMap) Struct(value reflect.Value) error  { return nil }
func (s queryMap) Struct(value reflect.Value) error { return nil }

func (s bodyMap) StructField(field reflect.StructField, value reflect.Value) error {
	tags, err := structtag.Parse(string(field.Tag))
	if err != nil {
		return err
	}

	jsonTag, err := tags.Get("json")
	if err != nil {
		// skip fields without json tag
		return nil
	}

	switch field.Type.Kind() {
	case reflect.String:
		s[jsonTag.Name] = "string"
	case reflect.Bool:
		s[jsonTag.Name] = "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s[jsonTag.Name] = "integer"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s[jsonTag.Name] = "integer"
	case reflect.Float32, reflect.Float64:
		s[jsonTag.Name] = "float"
	case reflect.Slice:
		if field.Type.Elem().Kind() == reflect.Uint8 {
			s[jsonTag.Name] = "binary"
		} else {
			s[jsonTag.Name] = "array"
		}
	case reflect.Array:
		s[jsonTag.Name] = "array"
	case reflect.Struct:
		if field.Type.Name() == "Time" {
			s[jsonTag.Name] = "date-time"
		} else {
			s[jsonTag.Name] = "object"
		}
	case reflect.Map, reflect.Interface:
		s[jsonTag.Name] = "object"
	case reflect.Ptr:
		if field.Type.Elem().Kind() == reflect.Struct {
			s[jsonTag.Name] = "object"
		}
	default:
		return fmt.Errorf("kind %q is not supported as json param", field.Type.Kind())
	}

	return nil
}

func (s queryMap) StructField(field reflect.StructField, value reflect.Value) error {
	tags, err := structtag.Parse(string(field.Tag))
	if err != nil {
		return err
	}

	queryTag, err := tags.Get("query")
	if err != nil {
		// skip fields without query tag
		return nil
	}

	kind := field.Type.Kind()
	switch kind {
	case reflect.String:
		s[queryTag.Name] = "string"
	case reflect.Bool:
		s[queryTag.Name] = "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s[queryTag.Name] = "integer"
	case reflect.Float32, reflect.Float64:
		s[queryTag.Name] = "float"
	default:
		return fmt.Errorf("kind %q is not supported as query param", kind)
	}

	return nil
}
