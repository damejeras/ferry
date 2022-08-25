package ferry

import (
	"reflect"

	"github.com/fatih/structtag"
	"github.com/mitchellh/reflectwalk"
)

type bodyWalk map[string]string

func (s bodyWalk) Struct(value reflect.Value) error {
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

type queryWalk map[string]string

func (s queryWalk) Struct(value reflect.Value) error {
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

func reflectBody(v interface{}) (map[string]string, error) {
	mapping := make(bodyWalk)

	if err := reflectwalk.Walk(v, mapping); err != nil {
		return nil, err
	}

	return mapping, nil
}

func reflectQuery(v interface{}) (map[string]string, error) {
	mapping := make(queryWalk)

	if err := reflectwalk.Walk(v, mapping); err != nil {
		return nil, err
	}

	return mapping, nil
}
