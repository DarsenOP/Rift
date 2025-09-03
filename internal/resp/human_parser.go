package resp

import (
	"bufio"
	"strings"
)

func ParseHuman(reader *bufio.Reader) (Value, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return Value{}, err
	}

	line = strings.TrimSpace(line)
	if line == "" {
		return Value{Typ: "array", Array: []Value{}}, nil
	}

	// Split into command and arguments
	parts := strings.Fields(line)
	array := make([]Value, len(parts))
	for i, part := range parts {
		array[i] = Value{Typ: "bulk", Str: part}
	}

	return Value{Typ: "array", Array: array}, nil
}
