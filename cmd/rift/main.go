package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/DarsenOP/Rift/internal/resp"
	"github.com/DarsenOP/Rift/pkg/version"
)

func main() {
	versionFlag := flag.Bool("version", false, "Print version information")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("Rift version %s\n", version.Version)
		os.Exit(0)
	}

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

	for _, input := range testInputs {
		fmt.Printf("Testing: %q\n", input)

		reader := bufio.NewReader(strings.NewReader(input))
		result, err := resp.Parse(reader)

		if err != nil {
			fmt.Printf("  Error: %v\n\n", err)
		} else {
			printValue(result, 0)
			fmt.Println()
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
