package resp

import (
	"io"
	"strconv"
)

// WriteValue writes any RESP value to the writer
func WriteValue(w io.Writer, v Value) error {
	switch v.Typ {
	case "simple":
		return WriteSimpleString(w, v.Str)
	case "error":
		return WriteError(w, v.Str)
	case "integer":
		return WriteInteger(w, v.Num)
	case "bulk":
		return WriteBulkString(w, v.Str)
	case "null":
		switch v.NullTyp {
		case "bulk":
			return WriteBulkNull(w)
		case "array":
			return WriteArrayNull(w)
		default:
			return WriteError(w, "ERR unknown null type")
		}
	case "array":
		return WriteArray(w, v.Array)
	default:
		return WriteError(w, "ERR unknown type")
	}
}

func WriteSimpleString(w io.Writer, s string) error {
	_, err := w.Write([]byte("+" + s + "\r\n"))
	return err
}

func WriteError(w io.Writer, s string) error {
	_, err := w.Write([]byte("-" + s + "\r\n"))
	return err
}

func WriteInteger(w io.Writer, n int) error {
	_, err := w.Write([]byte(":" + strconv.Itoa(n) + "\r\n"))
	return err
}

func WriteBulkString(w io.Writer, s string) error {
	if s == "" {
		_, err := w.Write([]byte("$0\r\n\r\n"))
		return err
	}

	_, err := w.Write([]byte("$" + strconv.Itoa(len(s)) + "\r\n" + s + "\r\n"))
	return err
}

// TODO: RESP3 would use '_' for null, but we're using RESP2 for compatibility
// For now, we use $-1 and *-1 for bulk/array nulls for simplicity
func WriteBulkNull(w io.Writer) error {
	_, err := w.Write([]byte("$-1\r\n"))
	return err
}

func WriteArrayNull(w io.Writer) error {
	_, err := w.Write([]byte("*-1\r\n"))
	return err
}

func WriteArray(w io.Writer, array []Value) error {
	_, err := w.Write([]byte("*" + strconv.Itoa(len(array)) + "\r\n"))
	if err != nil {
		return err
	}

	for _, element := range array {
		err := WriteValue(w, element)
		if err != nil {
			return err
		}
	}

	return nil
}
