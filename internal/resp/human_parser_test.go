package resp

import (
	"bufio"
	"strings"
	"testing"
)

func TestParseHuman(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Value
		wantErr  bool
	}{
		{
			name:  "single command",
			input: "PING\r\n",
			expected: Value{
				Typ: "array",
				Array: []Value{
					{Typ: "bulk", Str: "PING"},
				},
			},
			wantErr: false,
		},
		{
			name:  "command with args",
			input: "SET key value\r\n",
			expected: Value{
				Typ: "array",
				Array: []Value{
					{Typ: "bulk", Str: "SET"},
					{Typ: "bulk", Str: "key"},
					{Typ: "bulk", Str: "value"},
				},
			},
			wantErr: false,
		},
		{
			name:  "multiple spaces",
			input: "SET  key   value  \r\n",
			expected: Value{
				Typ: "array",
				Array: []Value{
					{Typ: "bulk", Str: "SET"},
					{Typ: "bulk", Str: "key"},
					{Typ: "bulk", Str: "value"},
				},
			},
			wantErr: false,
		},
		{
			name:     "empty input",
			input:    "\r\n",
			expected: Value{Typ: "array", Array: []Value{}},
			wantErr:  false,
		},
		{
			name:    "EOF error",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			result, err := ParseHuman(reader)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseHuman() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !valuesEqual(result, tt.expected) {
				t.Errorf("ParseHuman() = %+v, want %+v", result, tt.expected)
			}
		})
	}
}
