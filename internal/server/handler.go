package server

import (
	"strconv"
	"strings"
	"time"

	"github.com/DarsenOP/Rift/internal/resp"
	"github.com/DarsenOP/Rift/internal/storage"
)

func HandleCommand(store *storage.Store, value resp.Value) resp.Value {
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
	case "SET":
		return HandleSET(store, value.Array[1:])
	case "GET":
		return HandleGET(store, value.Array[1:])
	case "DEL":
		return HandleDEL(store, value.Array[1:])
	case "EXISTS":
		return HandleEXISTS(store, value.Array[1:])
	case "TTL":
		return HandleTTL(store, value.Array[1:])
	case "EXPIRE":
		return HandleEXPIRE(store, value.Array[1:])
	case "LPUSH":
		return HandleLPUSH(store, value.Array[1:])
	case "RPUSH":
		return HandleRPUSH(store, value.Array[1:])
	case "LPOP":
		return HandleLPOP(store, value.Array[1:])
	case "RPOP":
		return HandleRPOP(store, value.Array[1:])
	case "LRANGE":
		return HandleLRANGE(store, value.Array[1:])
	case "LLEN":
		return HandleLLEN(store, value.Array[1:])
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

func HandleSET(store *storage.Store, args []resp.Value) resp.Value {
	if len(args) < 2 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'set' command"}
	}

	if args[0].Typ != "bulk" || args[1].Typ != "bulk" {
		return resp.Value{Typ: "error", Str: "ERR arguments should be bulk strings"}
	}

	key, value := args[0].Str, args[1].Str
	var ttl time.Duration

	// parse optional flags
	if len(args) > 2 {
		if len(args) != 4 {
			return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'set' with expiration"}
		}

		flag, flagVal := strings.ToUpper(args[2].Str), args[3].Str
		if args[2].Typ != "bulk" || args[3].Typ != "bulk" {
			return resp.Value{Typ: "error", Str: "ERR expiration flag and value must be bulk strings"}
		}

		var multiplier time.Duration
		switch flag {
		case "EX":
			multiplier = time.Second
		case "PX":
			multiplier = time.Millisecond
		default:
			return resp.Value{Typ: "error", Str: "ERR unsupported option"}
		}

		n, err := strconv.Atoi(flagVal)
		if err != nil || n <= 0 {
			return resp.Value{Typ: "error", Str: "ERR value is not an integer or out of range"}
		}
		ttl = time.Duration(n) * multiplier
	}

	// atomic write
	if ttl > 0 {
		store.SetEX(key, value, ttl)
	} else {
		store.Set(key, value)
	}

	return resp.Value{Typ: "simple", Str: "OK"}
}

func HandleGET(store *storage.Store, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'get' command"}
	}

	if args[0].Typ != "bulk" {
		return resp.Value{Typ: "error", Str: "ERR argument should be a bulk string"}
	}

	value, ok := store.Get(args[0].Str)

	if !ok || value == "" {
		return resp.Value{Typ: "null", NullTyp: "bulk"}
	}

	return resp.Value{Typ: "bulk", Str: value}
}

func HandleDEL(store *storage.Store, args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'del' command"}
	}

	for _, arg := range args {
		if arg.Typ != "bulk" {
			return resp.Value{Typ: "error", Str: "ERR arguments should be bulk strings"}
		}
	}

	deleted := 0
	for _, arg := range args {
		if store.Del(arg.Str) {
			deleted++
		}
	}
	return resp.Value{Typ: "integer", Num: deleted}
}

func HandleEXISTS(store *storage.Store, args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'exists' command"}
	}

	for _, arg := range args {
		if arg.Typ != "bulk" {
			return resp.Value{Typ: "error", Str: "ERR arguments should be bulk strings"}
		}
	}

	count := 0
	for _, arg := range args {
		if store.Exists(arg.Str) {
			count++
		}
	}
	return resp.Value{Typ: "integer", Num: count}
}

func HandleTTL(store *storage.Store, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'ttl' command"}
	}
	if args[0].Typ != "bulk" {
		return resp.Value{Typ: "error", Str: "ERR argument should be a bulk string"}
	}

	ttl, _ := store.TTL(args[0].Str)

	switch ttl {
	case -2:
		return resp.Value{Typ: "integer", Num: -2}
	case -1:
		return resp.Value{Typ: "integer", Num: -1}
	default:
		return resp.Value{Typ: "integer", Num: int(ttl.Seconds())}
	}
}

func HandleEXPIRE(store *storage.Store, args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'expire' command"}
	}

	if args[0].Typ != "bulk" || args[1].Typ != "bulk" {
		return resp.Value{Typ: "error", Str: "ERR arguments should be bulk strings"}
	}

	// parse seconds
	secs, err := strconv.Atoi(args[1].Str)
	if err != nil || secs <= 0 {
		return resp.Value{Typ: "error", Str: "ERR value is not an integer or out of range"}
	}

	ok := store.Expire(args[0].Str, time.Duration(secs)*time.Second)
	if ok {
		return resp.Value{Typ: "integer", Num: 1}
	}

	return resp.Value{Typ: "integer", Num: 0}
}
