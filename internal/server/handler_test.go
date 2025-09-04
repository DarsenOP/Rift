package server

import (
	"testing"

	"github.com/DarsenOP/Rift/internal/resp"
	"github.com/DarsenOP/Rift/internal/storage"
)

func TestHandleCommand(t *testing.T) {
	store := storage.New()

	tests := []struct {
		name     string
		input    resp.Value
		expected resp.Value
	}{
		{
			name:     "empty array",
			input:    resp.Value{Typ: "array", Array: []resp.Value{}},
			expected: resp.Value{Typ: "array", Array: []resp.Value{}},
		},
		{
			name: "PING without args",
			input: resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "bulk", Str: "PING"},
				},
			},
			expected: resp.Value{Typ: "simple", Str: "PONG"},
		},
		{
			name: "PING with string arg",
			input: resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "bulk", Str: "PING"},
					{Typ: "bulk", Str: "hello"},
				},
			},
			expected: resp.Value{Typ: "bulk", Str: "hello"},
		},
		{
			name: "PING with integer arg",
			input: resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "bulk", Str: "PING"},
					{Typ: "integer", Num: 42},
				},
			},
			expected: resp.Value{Typ: "integer", Num: 42},
		},
		{
			name: "PING too many args",
			input: resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "bulk", Str: "PING"},
					{Typ: "bulk", Str: "a"},
					{Typ: "bulk", Str: "b"},
				},
			},
			expected: resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'ping' command"},
		},
		{
			name: "COMMAND",
			input: resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "bulk", Str: "COMMAND"},
				},
			},
			expected: resp.Value{Typ: "array", Array: []resp.Value{}},
		},
		{
			name: "unknown command",
			input: resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "bulk", Str: "BLAH"},
				},
			},
			expected: resp.Value{Typ: "error", Str: "ERR unknown command 'BLAH'"},
		},
		{
			name: "non-bulk command",
			input: resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "simple", Str: "PING"},
				},
			},
			expected: resp.Value{Typ: "error", Str: "ERR command must be a bulk string"},
		},
		{
			name: "SET OK",
			input: resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "bulk", Str: "SET"},
					{Typ: "bulk", Str: "mykey"},
					{Typ: "bulk", Str: "myvalue"},
				},
			},
			expected: resp.Value{Typ: "simple", Str: "OK"},
		},
		{
			name: "SET wrong arity (0 args)",
			input: resp.Value{
				Typ:   "array",
				Array: []resp.Value{{Typ: "bulk", Str: "SET"}},
			},
			expected: resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'set' command"},
		},
		{
			name: "SET wrong arity (3 args)",
			input: resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "bulk", Str: "SET"},
					{Typ: "bulk", Str: "k"},
					{Typ: "bulk", Str: "v"},
					{Typ: "bulk", Str: "extra"},
				},
			},
			expected: resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'set' command"},
		},
		{
			name: "SET non-bulk key",
			input: resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "bulk", Str: "SET"},
					{Typ: "integer", Num: 123},
					{Typ: "bulk", Str: "v"},
				},
			},
			expected: resp.Value{Typ: "error", Str: "ERR arguments should be bulk strings"},
		},
		{
			name: "GET missing key",
			input: resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "bulk", Str: "GET"},
					{Typ: "bulk", Str: "nosuch"},
				},
			},
			expected: resp.Value{Typ: "null", NullTyp: "bulk"}, // $-1\r\n
		},
		{
			name: "GET existing key",
			input: resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "bulk", Str: "GET"},
					{Typ: "bulk", Str: "mykey"},
				},
			},
			expected: resp.Value{Typ: "bulk", Str: "myvalue"},
		},
		{
			name: "GET wrong arity",
			input: resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "bulk", Str: "GET"},
					{Typ: "bulk", Str: "k1"},
					{Typ: "bulk", Str: "k2"},
				},
			},
			expected: resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'get' command"},
		},
		{
			name: "GET non-bulk key",
			input: resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "bulk", Str: "GET"},
					{Typ: "integer", Num: 123},
				},
			},
			expected: resp.Value{Typ: "error", Str: "ERR argument should be a bulk string"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HandleCommand(store, tt.input)

			if !valuesEqual(result, tt.expected) {
				t.Errorf("HandleCommand() = %+v, want %+v", result, tt.expected)
			}
		})
	}
}

// Helper function from parser_test.go
func valuesEqual(a, b resp.Value) bool {
	if a.Typ != b.Typ {
		return false
	}

	switch a.Typ {
	case "simple", "error", "bulk":
		return a.Str == b.Str
	case "integer":
		return a.Num == b.Num
	case "null":
		return a.NullTyp == b.NullTyp
	case "array":
		if len(a.Array) != len(b.Array) {
			return false
		}
		for i := range a.Array {
			if !valuesEqual(a.Array[i], b.Array[i]) {
				return false
			}
		}
		return true
	default:
		return false
	}
}
