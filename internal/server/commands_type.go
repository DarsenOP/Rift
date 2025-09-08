package server

import (
	"github.com/DarsenOP/Rift/internal/resp"
	"github.com/DarsenOP/Rift/internal/storage"
)

func HandleTYPE(s *storage.Store, args []resp.Value) resp.Value {
	if len(args) != 1 || args[0].Typ != "bulk" {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'TYPE' command"}
	}
	t := s.Type(args[0].Str)
	return resp.Value{Typ: "bulk", Str: t}
}

func HandleRENAME(s *storage.Store, args []resp.Value) resp.Value {
	if len(args) != 2 || args[0].Typ != "bulk" || args[1].Typ != "bulk" {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'RENAME' command"}
	}
	if err := s.Rename(args[0].Str, args[1].Str); err != nil {
		if err == storage.ErrNotFound {
			return resp.Value{Typ: "error", Str: "ERR no such key"}
		}
		return resp.Value{Typ: "error", Str: "ERR " + err.Error()}
	}
	return resp.Value{Typ: "simple", Str: "OK"}
}

func HandleRENAMENX(s *storage.Store, args []resp.Value) resp.Value {
	if len(args) != 2 || args[0].Typ != "bulk" || args[1].Typ != "bulk" {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'RENAMENX' command"}
	}
	ok, err := s.RenameNX(args[0].Str, args[1].Str)
	if err != nil {
		if err == storage.ErrNotFound {
			return resp.Value{Typ: "error", Str: "ERR no such key"}
		}
		return resp.Value{Typ: "error", Str: "ERR " + err.Error()}
	}
	return resp.Value{Typ: "integer", Num: boolToInt(ok)}
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
