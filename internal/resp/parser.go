package resp

import (
	"bufio"
	"errors"
	"io"
	"strconv"
)

type Value struct {
	Typ     string
	Str     string
	Num     int
	Array   []Value
	NullTyp string
}

var ErrInvalidSyntax = errors.New("invalid syntax")

func Parse(reader *bufio.Reader) (Value, error) {
	firstByte, err := reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch firstByte {
	case '+':
		return parseSimpleString(reader)
	case '-':
		return parseSimpleError(reader)
	case ':':
		return parseInteger(reader)
	case '$':
		return parseBulkString(reader)
	case '*':
		return parseArray(reader)
	default:
		return Value{}, ErrInvalidSyntax
	}
}

func parseSimpleString(reader *bufio.Reader) (Value, error) {
	line, err := readLine(reader)
	if err != nil {
		return Value{}, err
	}

	return Value{
		Typ: "simple",
		Str: string(line),
	}, nil
}

func parseSimpleError(reader *bufio.Reader) (Value, error) {
	line, err := readLine(reader)
	if err != nil {
		return Value{}, err
	}

	return Value{
		Typ: "error",
		Str: string(line),
	}, nil
}

func parseInteger(reader *bufio.Reader) (Value, error) {
	line, err := readLine(reader)
	if err != nil {
		return Value{}, err
	}

	num, err := strconv.Atoi(string(line))
	if err != nil {
		return Value{}, ErrInvalidSyntax
	}

	return Value{
		Typ: "integer",
		Num: num,
	}, nil
}

func parseBulkString(reader *bufio.Reader) (Value, error) {
	lenStr, err := readLine(reader)
	if err != nil {
		return Value{}, err
	}

	length, err := strconv.Atoi(string(lenStr))
	if err != nil {
		return Value{}, ErrInvalidSyntax
	}

	if length == -1 {
		return Value{Typ: "null", NullTyp: "bulk"}, nil
	}

	data := make([]byte, length)
	_, err = io.ReadFull(reader, data)
	if err != nil {
		return Value{}, err
	}

	crlf := make([]byte, 2)
	_, err = io.ReadFull(reader, crlf)
	if err != nil || crlf[0] != '\r' || crlf[1] != '\n' {
		return Value{}, ErrInvalidSyntax
	}

	return Value{
		Typ: "bulk",
		Str: string(data),
	}, nil
}

func parseArray(reader *bufio.Reader) (Value, error) {
	lenStr, err := readLine(reader)
	if err != nil {
		return Value{}, err
	}

	count, err := strconv.Atoi(string(lenStr))
	if err != nil {
		return Value{}, ErrInvalidSyntax
	}

	if count == -1 {
		return Value{Typ: "null", NullTyp: "array"}, nil
	}

	if count == 0 {
		return Value{Typ: "array", Array: []Value{}}, nil
	}

	array := make([]Value, count)
	for i := 0; i < count; i++ {
		element, err := Parse(reader)
		if err != nil {
			return Value{}, err
		}
		array[i] = element
	}

	return Value{
		Typ:   "array",
		Array: array,
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
