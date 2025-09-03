package resp

import (
	"bytes"
	"testing"
)

func TestWriteValue(t *testing.T) {
	tests := []struct {
		name     string
		value    Value
		expected string
		wantErr  bool
	}{
		{
			name:     "simple string",
			value:    Value{Typ: "simple", Str: "OK"},
			expected: "+OK\r\n",
			wantErr:  false,
		},
		{
			name:     "error",
			value:    Value{Typ: "error", Str: "Error message"},
			expected: "-Error message\r\n",
			wantErr:  false,
		},
		{
			name:     "integer",
			value:    Value{Typ: "integer", Num: 42},
			expected: ":42\r\n",
			wantErr:  false,
		},
		{
			name:     "bulk string",
			value:    Value{Typ: "bulk", Str: "hello"},
			expected: "$5\r\nhello\r\n",
			wantErr:  false,
		},
		{
			name:     "empty bulk string",
			value:    Value{Typ: "bulk", Str: ""},
			expected: "$0\r\n\r\n",
			wantErr:  false,
		},
		{
			name:     "bulk null",
			value:    Value{Typ: "null", NullTyp: "bulk"},
			expected: "$-1\r\n",
			wantErr:  false,
		},
		{
			name:     "array null",
			value:    Value{Typ: "null", NullTyp: "array"},
			expected: "*-1\r\n",
			wantErr:  false,
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
			wantErr:  false,
		},
		{
			name:     "unknown type",
			value:    Value{Typ: "unknown"},
			expected: "-ERR unknown type\r\n",
			wantErr:  false,
		},
		{
			name:     "unknown null type",
			value:    Value{Typ: "null", NullTyp: "invalid"},
			expected: "-ERR unknown null type\r\n",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := WriteValue(&buf, tt.value)

			if (err != nil) != tt.wantErr {
				t.Errorf("WriteValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if buf.String() != tt.expected {
				t.Errorf("WriteValue() = %q, want %q", buf.String(), tt.expected)
			}
		})
	}
}
