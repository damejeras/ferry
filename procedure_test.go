package ferry

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func (t testService) TestProcedureWithoutParams(ctx context.Context, r *empty) (*testPayload, error) {
	return &testPayload{Value: "test_data"}, nil
}

func (t testService) TestProcedureWithParams(ctx context.Context, r *jsonRequest) (*testPayload, error) {
	return &testPayload{Value: r.Value}, nil
}

func TestProcedure(t *testing.T) {
	t.Run("returns response without request params", func(t *testing.T) {
		t.Parallel()
		router := NewRouter()
		svc := testService{}
		router.Register(Procedure(svc.TestProcedureWithoutParams))
		rr := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/TestProcedureWithoutParams", nil)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r = r.WithContext(ctx)

		router.ServeHTTP(rr, r)

		content := rr.Body.Bytes()
		expected := `{"value":"test_data"}`
		if string(content) != expected {
			t.Errorf("unexpected response, got %s", content)
		}
	})

	t.Run("returns response with request params", func(t *testing.T) {
		t.Parallel()
		router := NewRouter()
		svc := testService{}
		router.Register(Procedure(svc.TestProcedureWithParams))
		payload, err := json.Marshal(&jsonRequest{Value: "test_data"})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		rr := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/TestProcedureWithParams", bytes.NewReader(payload))
		r.Header.Set("Content-Type", "application/json")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r = r.WithContext(ctx)

		router.ServeHTTP(rr, r)

		if rr.Code != http.StatusOK {
			t.Errorf("unexpected response code, got %d", rr.Code)
		}

		content := rr.Body.Bytes()
		expected := `{"value":"test_data"}`
		if string(content) != expected {
			t.Errorf("unexpected response, got %s", content)
		}

		if rr.Header().Get("Content-Type") != "application/json; charset=utf-8" {
			t.Errorf("unexpected content type, got %s", rr.Header().Get("Content-Type"))
		}
	})

	t.Run("returns error if request params invalid", func(t *testing.T) {
		t.Parallel()
		router := NewRouter()
		svc := testService{}
		router.Register(Procedure(svc.TestProcedureWithParams))
		rr := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/TestProcedureWithParams", nil)
		r.Header.Set("Content-Type", "application/json")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r = r.WithContext(ctx)

		router.ServeHTTP(rr, r)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("unexpected response code, got %d", rr.Code)
		}
	})
}
