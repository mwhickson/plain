package value

import "fmt"

const (
	ValNone ValueType = iota
	ValBoolean
	ValClosure
	ValFunction
	ValList
	ValNumeric
	ValString
	ValStruct
)

type ClosureObject struct {
	Fn   *FunctionObject
	Free []Value
}

type FunctionObject struct {
	Name                  string
	InstructionStartIndex int
	ParamCount            int
	LocalVarCount         int
}

// Value represents a dynamic scalar or composite object in the PLAIN VM.
type Value struct {
	Type ValueType
	Num  float64
	Str  string
	Bool bool
	Obj  any // Holds List ([]Value) or Struct (map[string]Value)
}

// ValueType represents the internal type tag for PLAIN runtime values.
type ValueType byte

// Helper Constructors

func NewNone() Value {
	return Value{Type: ValNone}
}

func NewBoolean(b bool) Value {
	var numVal float64
	if b {
		numVal = 1
	} else {
		numVal = 0
	}
	return Value{Type: ValBoolean, Num: numVal, Bool: b}
}

func NewClosure(fn *FunctionObject, free []Value) Value {
	return Value{
		Type: ValClosure,
		Obj:  &ClosureObject{Fn: fn, Free: free},
	}
}

func NewNumeric(n float64) Value {
	return Value{Type: ValNumeric, Num: n}
}

func NewString(s string) Value {
	return Value{Type: ValString, Str: s} // assuming Str field exists
}

func (v Value) MatchesTypeName(typeName string) bool {
	switch typeName {
	case "Boolean":
		return v.Type == ValBoolean
	case "None":
		return v.Type == ValNone
	case "Numeric":
		return v.Type == ValNumeric
	case "String":
		return v.Type == ValString
	default:
		return true // Custom structs/lists or untyped can pass for now
	}
}

// String returns a human-readable representation of the Value.
func (v Value) String() string {
	switch v.Type {
	case ValNone:
		return "None"
	case ValBoolean:
		if v.Bool {
			return "True"
		}
		return "False"
	case ValClosure:
		cl := v.Obj.(*ClosureObject)
		return fmt.Sprintf("<closure %s>", cl.Fn.Name)
	case ValFunction:
		fn := v.Obj.(*FunctionObject)
		return fmt.Sprintf("<fn %s>", fn.Name)
	case ValList:
		return "[List]"
	case ValNumeric:
		return fmt.Sprintf("%g", v.Num)
	case ValString:
		return v.Str
	case ValStruct:
		return "[Struct]"
	default:
		return "Unknown"
	}
}

func (vt ValueType) String() string {
	switch vt {
	case ValNone:
		return "None"
	case ValBoolean:
		return "Boolean"
	case ValClosure:
		return "Closure"
	case ValFunction:
		return "Function"
	case ValList:
		return "List"
	case ValNumeric:
		return "Numeric"
	case ValString:
		return "String"
	case ValStruct:
		return "Struct"
	default:
		return "Unknown"
	}
}
