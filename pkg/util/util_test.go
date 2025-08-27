package util

import "testing"

func TestReverseString(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{"basic", "hello", "olleh"},
		{"empty", "", ""},
		{"unicode", "Hello, 世界", "界世 ,olleH"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReverseString(tt.args); got != tt.want {
				t.Errorf("ReverseString() = %v, want %v", got, tt.want)
			}
		})
	}
}
