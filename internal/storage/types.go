package storage

import "time"

// DataType represents the type of data stored
type DataType string

const (
	StringType DataType = "string"
	ListType   DataType = "list"
	HashType   DataType = "hash"
	SetType    DataType = "set"
)

// StringValue represents a string data type
type StringValue struct {
	Value  string
	Expiry *time.Time
}

// ListValue represents a list data type (slice of strings)
type ListValue struct {
	Values []string
	Expiry *time.Time
}

// HashValue represents a hash data type (map of string fields)
type HashValue struct {
	Fields map[string]string
	Expiry *time.Time
}

// SetValue represents a set data type (map for O(1) lookups)
type SetValue struct {
	Members map[string]struct{}
	Expiry  *time.Time
}

// Value represents a generic stored value with type information
type Value struct {
	Type   DataType
	String *StringValue
	List   *ListValue
	Hash   *HashValue
	Set    *SetValue
}

// Helper methods for value creation
func NewStringValue(value string) *Value {
	return &Value{
		Type:   StringType,
		String: &StringValue{Value: value},
	}
}

func NewListValue(values []string) *Value {
	return &Value{
		Type: ListType,
		List: &ListValue{Values: values},
	}
}

func NewHashValue(fields map[string]string) *Value {
	return &Value{
		Type: HashType,
		Hash: &HashValue{Fields: fields},
	}
}

func NewSetValue(members []string) *Value {
	memberMap := make(map[string]struct{})
	for _, m := range members {
		memberMap[m] = struct{}{}
	}
	return &Value{
		Type: SetType,
		Set:  &SetValue{Members: memberMap},
	}
}
