package object

import "fmt"

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL_OBJ    = "NULL"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (integer *Integer) Type() ObjectType { return INTEGER_OBJ }
func (integer *Integer) Inspect() string  { return fmt.Sprintf("%d", integer.Value) }

type Boolean struct {
	Value bool
}

func (boolean *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (boolean *Boolean) Inspect() string  { return fmt.Sprintf("%t", boolean.Value) }

type Null struct{}

func (null *Null) Type() ObjectType { return NULL_OBJ }
func (null *Null) Inspect() string  { return fmt.Sprintf("null") }