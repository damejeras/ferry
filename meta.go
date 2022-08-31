package ferry

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// meta contains information about service method.
type meta struct {
	name  string
	body  map[string]string
	query map[string]string
}

// buildMeta uses reflection to determine service name and method name.
func buildMeta(function, request interface{}) (meta, error) {
	name := runtime.FuncForPC(reflect.ValueOf(function).Pointer()).Name()

	// -fm is a suffix for functions that have receiver.
	nameParts := strings.Split(strings.TrimSuffix(name, "-fm"), ".")
	if len(nameParts) < 2 {
		return meta{}, fmt.Errorf("can not use %q as handler", name)
	}

	m := meta{
		name: nameParts[len(nameParts)-1],
	}

	var err error
	if m.body, err = jsonMapping(request); err != nil {
		return meta{}, fmt.Errorf("can not create json mapping: %w", err)
	}

	if m.query, err = queryMapping(request); err != nil {
		return meta{}, fmt.Errorf("can not create query mapping: %w", err)
	}

	return m, nil
}
