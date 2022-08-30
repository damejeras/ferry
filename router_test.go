package ferry

type testService struct{}

type testPayload struct {
	Value string `json:"value"`
}

type queryRequest struct {
	Value string `query:"value"`
}

type jsonRequest struct {
	Value string `json:"value"`
}

type empty struct{}
