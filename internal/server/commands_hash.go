package server

import (
	"github.com/DarsenOP/Rift/internal/resp"
	"github.com/DarsenOP/Rift/internal/storage"
)

// HSET key field value [field value ...]
func HandleHSET(s *storage.Store, args []resp.Value) resp.Value {
	if len(args) < 3 || len(args)%2 == 0 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'HSET' command"}
	}
	key := args[0].Str
	fv := make([]string, len(args)-1)
	for i := 1; i < len(args); i++ {
		if args[i].Typ != "bulk" {
			return resp.Value{Typ: "error", Str: "ERR arguments must be bulk strings"}
		}
		fv[i-1] = args[i].Str
	}
	added, err := s.HSet(key, fv...)
	if err != nil {
		if err == storage.ErrWrongType {
			return wrongType()
		}
		return resp.Value{Typ: "error", Str: "ERR " + err.Error()}
	}
	return resp.Value{Typ: "integer", Num: added}
}

// HGET key field
func HandleHGET(s *storage.Store, args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'HGET' command"}
	}
	if args[0].Typ != "bulk" || args[1].Typ != "bulk" {
		return resp.Value{Typ: "error", Str: "ERR arguments must be bulk strings"}
	}
	val, err := s.HGet(args[0].Str, args[1].Str)
	if err != nil {
		if err == storage.ErrNotFound {
			return nullBulk()
		}
		if err == storage.ErrWrongType {
			return wrongType()
		}
		return resp.Value{Typ: "error", Str: "ERR " + err.Error()}
	}
	return resp.Value{Typ: "bulk", Str: val}
}

// HGETALL key  -> alternating field/value array
func HandleHGETALL(s *storage.Store, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'HGETALL' command"}
	}
	if args[0].Typ != "bulk" {
		return resp.Value{Typ: "error", Str: "ERR key must be a bulk string"}
	}
	arr, err := s.HGetAll(args[0].Str)
	if err != nil {
		if err == storage.ErrNotFound {
			return resp.Value{Typ: "array", Array: []resp.Value{}}
		}
		if err == storage.ErrWrongType {
			return wrongType()
		}
		return resp.Value{Typ: "error", Str: "ERR " + err.Error()}
	}
	out := make([]resp.Value, len(arr))
	for i, v := range arr {
		out[i] = resp.Value{Typ: "bulk", Str: v}
	}
	return resp.Value{Typ: "array", Array: out}
}

// HDEL key field [field ...]
func HandleHDEL(s *storage.Store, args []resp.Value) resp.Value {
	if len(args) < 2 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'HDEL' command"}
	}
	fields := make([]string, len(args)-1)
	for i := 1; i < len(args); i++ {
		if args[i].Typ != "bulk" {
			return resp.Value{Typ: "error", Str: "ERR arguments must be bulk strings"}
		}
		fields[i-1] = args[i].Str
	}
	removed, err := s.HDel(args[0].Str, fields...)
	if err != nil {
		if err == storage.ErrWrongType {
			return wrongType()
		}
		return resp.Value{Typ: "error", Str: "ERR " + err.Error()}
	}
	return resp.Value{Typ: "integer", Num: removed}
}

// HEXISTS key field
func HandleHEXISTS(s *storage.Store, args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'HEXISTS' command"}
	}
	if args[0].Typ != "bulk" || args[1].Typ != "bulk" {
		return resp.Value{Typ: "error", Str: "ERR arguments must be bulk strings"}
	}
	ok, err := s.HExists(args[0].Str, args[1].Str)
	if err != nil {
		if err == storage.ErrWrongType {
			return wrongType()
		}
		return resp.Value{Typ: "error", Str: "ERR " + err.Error()}
	}
	if ok {
		return resp.Value{Typ: "integer", Num: 1}
	}
	return resp.Value{Typ: "integer", Num: 0}
}

// HLEN key
func HandleHLEN(s *storage.Store, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'HLEN' command"}
	}
	if args[0].Typ != "bulk" {
		return resp.Value{Typ: "error", Str: "ERR key must be a bulk string"}
	}
	l, err := s.HLen(args[0].Str)
	if err != nil {
		if err == storage.ErrNotFound {
			return resp.Value{Typ: "integer", Num: 0}
		}
		if err == storage.ErrWrongType {
			return wrongType()
		}
		return resp.Value{Typ: "error", Str: "ERR " + err.Error()}
	}
	return resp.Value{Typ: "integer", Num: l}
}

// helpers
func wrongType() resp.Value {
	return resp.Value{Typ: "error", Str: "WRONGTYPE Operation against a key holding the wrong kind of value"}
}

func nullBulk() resp.Value {
	return resp.Value{Typ: "null", NullTyp: "bulk"}
}
