package ferry

import (
	"net/http/httptest"
	"testing"
)

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

func TestRouter(t *testing.T) {
	t.Run("creates api spec", func(t *testing.T) {
		t.Parallel()
		svc := testService{}
		router := NewRouter(WithServiceDiscovery)

		router.Register(Procedure(svc.TestProcedureWithParams))
		router.Register(Stream(svc.StreamOneEvent))

		rr := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		router.ServeHTTP(rr, r)
		content := rr.Body.Bytes()
		expected := `[
  {
    "method": "POST",
    "path": "http://example.com/testService.TestProcedureWithParams",
    "body": {
      "value": "string"
    }
  },
  {
    "method": "GET",
    "path": "http://example.com/testService.StreamOneEvent",
    "query": {
      "value": "string"
    }
  }
]`
		if string(content) != expected {
			t.Errorf("unexpected response, got %s", content)
		}
	})
}
