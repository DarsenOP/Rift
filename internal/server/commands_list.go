package server

import (
	"strconv"

	"github.com/DarsenOP/Rift/internal/resp"
	"github.com/DarsenOP/Rift/internal/storage"
)

func HandleLPUSH(store *storage.Store, args []resp.Value) resp.Value {
	if len(args) < 2 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'LPUSH' command"}
	}

	key := args[0].Str
	values := make([]string, len(args)-1)
	for i, arg := range args[1:] {
		if arg.Typ != "bulk" {
			return resp.Value{Typ: "error", Str: "ERR arguments should be bulk strings"}
		}
		values[i] = arg.Str
	}

	length, err := store.LPush(key, values...)
	if err != nil {
		if err == storage.ErrWrongType {
			return resp.Value{Typ: "error", Str: "WRONGTYPE Operation against a key holding the wrong kind of value"}
		}
		return resp.Value{Typ: "error", Str: "ERR " + err.Error()}
	}

	return resp.Value{Typ: "integer", Num: length}
}

func HandleRPUSH(store *storage.Store, args []resp.Value) resp.Value {
	if len(args) < 2 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'RPUSH' command"}
	}

	key := args[0].Str
	values := make([]string, len(args)-1)
	for i, arg := range args[1:] {
		if arg.Typ != "bulk" {
			return resp.Value{Typ: "error", Str: "ERR arguments should be bulk strings"}
		}
		values[i] = arg.Str
	}

	length, err := store.RPush(key, values...)
	if err != nil {
		if err == storage.ErrWrongType {
			return resp.Value{Typ: "error", Str: "WRONGTYPE Operation against a key holding the wrong kind of value"}
		}
		return resp.Value{Typ: "error", Str: "ERR " + err.Error()}
	}

	return resp.Value{Typ: "integer", Num: length}
}

func HandleLPOP(store *storage.Store, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'LPOP' command"}
	}

	if args[0].Typ != "bulk" {
		return resp.Value{Typ: "error", Str: "ERR argument should be a bulk string"}
	}

	value, err := store.LPop(args[0].Str)
	if err != nil {
		if err == storage.ErrWrongType {
			return resp.Value{Typ: "error", Str: "WRONGTYPE Operation against a key holding the wrong kind of value"}
		}
		if err == storage.ErrNotFound {
			return resp.Value{Typ: "null", NullTyp: "bulk"} // Redis returns nil for non-existent keys
		}
		return resp.Value{Typ: "error", Str: "ERR " + err.Error()}
	}

	if value == "" {
		return resp.Value{Typ: "null", NullTyp: "bulk"}
	}

	return resp.Value{Typ: "bulk", Str: value}
}

func HandleRPOP(store *storage.Store, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'RPOP' command"}
	}

	if args[0].Typ != "bulk" {
		return resp.Value{Typ: "error", Str: "ERR argument should be a bulk string"}
	}

	value, err := store.RPop(args[0].Str)
	if err != nil {
		if err == storage.ErrWrongType {
			return resp.Value{Typ: "error", Str: "WRONGTYPE Operation against a key holding the wrong kind of value"}
		}
		if err == storage.ErrNotFound {
			return resp.Value{Typ: "null", NullTyp: "bulk"} // Redis returns nil for non-existent keys
		}
		return resp.Value{Typ: "error", Str: "ERR " + err.Error()}
	}

	if value == "" {
		return resp.Value{Typ: "null", NullTyp: "bulk"}
	}

	return resp.Value{Typ: "bulk", Str: value}
}

func HandleLRANGE(store *storage.Store, args []resp.Value) resp.Value {
	if len(args) != 3 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'LRANGE' command"}
	}

	if args[0].Typ != "bulk" || args[1].Typ != "bulk" || args[2].Typ != "bulk" {
		return resp.Value{Typ: "error", Str: "ERR arguments should be bulk strings"}
	}

	start, err1 := strconv.Atoi(args[1].Str)
	stop, err2 := strconv.Atoi(args[2].Str)

	if err1 != nil || err2 != nil {
		return resp.Value{Typ: "error", Str: "ERR value is not an integer or out of range"}
	}

	elements, err := store.LRange(args[0].Str, start, stop)
	if err != nil {
		if err == storage.ErrWrongType {
			return resp.Value{Typ: "error", Str: "WRONGTYPE Operation against a key holding the wrong kind of value"}
		}
		if err == storage.ErrNotFound {
			return resp.Value{Typ: "array", Array: []resp.Value{}} // Redis returns empty array for non-existent keys
		}
		return resp.Value{Typ: "error", Str: "ERR " + err.Error()}
	}

	// Build RESP array response
	respArray := make([]resp.Value, len(elements))
	for i, element := range elements {
		respArray[i] = resp.Value{Typ: "bulk", Str: element}
	}

	return resp.Value{Typ: "array", Array: respArray}
}

func HandleLLEN(store *storage.Store, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'LLEN' command"}
	}

	if args[0].Typ != "bulk" {
		return resp.Value{Typ: "error", Str: "ERR argument should be a bulk string"}
	}

	length, err := store.LLen(args[0].Str)
	if err != nil {
		if err == storage.ErrWrongType {
			return resp.Value{Typ: "error", Str: "WRONGTYPE Operation against a key holding the wrong kind of value"}
		}
		if err == storage.ErrNotFound {
			return resp.Value{Typ: "integer", Num: 0} // Redis returns 0 for non-existent keys
		}
		return resp.Value{Typ: "error", Str: "ERR " + err.Error()}
	}

	return resp.Value{Typ: "integer", Num: length}
}
