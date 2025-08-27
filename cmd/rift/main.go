package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/DarsenOP/Rift/internal/resp"
	"github.com/DarsenOP/Rift/pkg/version"
)

func main() {
	versionFlag := flag.Bool("version", false, "Print version information")
	testParserFlag := flag.Bool("test-parser", false, "Test RESP parser")
	testWriterFlag := flag.Bool("test-writer", false, "Test RESP writer")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("Rift version %s\n", version.Version)
		os.Exit(0)
	}

	if *testWriterFlag {
		testWriter()
	}
	if *testParserFlag {
		testParser()
	}
}

func testParser() {
	// Test cases for ALL RESP types including nested arrays
	testInputs := []string{
		// Simple Strings
		"+OK\r\n",
		"+HELLO WORLD\r\n",
		"+\r\n",

		// Simple Errors
		"-Error message\r\n",
		"-ERR unknown command\r\n",

		// Integers
		":1000\r\n",
		":-42\r\n",
		":0\r\n",

		// Bulk Strings
		"$5\r\nhello\r\n",
		"$0\r\n\r\n",
		"$-1\r\n", // Null bulk string
		"$11\r\nhello world\r\n",

		// Arrays
		"*1\r\n$4\r\nPING\r\n",
		"*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n",
		"*3\r\n$3\r\nSET\r\n$5\r\nhello\r\n$5\r\nworld\r\n",
		"*0\r\n",  // Empty array
		"*-1\r\n", // Null array

		// Nested Arrays (complex example)
		"*2\r\n*2\r\n:1\r\n:2\r\n*2\r\n:3\r\n:4\r\n",

		// Mixed types in array
		"*5\r\n+hello\r\n-err\r\n:42\r\n$5\r\nworld\r\n*2\r\n:1\r\n:2\r\n",
	}

	fmt.Println("=== TESTING RESP PARSER ===")

	for _, input := range testInputs {
		fmt.Printf("Input: %q\n", input)

		reader := bufio.NewReader(strings.NewReader(input))
		result, err := resp.Parse(reader)

		if err != nil {
			fmt.Printf("  Error: %v\n\n", err)
		} else {
			printValue(result, 0)

			// Test round-trip: parse then write back
			var buf bytes.Buffer
			err := resp.WriteValue(&buf, result)
			if err != nil {
				fmt.Printf("  Round-trip error: %v\n\n", err)
			} else {
				fmt.Printf("  Round-trip: %q\n\n", buf.String())
			}
		}
	}
}

func testWriter() {
	fmt.Println("=== TESTING RESP WRITER ===")

	// Test cases for all RESP types
	tests := []struct {
		name  string
		value resp.Value
	}{
		{"Simple String", resp.Value{Typ: "simple", Str: "OK"}},
		{"Error", resp.Value{Typ: "error", Str: "Error message"}},
		{"Integer", resp.Value{Typ: "integer", Num: 42}},
		{"Bulk String", resp.Value{Typ: "bulk", Str: "hello"}},
		{"Empty Bulk", resp.Value{Typ: "bulk", Str: ""}},
		{"Null", resp.Value{Typ: "null"}},
		{
			"PING Array",
			resp.Value{Typ: "array", Array: []resp.Value{{Typ: "bulk", Str: "PING"}}},
		},
		{
			"Mixed Array",
			resp.Value{
				Typ: "array",
				Array: []resp.Value{
					{Typ: "simple", Str: "hello"},
					{Typ: "error", Str: "err"},
					{Typ: "integer", Num: 42},
					{Typ: "bulk", Str: "world"},
				},
			},
		},
		{"Unknown Type", resp.Value{Typ: "unknown"}},
	}

	for _, tt := range tests {
		fmt.Printf("Test: %s\n", tt.name)

		var buf bytes.Buffer
		err := resp.WriteValue(&buf, tt.value)

		if err != nil {
			fmt.Printf("  Error: %v\n\n", err)
		} else {
			output := buf.String()
			fmt.Printf("  Output: %q\n", output)
			fmt.Printf("  Human: ")
			for _, b := range output {
				switch b {
				case '\r':
					fmt.Print("\\r")
				case '\n':
					fmt.Print("\\n")
				default:
					fmt.Printf("%c", b)
				}
			}
			fmt.Printf("\n\n")
		}
	}
}

// printValue recursively prints any RESP value with indentation
func printValue(v resp.Value, depth int) {
	indent := strings.Repeat("  ", depth)

	switch v.Typ {
	case "simple":
		fmt.Printf("%sSimple String: %q\n", indent, v.Str)
	case "error":
		fmt.Printf("%sError: %q\n", indent, v.Str)
	case "integer":
		fmt.Printf("%sInteger: %d\n", indent, v.Num)
	case "bulk":
		fmt.Printf("%sBulk String: %q\n", indent, v.Str)
	case "null":
		fmt.Printf("%sNull\n", indent)
	case "array":
		fmt.Printf("%sArray (%d elements):\n", indent, len(v.Array))
		for i, elem := range v.Array {
			fmt.Printf("%s  [%d]: ", indent, i)
			printValue(elem, depth+2)
		}
	default:
		fmt.Printf("%sUnknown type: %s\n", indent, v.Typ)
	}
}
