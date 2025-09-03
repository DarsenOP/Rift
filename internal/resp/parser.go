package resp

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
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

var (
	ErrCRLFMissing        = errors.New("ERR protocol error: CRLF not found")
	ErrInvalidSyntax      = errors.New("ERR Protocol error: invalid multibulk length")
	ErrInvalidBulkLength  = errors.New("ERR Protocol error: invalid bulk length")
	ErrInvalidArrayLength = errors.New("ERR Protocol error: invalid multibulk length")
)

func ParseRESP(reader *bufio.Reader) (Value, error) {
	firstByte, err := reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	if firstByte != '*' {
		return Value{}, fmt.Errorf("ERR Protocol error: expected '*', got '%c'", firstByte)
	}

	return ParseArray(reader)
}

func Parse(reader *bufio.Reader) (Value, error) {
	firstByte, err := reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch firstByte {
	case '+':
		return ParseSimpleString(reader)
	case '-':
		return ParseSimpleError(reader)
	case ':':
		return ParseInteger(reader)
	case '$':
		return ParseBulkString(reader)
	case '*':
		return ParseArray(reader)
	default:
		return Value{}, fmt.Errorf("ERR Protocol error: expected '+', '-', ':', '$', '*', got '%c'", firstByte)
	}
}

func ParseSimpleString(reader *bufio.Reader) (Value, error) {
	line, err := ReadLine(reader)
	if err != nil {
		return Value{}, err
	}

	return Value{
		Typ: "simple",
		Str: string(line),
	}, nil
}

func ParseSimpleError(reader *bufio.Reader) (Value, error) {
	line, err := ReadLine(reader)
	if err != nil {
		return Value{}, err
	}

	return Value{
		Typ: "error",
		Str: string(line),
	}, nil
}

func ParseInteger(reader *bufio.Reader) (Value, error) {
	line, err := ReadLine(reader)
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

func ParseBulkString(reader *bufio.Reader) (Value, error) {
	lenStr, err := ReadLine(reader)
	if err != nil {
		return Value{}, err
	}

	length, err := strconv.Atoi(string(lenStr))
	if err != nil {
		return Value{}, ErrInvalidBulkLength
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

func ParseArray(reader *bufio.Reader) (Value, error) {
	lenStr, err := ReadLine(reader)
	if err != nil {
		return Value{}, err
	}

	count, err := strconv.Atoi(string(lenStr))
	if err != nil {
		return Value{}, ErrInvalidArrayLength
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
func ReadLine(reader *bufio.Reader) ([]byte, error) {
	line, err := reader.ReadBytes('\n')
	if err != nil {
		// io.EOF with no data is fine; let caller decide.
		return nil, err
	}

	// Must end with \r\n
	if len(line) < 2 || line[len(line)-2] != '\r' {
		return nil, ErrCRLFMissing
	}

	return bytes.TrimSuffix(line, []byte("\r\n")), nil
}
