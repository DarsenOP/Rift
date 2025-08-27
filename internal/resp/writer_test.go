package resp

import (
	"bytes"
	"testing"
)

func TestWriter(t *testing.T) {
	tests := []struct {
		name     string
		value    Value
		expected string
	}{
		{
			name:     "simple string",
			value:    Value{Typ: "simple", Str: "OK"},
			expected: "+OK\r\n",
		},
		{
			name:     "error",
			value:    Value{Typ: "error", Str: "Error message"},
			expected: "-Error message\r\n",
		},
		{
			name:     "integer",
			value:    Value{Typ: "integer", Num: 42},
			expected: ":42\r\n",
		},
		{
			name:     "bulk string",
			value:    Value{Typ: "bulk", Str: "hello"},
			expected: "$5\r\nhello\r\n",
		},
		{
			name:     "empty bulk string",
			value:    Value{Typ: "bulk", Str: ""},
			expected: "$0\r\n\r\n",
		},
		{
			name:     "null",
			value:    Value{Typ: "null", NullTyp: "bulk"},
			expected: "$-1\r\n",
		},
		{
			name: "array",
			value: Value{
				Typ: "array",
				Array: []Value{
					{Typ: "bulk", Str: "PING"},
				},
			},
			expected: "*1\r\n$4\r\nPING\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := WriteValue(&buf, tt.value)
			if err != nil {
				t.Errorf("WriteValue() error = %v", err)
				return
			}

			if buf.String() != tt.expected {
				t.Errorf("WriteValue() = %q, want %q", buf.String(), tt.expected)
			}
		})
	}
}
