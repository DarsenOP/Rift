package resp

import (
	"bufio"
	"errors"
)

type Value struct {
	Typ   string
	Str   string
	Array []Value
}

var ErrInvalidSyntax = errors.New("invalid syntax")

func Parse(reader *bufio.Reader) (Value, error) {
	// Start with simple string parsing: "+OK\r\n"
	firstByte, err := reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch firstByte {
	case '+':
		return parseSimpleString(reader)
	default:
		return Value{}, ErrInvalidSyntax
	}
	// Then implement arrays: "*1\r\n$4\r\nPING\r\n"
}

func parseSimpleString(reader *bufio.Reader) (Value, error) {
	// Read until we find \r\n
	line, err := readLine(reader)
	if err != nil {
		return Value{}, err
	}

	return Value{
		Typ: "simple",
		Str: string(line),
	}, nil
}

// readLine reads until \r\n and returns the line without the terminator
func readLine(reader *bufio.Reader) ([]byte, error) {
	var line []byte

	for {
		b, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}

		line = append(line, b)

		// Check if we've reached \r\n
		if len(line) >= 2 && line[len(line)-2] == '\r' && line[len(line)-1] == '\n' {
			// Return the line without the \r\n
			return line[:len(line)-2], nil
		}
	}
}
