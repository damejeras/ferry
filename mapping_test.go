package ferry

import (
	"reflect"
	"testing"
)

type exoticJSONRequest struct {
	Number uint8   `json:"number"`
	Float  float64 `json:"float"`
	Data   []byte  `json:"data"`
}

type exoticQueryRequest struct {
	Number int8    `query:"number"`
	Float  float64 `query:"float"`
}

type invalidQueryRequest struct {
	Data []byte `query:"data"`
}

func TestJSONMapping(t *testing.T) {
	testCases := []struct {
		name     string
		input    interface{}
		expected map[string]string
	}{
		{
			name:  "exotic",
			input: exoticJSONRequest{},
			expected: map[string]string{
				"number": "integer",
				"float":  "float",
				"data":   "binary",
			},
		},
		{
			name:     "empty",
			input:    empty{},
			expected: map[string]string{},
		},
		{
			name:  "string",
			input: testPayload{},
			expected: map[string]string{
				"value": "string",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mapping, err := jsonMapping(testCase.input)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if !reflect.DeepEqual(mapping, testCase.expected) {
				t.Errorf("unexpected mapping, got %v", mapping)
			}
		})
	}
}

func TestQueryMapping(t *testing.T) {
	t.Run("test mapping", func(t *testing.T) {
		testCases := []struct {
			name     string
			input    interface{}
			expected map[string]string
		}{
			{
				name:     "json",
				input:    exoticJSONRequest{},
				expected: map[string]string{},
			},
			{
				name:     "empty",
				input:    empty{},
				expected: map[string]string{},
			},
			{
				name:  "query",
				input: queryRequest{},
				expected: map[string]string{
					"value": "string",
				},
			},
			{
				name:  "exotic",
				input: exoticQueryRequest{},
				expected: map[string]string{
					"number": "integer",
					"float":  "float",
				},
			},
		}

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				mapping, err := queryMapping(testCase.input)
				if err != nil {
					t.Errorf("unexpected error: %s", err)
				}

				if !reflect.DeepEqual(mapping, testCase.expected) {
					t.Errorf("unexpected mapping, got %v", mapping)
				}
			})
		}
	})

	t.Run("test invalid", func(t *testing.T) {
		_, err := queryMapping(invalidQueryRequest{})
		if err == nil {
			t.Errorf("expected error")
		}
	})
}
