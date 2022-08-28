package ferry

import (
	"context"
	"reflect"
	"testing"
)

type testMeta struct{}

func (*testMeta) TestProcedure(ctx context.Context, _ *empty) (*empty, error) {
	return &empty{}, nil
}

func testProc(ctx context.Context, _ *empty) (*empty, error) {
	return &empty{}, nil
}

func TestBuildMeta(t *testing.T) {
	t.Run("handles anonymous functions", func(t *testing.T) {
		m, err := buildMeta(func(ctx context.Context, _ *empty) (*empty, error) {
			return &empty{}, nil
		}, empty{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		expected := serviceMeta{
			reflectedName: "github.com/damejeras/ferry.TestBuildMeta.func1.1",
			methodName:    "1",
			body:          make(map[string]string),
			query:         make(map[string]string),
		}

		if !reflect.DeepEqual(m, expected) {
			t.Errorf("got %+v", m)
		}

		m, err = buildMeta(func(ctx context.Context, _ *empty) (*empty, error) {
			return &empty{}, nil
		}, empty{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		expected = serviceMeta{
			reflectedName: "github.com/damejeras/ferry.TestBuildMeta.func1.2",
			methodName:    "2",
			body:          make(map[string]string),
			query:         make(map[string]string),
		}

		if !reflect.DeepEqual(m, expected) {
			t.Errorf("got %+v", m)
		}
	})

	t.Run("handles function with pointer receiver", func(t *testing.T) {
		svc := new(testMeta)
		m, err := buildMeta(svc.TestProcedure, empty{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		expected := serviceMeta{
			reflectedName: "github.com/damejeras/ferry.(*testMeta).TestProcedure-fm",
			methodName:    "TestProcedure",
			body:          make(map[string]string),
			query:         make(map[string]string),
		}

		if !reflect.DeepEqual(m, expected) {
			t.Errorf("got %+v", m)
		}
	})

	t.Run("handles function with value receiver", func(t *testing.T) {
		svc := testService{}
		m, err := buildMeta(svc.StreamOneEvent, empty{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		expected := serviceMeta{
			reflectedName: "github.com/damejeras/ferry.testService.StreamOneEvent-fm",
			methodName:    "StreamOneEvent",
			body:          make(map[string]string),
			query:         make(map[string]string),
		}

		if !reflect.DeepEqual(m, expected) {
			t.Errorf("got %+v", m)
		}
	})

	t.Run("handles function without receiver", func(t *testing.T) {
		m, err := buildMeta(testProc, empty{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		expected := serviceMeta{
			reflectedName: "github.com/damejeras/ferry.testProc",
			methodName:    "testProc",
			body:          make(map[string]string),
			query:         make(map[string]string),
		}

		if !reflect.DeepEqual(m, expected) {
			t.Errorf("got %+v", m)
		}
	})
}
