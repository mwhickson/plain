package value

import "fmt"

// ValueType represents the internal type tag for PLAIN runtime values.
type ValueType byte

const (
	ValNone ValueType = iota
	ValNumeric
	ValString
	ValBoolean
	ValList
	ValStruct
	ValFunction
)

// Value represents a dynamic scalar or composite object in the PLAIN VM.
type Value struct {
	Type ValueType
	Num  float64
	Str  string
	Bool bool
	Obj  any // Holds List ([]Value) or Struct (map[string]Value)
}

func (vt ValueType) String() string {
	switch vt {
	case ValNone:
		return "None"
	case ValNumeric:
		return "Numeric"
	case ValString:
		return "String"
	case ValBoolean:
		return "Boolean"
	case ValList:
		return "List"
	case ValStruct:
		return "Struct"
	default:
		return "Unknown"
	}
}

// Helper Constructors

func NewNone() Value {
	return Value{Type: ValNone}
}

func NewNumeric(n float64) Value {
	return Value{Type: ValNumeric, Num: n}
}

func NewString(s string) Value {
	return Value{Type: ValString, Str: s} // assuming Str field exists
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

// String returns a human-readable representation of the Value.
func (v Value) String() string {
	switch v.Type {
	case ValNone:
		return "None"
	case ValNumeric:
		return fmt.Sprintf("%g", v.Num)
	case ValString:
		return v.Str
	case ValBoolean:
		if v.Bool {
			return "True"
		}
		return "False"
	case ValList:
		return "[List]"
	case ValStruct:
		return "[Struct]"
	case ValFunction:
		fn := v.Obj.(*FunctionObject)
		return fmt.Sprintf("<fn %s>", fn.Name)
	default:
		return "Unknown"
	}
}

type FunctionObject struct {
	Name                  string
	InstructionStartIndex int
	ParamCount            int
}

func (v Value) MatchesTypeName(typeName string) bool {
	switch typeName {
	case "Numeric":
		return v.Type == ValNumeric
	case "String":
		return v.Type == ValString
	case "Boolean":
		return v.Type == ValBoolean
	case "None":
		return v.Type == ValNone
	default:
		return true // Custom structs/lists or untyped can pass for now
	}
}
