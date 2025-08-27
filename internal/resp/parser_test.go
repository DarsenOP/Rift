package resp

import (
	"bufio"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Value
		wantErr  bool
	}{
		// Simple Strings
		{
			name:     "valid simple string",
			input:    "+OK\r\n",
			expected: Value{Typ: "simple", Str: "OK"},
			wantErr:  false,
		},
		{
			name:     "empty simple string",
			input:    "+\r\n",
			expected: Value{Typ: "simple", Str: ""},
			wantErr:  false,
		},

		// Simple Errors
		{
			name:     "valid simple error",
			input:    "-Error message\r\n",
			expected: Value{Typ: "error", Str: "Error message"},
			wantErr:  false,
		},
		{
			name:     "empty simple error",
			input:    "-\r\n",
			expected: Value{Typ: "error", Str: ""},
			wantErr:  false,
		},

		// Integers
		{
			name:     "positive integer",
			input:    ":1000\r\n",
			expected: Value{Typ: "integer", Num: 1000},
			wantErr:  false,
		},
		{
			name:     "negative integer",
			input:    ":-42\r\n",
			expected: Value{Typ: "integer", Num: -42},
			wantErr:  false,
		},
		{
			name:     "zero integer",
			input:    ":0\r\n",
			expected: Value{Typ: "integer", Num: 0},
			wantErr:  false,
		},

		// Bulk Strings
		{
			name:     "valid bulk string",
			input:    "$5\r\nhello\r\n",
			expected: Value{Typ: "bulk", Str: "hello"},
			wantErr:  false,
		},
		{
			name:     "empty bulk string",
			input:    "$0\r\n\r\n",
			expected: Value{Typ: "bulk", Str: ""},
			wantErr:  false,
		},
		{
			name:     "null bulk string",
			input:    "$-1\r\n",
			expected: Value{Typ: "null"},
			wantErr:  false,
		},
		{
			name:     "long bulk string",
			input:    "$11\r\nhello world\r\n",
			expected: Value{Typ: "bulk", Str: "hello world"},
			wantErr:  false,
		},

		// Arrays
		{
			name:  "PING command array",
			input: "*1\r\n$4\r\nPING\r\n",
			expected: Value{
				Typ: "array",
				Array: []Value{
					{Typ: "bulk", Str: "PING"},
				},
			},
			wantErr: false,
		},
		{
			name:  "SET command array",
			input: "*3\r\n$3\r\nSET\r\n$5\r\nhello\r\n$5\r\nworld\r\n",
			expected: Value{
				Typ: "array",
				Array: []Value{
					{Typ: "bulk", Str: "SET"},
					{Typ: "bulk", Str: "hello"},
					{Typ: "bulk", Str: "world"},
				},
			},
			wantErr: false,
		},
		{
			name:     "empty array",
			input:    "*0\r\n",
			expected: Value{Typ: "array", Array: []Value{}},
			wantErr:  false,
		},
		{
			name:     "null array",
			input:    "*-1\r\n",
			expected: Value{Typ: "null"},
			wantErr:  false,
		},

		// Nested Arrays
		{
			name:  "nested arrays",
			input: "*2\r\n*2\r\n:1\r\n:2\r\n*2\r\n:3\r\n:4\r\n",
			expected: Value{
				Typ: "array",
				Array: []Value{
					{
						Typ: "array",
						Array: []Value{
							{Typ: "integer", Num: 1},
							{Typ: "integer", Num: 2},
						},
					},
					{
						Typ: "array",
						Array: []Value{
							{Typ: "integer", Num: 3},
							{Typ: "integer", Num: 4},
						},
					},
				},
			},
			wantErr: false,
		},

		// Mixed types in array
		{
			name:  "mixed types array",
			input: "*5\r\n+hello\r\n-err\r\n:42\r\n$5\r\nworld\r\n*2\r\n:1\r\n:2\r\n",
			expected: Value{
				Typ: "array",
				Array: []Value{
					{Typ: "simple", Str: "hello"},
					{Typ: "error", Str: "err"},
					{Typ: "integer", Num: 42},
					{Typ: "bulk", Str: "world"},
					{
						Typ: "array",
						Array: []Value{
							{Typ: "integer", Num: 1},
							{Typ: "integer", Num: 2},
						},
					},
				},
			},
			wantErr: false,
		},

		// Error cases
		{
			name:    "missing terminator",
			input:   "+OK",
			wantErr: true,
		},
		{
			name:    "invalid prefix",
			input:   "OK\r\n",
			wantErr: true,
		},
		{
			name:    "invalid integer",
			input:   ":not_a_number\r\n",
			wantErr: true,
		},
		{
			name:    "invalid bulk string length",
			input:   "$not_a_number\r\n",
			wantErr: true,
		},
		{
			name:    "bulk string missing data",
			input:   "$5\r\nhell\r\n", // Too short
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			result, err := Parse(reader)

			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !valuesEqual(result, tt.expected) {
					t.Errorf("Parse() = %+v, want %+v", result, tt.expected)
				}
			}
		})
	}
}

// valuesEqual compares two Value structs for equality
func valuesEqual(a, b Value) bool {
	if a.Typ != b.Typ {
		return false
	}

	switch a.Typ {
	case "simple", "error", "bulk":
		return a.Str == b.Str
	case "integer":
		return a.Num == b.Num
	case "null":
		return true // Both are null
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
