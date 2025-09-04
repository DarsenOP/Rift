package server

import (
	"strconv"
	"testing"

	"github.com/DarsenOP/Rift/internal/resp"
	"github.com/DarsenOP/Rift/internal/storage"
)

func TestHandleCommand(t *testing.T) {
	store := storage.New()

	for i := 1; i <= 10; i++ {
		key := "k" + strconv.Itoa(i)
		val := "v" + strconv.Itoa(i)
		store.Set(key, val)
	}

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
			name: "SET OK 1",
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
			name: "SET OK 2",
			input: resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "bulk", Str: "SET"},
					{Typ: "bulk", Str: "mykey2"},
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
			expected: resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'set' with expiration"},
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
		{
			name: "DEL one existing",
			input: resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "bulk", Str: "DEL"},
					{Typ: "bulk", Str: "mykey"},
				},
			},
			expected: resp.Value{Typ: "integer", Num: 1},
		},
		{
			name: "DEL two, one missing",
			input: resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "bulk", Str: "DEL"},
					{Typ: "bulk", Str: "mykey2"},
					{Typ: "bulk", Str: "nosuch"},
				},
			},
			expected: resp.Value{Typ: "integer", Num: 1},
		},
		{
			name: "DEL multiple all missing",
			input: resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "bulk", Str: "DEL"},
					{Typ: "bulk", Str: "nosuch1"},
					{Typ: "bulk", Str: "nosuch2"},
				},
			},
			expected: resp.Value{Typ: "integer", Num: 0},
		},
		{
			name: "DEL no args",
			input: resp.Value{
				Typ:   "array",
				Array: []resp.Value{{Typ: "bulk", Str: "DEL"}},
			},
			expected: resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'del' command"},
		},
		{
			name: "DEL non-bulk key",
			input: resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "bulk", Str: "DEL"},
					{Typ: "integer", Num: 123},
				},
			},
			expected: resp.Value{Typ: "error", Str: "ERR arguments should be bulk strings"},
		},
		{
			name: "EXISTS one present",
			input: resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "bulk", Str: "EXISTS"},
					{Typ: "bulk", Str: "k3"},
				},
			},
			expected: resp.Value{Typ: "integer", Num: 1},
		},
		{
			name: "EXISTS two, one missing",
			input: resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "bulk", Str: "EXISTS"},
					{Typ: "bulk", Str: "k5"},
					{Typ: "bulk", Str: "nosuch"},
				},
			},
			expected: resp.Value{Typ: "integer", Num: 1},
		},
		{
			name: "EXISTS all missing",
			input: resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "bulk", Str: "EXISTS"},
					{Typ: "bulk", Str: "nosuch1"},
					{Typ: "bulk", Str: "nosuch2"},
				},
			},
			expected: resp.Value{Typ: "integer", Num: 0},
		},
		{
			name: "EXISTS no arguments",
			input: resp.Value{
				Typ:   "array",
				Array: []resp.Value{{Typ: "bulk", Str: "EXISTS"}},
			},
			expected: resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'exists' command"},
		},
		{
			name: "EXISTS non-bulk argument",
			input: resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "bulk", Str: "EXISTS"},
					{Typ: "integer", Num: 123},
				},
			},
			expected: resp.Value{Typ: "error", Str: "ERR arguments should be bulk strings"},
		},
		{
			name: "SET with EX seconds",
			input: resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "bulk", Str: "SET"},
					{Typ: "bulk", Str: "k"},
					{Typ: "bulk", Str: "v"},
					{Typ: "bulk", Str: "EX"},
					{Typ: "bulk", Str: "2"},
				},
			},
			expected: resp.Value{Typ: "simple", Str: "OK"},
		},
		{
			name: "SET with PX milliseconds",
			input: resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "bulk", Str: "SET"},
					{Typ: "bulk", Str: "k"},
					{Typ: "bulk", Str: "v"},
					{Typ: "bulk", Str: "PX"},
					{Typ: "bulk", Str: "2000"},
				},
			},
			expected: resp.Value{Typ: "simple", Str: "OK"},
		},
		{
			name: "SET invalid flag",
			input: resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "bulk", Str: "SET"},
					{Typ: "bulk", Str: "k"},
					{Typ: "bulk", Str: "v"},
					{Typ: "bulk", Str: "XX"},
					{Typ: "bulk", Str: "10"},
				},
			},
			expected: resp.Value{Typ: "error", Str: "ERR unsupported option"},
		},
		{
			name: "SET negative TTL",
			input: resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "bulk", Str: "SET"},
					{Typ: "bulk", Str: "k"},
					{Typ: "bulk", Str: "v"},
					{Typ: "bulk", Str: "EX"},
					{Typ: "bulk", Str: "-1"},
				},
			},
			expected: resp.Value{Typ: "error", Str: "ERR value is not an integer or out of range"},
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

func TestTTLExpire(t *testing.T) {
	store := storage.New()
	defer store.Shutdown()

	// seed a key
	store.Set("k", "v")

	tests := []struct {
		name     string
		input    resp.Value
		expected resp.Value
	}{
		{
			name: "TTL on existing key without expiry",
			input: resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "bulk", Str: "TTL"},
					{Typ: "bulk", Str: "k"},
				},
			},
			expected: resp.Value{Typ: "integer", Num: -1},
		},
		{
			name: "TTL on missing key",
			input: resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "bulk", Str: "TTL"},
					{Typ: "bulk", Str: "nosuch"},
				},
			},
			expected: resp.Value{Typ: "integer", Num: -2},
		},
		{
			name: "EXPIRE set 5 s",
			input: resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "bulk", Str: "EXPIRE"},
					{Typ: "bulk", Str: "k"},
					{Typ: "bulk", Str: "5"},
				},
			},
			expected: resp.Value{Typ: "integer", Num: 1},
		},
		{
			name: "EXPIRE on missing key",
			input: resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "bulk", Str: "EXPIRE"},
					{Typ: "bulk", Str: "nosuch"},
					{Typ: "bulk", Str: "10"},
				},
			},
			expected: resp.Value{Typ: "integer", Num: 0},
		},
		{
			name: "EXPIRE wrong arity",
			input: resp.Value{
				Typ:   "array",
				Array: []resp.Value{{Typ: "bulk", Str: "EXPIRE"}},
			},
			expected: resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'expire' command"},
		},
		{
			name: "EXPIRE non-integer seconds",
			input: resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "bulk", Str: "EXPIRE"},
					{Typ: "bulk", Str: "k"},
					{Typ: "bulk", Str: "abc"},
				},
			},
			expected: resp.Value{Typ: "error", Str: "ERR value is not an integer or out of range"},
		},
		{
			name: "TTL wrong arity",
			input: resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "bulk", Str: "TTL"},
					{Typ: "bulk", Str: "k1"},
					{Typ: "bulk", Str: "k2"},
				},
			},
			expected: resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'ttl' command"},
		},
		{
			name: "TTL non-bulk key",
			input: resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "bulk", Str: "TTL"},
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
