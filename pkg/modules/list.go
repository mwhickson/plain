package modules

import (
	"fmt"
	"plain/pkg/value"
)

func GetListModule() value.Value {
	listModule := map[string]value.Value{
		"Len":      value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(listLen)},
		"Contains": value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(listContains)},
		"Append":   value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(listAppend)},
		"Reverse":  value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(listReverse)},
	}
	return value.Value{
		Type: value.ValStruct,
		Obj:  listModule,
	}
}

func listLen(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Type != value.ValList {
		return value.NewNone(), fmt.Errorf("List.Len() requires 1 List argument")
	}
	list := args[0].Obj.([]value.Value)
	return value.NewNumeric(float64(len(list))), nil
}

func listContains(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Type != value.ValList {
		return value.NewNone(), fmt.Errorf("List.Contains() requires (List, Value)")
	}
	list := args[0].Obj.([]value.Value)
	target := args[1]

	for _, item := range list {
		// Basic equality check matching VM Eq logic
		if item.Type == target.Type && item.Num == target.Num && item.Str == target.Str {
			return value.NewBoolean(true), nil
		}
	}
	return value.NewBoolean(false), nil
}

func listAppend(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Type != value.ValList {
		return value.NewNone(), fmt.Errorf("List.Append() requires (List, Value)")
	}
	list := args[0].Obj.([]value.Value)
	updated := append(list, args[1])
	return value.Value{Type: value.ValList, Obj: updated}, nil
}

func listReverse(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Type != value.ValList {
		return value.NewNone(), fmt.Errorf("List.Reverse() requires 1 List argument")
	}
	original := args[0].Obj.([]value.Value)
	reversed := make([]value.Value, len(original))
	for i, item := range original {
		reversed[len(original)-1-i] = item
	}
	return value.Value{Type: value.ValList, Obj: reversed}, nil
}
