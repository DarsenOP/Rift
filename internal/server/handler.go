package server

import (
	"strings"

	"github.com/DarsenOP/Rift/internal/resp"
)

func HandleCommand(value resp.Value) resp.Value {
	if len(value.Array) == 0 {
		return resp.Value{Typ: "array", Array: []resp.Value{}}
	}

	// Extract command name (first element of array)
	command := value.Array[0]
	if command.Typ != "bulk" {
		return resp.Value{Typ: "error", Str: "ERR command must be a bulk string"}
	}

	commandName := strings.ToUpper(command.Str)

	// Handle different commands
	switch commandName {
	case "PING":
		return HandlePING(value.Array[1:])
	case "COMMAND":
		return HandleCOMMAND()
	default:
		return resp.Value{Typ: "error", Str: "ERR unknown command '" + command.Str + "'"}
	}
}

func HandlePING(args []resp.Value) resp.Value {
	if len(args) == 0 {
		// Simple PING: return PONG
		return resp.Value{Typ: "simple", Str: "PONG"}
	}

	if len(args) == 1 {
		// PING with message: return the message
		if args[0].Typ == "integer" {
			return resp.Value{Typ: "integer", Num: args[0].Num}
		}

		return resp.Value{Typ: "bulk", Str: args[0].Str}
	}

	return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'ping' command"}
}

func HandleCOMMAND() resp.Value {
	// Return empty array for now (basic Redis compatibility)
	return resp.Value{Typ: "array", Array: []resp.Value{}}
}
