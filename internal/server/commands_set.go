package server

import (
	"github.com/DarsenOP/Rift/internal/resp"
	"github.com/DarsenOP/Rift/internal/storage"
)

// SADD key member [member ...]
func HandleSADD(s *storage.Store, args []resp.Value) resp.Value {
	if len(args) < 2 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'SADD' command"}
	}
	key := args[0].Str
	members := make([]string, len(args)-1)
	for i := 1; i < len(args); i++ {
		if args[i].Typ != "bulk" {
			return resp.Value{Typ: "error", Str: "ERR arguments must be bulk strings"}
		}
		members[i-1] = args[i].Str
	}
	added, err := s.SAdd(key, members...)
	if err != nil {
		if err == storage.ErrWrongType {
			return wrongType()
		}
		return resp.Value{Typ: "error", Str: "ERR " + err.Error()}
	}
	return resp.Value{Typ: "integer", Num: added}
}

// SREM key member [member ...]
func HandleSREM(s *storage.Store, args []resp.Value) resp.Value {
	if len(args) < 2 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'SREM' command"}
	}
	key := args[0].Str
	members := make([]string, len(args)-1)
	for i := 1; i < len(args); i++ {
		if args[i].Typ != "bulk" {
			return resp.Value{Typ: "error", Str: "ERR arguments must be bulk strings"}
		}
		members[i-1] = args[i].Str
	}
	removed, err := s.SRem(key, members...)
	if err != nil {
		if err == storage.ErrWrongType {
			return wrongType()
		}
		return resp.Value{Typ: "error", Str: "ERR " + err.Error()}
	}
	return resp.Value{Typ: "integer", Num: removed}
}

// SISMEMBER key member
func HandleSISMEMBER(s *storage.Store, args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'SISMEMBER' command"}
	}
	if args[0].Typ != "bulk" || args[1].Typ != "bulk" {
		return resp.Value{Typ: "error", Str: "ERR arguments must be bulk strings"}
	}
	ok, err := s.SIsMember(args[0].Str, args[1].Str)
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

// SMEMBERS key
func HandleSMEMBERS(s *storage.Store, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'SMEMBERS' command"}
	}
	if args[0].Typ != "bulk" {
		return resp.Value{Typ: "error", Str: "ERR key must be a bulk string"}
	}
	members, err := s.SMembers(args[0].Str)
	if err != nil {
		if err == storage.ErrNotFound {
			return resp.Value{Typ: "array", Array: []resp.Value{}}
		}
		if err == storage.ErrWrongType {
			return wrongType()
		}
		return resp.Value{Typ: "error", Str: "ERR " + err.Error()}
	}
	out := make([]resp.Value, len(members))
	for i, m := range members {
		out[i] = resp.Value{Typ: "bulk", Str: m}
	}
	return resp.Value{Typ: "array", Array: out}
}

// SCARD key
func HandleSCARD(s *storage.Store, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'SCARD' command"}
	}
	if args[0].Typ != "bulk" {
		return resp.Value{Typ: "error", Str: "ERR key must be a bulk string"}
	}
	n, err := s.SCard(args[0].Str)
	if err != nil {
		if err == storage.ErrWrongType {
			return wrongType()
		}
		return resp.Value{Typ: "error", Str: "ERR " + err.Error()}
	}
	return resp.Value{Typ: "integer", Num: n}
}

// SINTER key [key ...]
func HandleSINTER(s *storage.Store, args []resp.Value) resp.Value {
	if len(args) < 1 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'SINTER' command"}
	}
	keys := make([]string, len(args))
	for i, a := range args {
		if a.Typ != "bulk" {
			return resp.Value{Typ: "error", Str: "ERR arguments must be bulk strings"}
		}
		keys[i] = a.Str
	}
	inter, err := s.SInter(keys...)
	if err != nil {
		if err == storage.ErrWrongType {
			return wrongType()
		}
		return resp.Value{Typ: "error", Str: "ERR " + err.Error()}
	}
	out := make([]resp.Value, len(inter))
	for i, m := range inter {
		out[i] = resp.Value{Typ: "bulk", Str: m}
	}
	return resp.Value{Typ: "array", Array: out}
}
