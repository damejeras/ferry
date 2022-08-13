package ferry

import (
	"fmt"
	"strings"
)

// MetaKey is context value key for Meta.
type MetaKey struct{}

// Meta contains information about endpoint.
type Meta struct {
	Path    string
	Service string
	Method  string
}

func toMeta(path string) (Meta, error) {
	meta := Meta{}
	parts := strings.Split(path, ".")
	if len(parts) == 2 {
		meta.Method = parts[1]

		servicePath := strings.Split(strings.Trim(parts[0], "/"), "/")
		switch len(servicePath) {
		case 1:
			meta.Service = servicePath[0]
		default:
			meta.Service = servicePath[len(servicePath)-1]
			meta.Path = "/" + (strings.Trim(path, "/"))
		}

		return meta, nil
	}

	return Meta{}, fmt.Errorf("%q is not valid path, endpoint must have /prefix/Service.Method signature", path)
}
