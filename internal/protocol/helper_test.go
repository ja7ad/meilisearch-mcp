package protocol

import (
	"reflect"
	"testing"

	"github.com/meilisearch/meilisearch-go"
)

func TestSwapIndexes(t *testing.T) {
	cases := []struct {
		name     string
		input    []interface{}
		expected []*meilisearch.SwapIndexesParams
	}{
		{
			name: "two valid pairs",
			input: []interface{}{
				[]interface{}{"a", "b"},
				[]interface{}{"c", "d"},
			},
			expected: []*meilisearch.SwapIndexesParams{
				{Indexes: []string{"a", "b"}},
				{Indexes: []string{"c", "d"}},
			},
		},
		{
			name: "skip invalid length + non slice",
			input: []interface{}{
				[]interface{}{"onlyOne"}, // invalid length
				"notASlice",              // not a slice
				[]interface{}{"x", "y"},  // valid
			},
			expected: []*meilisearch.SwapIndexesParams{
				{Indexes: []string{"x", "y"}},
			},
		},
		{
			name: "mixed types inside pair",
			input: []interface{}{
				[]interface{}{"x", 123},   // second not string -> becomes empty
				[]interface{}{456, "bar"}, // first not string -> becomes empty
			},
			expected: []*meilisearch.SwapIndexesParams{
				{Indexes: []string{"x", ""}},
				{Indexes: []string{"", "bar"}},
			},
		},
		{
			name:  "nil element skipped",
			input: []interface{}{nil, []interface{}{"a", "b"}},
			expected: []*meilisearch.SwapIndexesParams{
				{Indexes: []string{"a", "b"}},
			},
		},
	}

	for _, tc := range cases {
		c := tc
		t.Run(c.name, func(t *testing.T) {
			got := swapIndexes(c.input)
			if !reflect.DeepEqual(got, c.expected) {
				// Simple diff output; rely on %#v formatting
				t.Fatalf("unexpected result.\nexpected: %#v\n     got: %#v", c.expected, got)
			}
		})
	}
}
