package resp

import (
	"bufio"
	"strings"
	"testing"
)

func TestParseSimpleString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Value
		wantErr  bool
	}{
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
				if result.Typ != tt.expected.Typ {
					t.Errorf("Parse() type = %v, want %v", result.Typ, tt.expected.Typ)
				}
				if result.Str != tt.expected.Str {
					t.Errorf("Parse() string = %v, want %v", result.Str, tt.expected.Str)
				}
			}
		})
	}
}
