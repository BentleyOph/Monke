package object

import (
	"fmt"

	"github.com/BentleyOph/monke/ast"
	"bytes"
	"strings"
)

type ObjectType string

// Define the object types
const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	NULL_OBJ    = "NULL"
	ERROR_OBJ   = "ERROR"
	FUNCTION_OBJ = "FUNCTION"
	STRING_OBJ = "STRING"	
	BUILTIN_OBJ = "BUILTIN"
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

// Boolean BOOLEANS
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}
func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

// Null NULL
type Null struct{}

func (n *Null) Type() ObjectType {
	return NULL_OBJ
}

func (n *Null) Inspect() string {
	return "null"
}


// ReturnValue RETURN VALUES
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type () ObjectType{
	return RETURN_VALUE_OBJ
}
func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

// Error ERRORS
type Error struct {
	Message string
}
func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}
func (e *Error) Inspect() string {
	return "ERROR:" + e.Message
}


// environment
type Environment struct {
	store map[string]Object
	outer *Environment // reference to the outer environment(enclosing environment)
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer:nil}
}

func(e *Environment) Get(name string)(Object, bool){
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok 
}
func (e *Environment) Set (name string, val Object) Object {
	e.store[name]= val
	return val 
}
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}


// String objects
type String struct {
	Value string
}

func (s *String) Type() ObjectType {
	return STRING_OBJ
}
func (s *String) Inspect() string {
	return s.Value
}








// Function objects
type Function struct {
	Parameters []*ast.Identifier
	Body *ast.BlockStatement
	Env *Environment //function's environment
}

func (f * Function) Type() ObjectType {
	return FUNCTION_OBJ
}
func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _,p := range f.Parameters{
		params = append(params, p.String())
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}


type BuiltinFunction func(args ... Object) Object 

type Builtin struct {
	Fn BuiltinFunction
}

func(b *Builtin) Type() ObjectType {
	return BUILTIN_OBJ
}
func (b *Builtin) Inspect() string{
	return "builtin function";
}