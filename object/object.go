package object

import (
	"bytes"
	"fmt"
	"interpreter/ast"
	"strings"
)

type ObjectType string

const (
    INTEGER_OBJ = "INTEGER"
    BOOLEAN_OBJ = "BOOLEAN"
    NULL_OBJ = "NULL"
    RETURN_VALUE_OBJ = "RETURN_VALUE"
    ERRROR_OBJ = "ERROR"
    FUNCTION_OBJ = "FUNCTION"
    STRING_OBJ = "STRING"
)

type Object interface {
    Type() ObjectType
    Inspect() string
}

type Integer struct {
    Value int64
}

func (i *Integer) Inspect() string {
    return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) Type() ObjectType {
    return INTEGER_OBJ
}

type Boolean struct {
    Value bool
}

func (b *Boolean) Inspect() string {
    return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) Type() ObjectType {
    return BOOLEAN_OBJ
}

type Null struct {}

func (n *Null) Inspect() string {
    return "null"
}

func (n *Null) Type() ObjectType {
    return NULL_OBJ
}

type RetrunValue struct {
    Value Object
}

func (rv *RetrunValue) Inspect() string {
    return rv.Value.Inspect()
}

func (rv *RetrunValue) Type() ObjectType {
    return RETURN_VALUE_OBJ
}

type Function struct {
    Parameters []*ast.Indentifier
    Body *ast.BlockStatement
    Env *Enviroment
}

func (f *Function) Inspect() string {
    var out bytes.Buffer
    var params []string

    for _, p := range f.Parameters {
        params = append(params, p.String())
    }

    out.WriteString("fn(")
    out.WriteString(strings.Join(params, ", "))
    out.WriteString(") {\n")
    out.WriteString(f.Body.String())
    out.WriteString("\n}")

    return out.String()
}

func (f *Function) Type() ObjectType {
    return FUNCTION_OBJ
}

type Error struct {
    Message string
}

func (e *Error) Inspect() string {
    return fmt.Sprintf("ERROR: %s", e.Message)
}

func (e *Error) Type() ObjectType {
    return ERRROR_OBJ
}

type String struct {
    Value string
}

func (s *String) Inspect() string {
    return s.Value
}

func (s *String) Type() ObjectType {
    return STRING_OBJ
}
