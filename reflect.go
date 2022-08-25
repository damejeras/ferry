package ferry

func requestJSONReflection[R any](r *R) (map[string]string, error) {
	mapping := make(map[string]string)
	// TODO: parse request params
	return mapping, nil
}

func requestQueryReflection[R any](r *R) (map[string]string, error) {
	mapping := make(map[string]string)
	// TODO: parse query params
	return mapping, nil
}
