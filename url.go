package ferry

import (
	"fmt"
	"strings"
)

type url struct {
	path    string
	service string
	method  string
}

func parseURL(path string) (*url, error) {
	p := new(url)
	parts := strings.Split(path, ".")
	if len(parts) == 2 {
		p.method = parts[1]

		servicePath := strings.Split(strings.Trim(parts[0], "/"), "/")
		switch len(servicePath) {
		case 1:
			p.service = servicePath[0]
		default:
			p.service = servicePath[len(servicePath)-1]
			p.path = "/" + (strings.Trim(path, "/"))
		}

		return p, nil
	}

	return nil, fmt.Errorf("%q is not valid path, endpoint must have /path/Service.Method signature", path)
}
