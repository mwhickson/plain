package modules

import (
	"fmt"
	"math"
	"plain/pkg/value"
)

func GetListModule() value.Value {
	listModule := map[string]value.Value{
		"Len":      value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(listLen)},
		"Contains": value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(listContains)},
		"Append":   value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(listAppend)},
		"Delete":   value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(listDelete)},
		"Reverse":  value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(listReverse)},
		"Update":   value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(listUpdate)},
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
	updated := make([]value.Value, len(list)+1)
	copy(updated, list)
	updated[len(list)] = args[1]
	return value.Value{Type: value.ValList, Obj: updated}, nil
}

func listDelete(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Type != value.ValList {
		return value.NewNone(), fmt.Errorf("List.Delete() requires (List, Numeric index)")
	}
	list := args[0].Obj.([]value.Value)
	index, err := listIndex(list, args[1], "List.Delete()")
	if err != nil {
		return value.NewNone(), err
	}

	updated := make([]value.Value, len(list)-1)
	copy(updated, list[:index])
	copy(updated[index:], list[index+1:])
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

func listUpdate(args []value.Value) (value.Value, error) {
	if len(args) != 3 || args[0].Type != value.ValList {
		return value.NewNone(), fmt.Errorf("List.Update() requires (List, Numeric index, Value)")
	}
	list := args[0].Obj.([]value.Value)
	index, err := listIndex(list, args[1], "List.Update()")
	if err != nil {
		return value.NewNone(), err
	}

	updated := make([]value.Value, len(list))
	copy(updated, list)
	updated[index] = args[2]
	return value.Value{Type: value.ValList, Obj: updated}, nil
}

func listIndex(list []value.Value, indexValue value.Value, function string) (int, error) {
	if indexValue.Type != value.ValNumeric {
		return 0, fmt.Errorf("%s index must be Numeric, got %s", function, indexValue.Type)
	}
	if math.IsNaN(indexValue.Num) || math.IsInf(indexValue.Num, 0) || math.Trunc(indexValue.Num) != indexValue.Num {
		return 0, fmt.Errorf("%s index must be a finite whole number, got %v", function, indexValue.Num)
	}
	if indexValue.Num < 1 || indexValue.Num > float64(len(list)) {
		return 0, fmt.Errorf("%s index %v out of bounds (length %d)", function, indexValue.Num, len(list))
	}
	return int(indexValue.Num) - 1, nil
}
