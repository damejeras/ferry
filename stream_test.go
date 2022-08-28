package ferry

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func (s testService) EmptyStreamForSixSeconds(ctx context.Context, _ *empty) (<-chan Event[empty], error) {
	c := make(chan Event[empty])

	go func() {
		defer close(c)

		select {
		case <-ctx.Done():
			return
		case <-time.After(6 * time.Second):
			return
		}
	}()

	return c, nil
}

func (s testService) StreamForSixSeconds(ctx context.Context, r *queryRequest) (<-chan Event[testPayload], error) {
	c := make(chan Event[testPayload])

	go func() {
		defer close(c)
		ctx, cancel := context.WithTimeout(ctx, time.Second*6)
		defer cancel()

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		var id int
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				id++
				c <- Event[testPayload]{
					ID:      strconv.Itoa(id),
					Payload: &testPayload{Value: r.Value},
				}
			}
		}
	}()

	return c, nil
}

func (s testService) StreamOneEvent(ctx context.Context, r *queryRequest) (<-chan Event[testPayload], error) {
	c := make(chan Event[testPayload])

	go func() {
		c <- Event[testPayload]{
			ID:      "1",
			Payload: &testPayload{Value: r.Value},
		}
		close(c)
	}()

	return c, nil
}

func (s testService) LeakyStream(ctx context.Context, r *empty) (<-chan Event[empty], error) {
	c := make(chan Event[empty])

	go func() {
		<-ctx.Done()
	}()

	return c, nil
}

func TestStream(t *testing.T) {
	t.Run("keep alive messages are sent each 5 seconds", func(t *testing.T) {
		t.Parallel()
		router := NewRouter()
		svc := testService{}
		router.Register(Stream(svc.EmptyStreamForSixSeconds))
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, httptest.NewRequest("GET", "/EmptyStreamForSixSeconds", nil))
		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		content, err := io.ReadAll(rr.Body)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		expected := `event: keep-alive

event: keep-alive

`

		if string(content) != expected {
			t.Errorf("unexpected response, got %s", content)
		}
	})

	t.Run("keep alive is not sent when there is activity", func(t *testing.T) {
		t.Parallel()
		router := NewRouter()
		svc := testService{}
		router.Register(Stream(svc.StreamForSixSeconds))
		rr := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/StreamForSixSeconds?value=test_data", nil)
		// ctx, cancel := context.WithTimeout(context.Background(), time.Second*6)
		// defer cancel()
		// r = r.WithContext(ctx)

		router.ServeHTTP(rr, r)
		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		content, err := io.ReadAll(rr.Body)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		expected := `event: keep-alive

id: 1
event: testPayload
data: {"value":"test_data"}

id: 2
event: testPayload
data: {"value":"test_data"}

id: 3
event: testPayload
data: {"value":"test_data"}

id: 4
event: testPayload
data: {"value":"test_data"}

id: 5
event: testPayload
data: {"value":"test_data"}

id: 6
event: testPayload
data: {"value":"test_data"}

`

		if string(content) != expected {
			t.Errorf("unexpected response, got %s", content)
		}
	})

	t.Run("panics on leaking stream channel", func(t *testing.T) {
		t.Parallel()
		router := NewRouter()
		svc := testService{}
		router.Register(Stream(svc.LeakyStream))
		rr := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/LeakyStream", nil)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r = r.WithContext(ctx)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		router.ServeHTTP(rr, r)
	})

	t.Run("receive one event", func(t *testing.T) {
		t.Parallel()
		router := NewRouter()
		svc := testService{}
		router.Register(Stream(svc.StreamOneEvent))
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, httptest.NewRequest("GET", "/StreamOneEvent?value=test", nil))
		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		content, err := io.ReadAll(rr.Body)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		expected := `event: keep-alive

id: 1
event: testPayload
data: {"value":"test"}

`

		if string(content) != expected {
			t.Errorf("unexpected response, got %s", content)
		}

		if rr.Header().Get("Content-Type") != "text/event-stream" {
			t.Error("expected Content-Type: text/event-stream")
		}

		if rr.Header().Get("Cache-Control") != "no-cache" {
			t.Error("expected Cache-Control: no-cache")
		}

		if rr.Header().Get("Connection") != "keep-alive" {
			t.Error("expected Connection: keep-alive")
		}
	})
}
