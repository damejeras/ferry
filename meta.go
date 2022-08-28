package ferry

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// serviceMeta contains meta information for the service method.
type serviceMeta struct {
	reflectedName string
	methodName    string
	body          map[string]string
	query         map[string]string
}

// path builds path for the service method.
func (s serviceMeta) path() string { return "/" + s.methodName }

// buildMeta uses reflection to determine service name and method name.
func buildMeta(function, request interface{}) (serviceMeta, error) {
	name := runtime.FuncForPC(reflect.ValueOf(function).Pointer()).Name()

	// -fm is a suffix for functions that have receiver.
	nameParts := strings.Split(strings.TrimSuffix(name, "-fm"), ".")
	if len(nameParts) < 2 {
		return serviceMeta{}, fmt.Errorf("can not use %q as handler", name)
	}

	meta := serviceMeta{
		reflectedName: name,
		methodName:    nameParts[len(nameParts)-1],
	}

	var err error
	if meta.body, err = jsonMapping(request); err != nil {
		return serviceMeta{}, fmt.Errorf("can not create json mapping: %w", err)
	}

	if meta.query, err = queryMapping(request); err != nil {
		return serviceMeta{}, fmt.Errorf("can not create query mapping: %w", err)
	}

	return meta, nil
}
