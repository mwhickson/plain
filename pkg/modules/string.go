package modules

import (
	"fmt"
	"strings"
	"plain/pkg/value"
)

func GetStringModule() value.Value {
	stringModule := map[string]value.Value{
		"Len":     value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(stringLen)},
		"Lower":   value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(stringLower)},
		"Upper":   value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(stringUpper)},
		"Left":    value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(stringLeft)},
		"Right":   value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(stringRight)},
		"Mid":     value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(stringMid)},
		"Split":   value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(stringSplit)},
		"Join":    value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(stringJoin)},
		"Replace": value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(stringReplace)},
	}
	return value.Value{
		Type: value.ValStruct,
		Obj:  stringModule,
	}
}

func stringLen(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Type != value.ValString {
		return value.NewNone(), fmt.Errorf("String.Len() requires 1 String argument")
	}
	runes := []rune(args[0].Str)
	return value.NewNumeric(float64(len(runes))), nil
}

func stringLower(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Type != value.ValString {
		return value.NewNone(), fmt.Errorf("String.Lower() requires 1 String argument")
	}
	return value.NewString(strings.ToLower(args[0].Str)), nil
}

func stringUpper(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Type != value.ValString {
		return value.NewNone(), fmt.Errorf("String.Upper() requires 1 String argument")
	}
	return value.NewString(strings.ToUpper(args[0].Str)), nil
}

func stringLeft(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Type != value.ValString || args[1].Type != value.ValNumeric {
		return value.NewNone(), fmt.Errorf("String.Left() requires (String, Numeric count)")
	}
	runes := []rune(args[0].Str)
	count := int(args[1].Num)
	if count < 0 {
		count = 0
	}
	if count > len(runes) {
		count = len(runes)
	}
	return value.NewString(string(runes[:count])), nil
}

func stringRight(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Type != value.ValString || args[1].Type != value.ValNumeric {
		return value.NewNone(), fmt.Errorf("String.Right() requires (String, Numeric count)")
	}
	runes := []rune(args[0].Str)
	count := int(args[1].Num)
	if count < 0 {
		count = 0
	}
	if count > len(runes) {
		count = len(runes)
	}
	start := len(runes) - count
	return value.NewString(string(runes[start:])), nil
}

func stringMid(args []value.Value) (value.Value, error) {
	if len(args) < 2 || len(args) > 3 || args[0].Type != value.ValString || args[1].Type != value.ValNumeric {
		return value.NewNone(), fmt.Errorf("String.Mid() requires (String, 1-based start [, length])")
	}
	runes := []rune(args[0].Str)
	// 1-based start adjustment
	start := int(args[1].Num) - 1
	if start < 0 {
		start = 0
	}
	if start >= len(runes) {
		return value.NewString(""), nil
	}

	if len(args) == 3 {
		if args[2].Type != value.ValNumeric {
			return value.NewNone(), fmt.Errorf("String.Mid() length must be Numeric")
		}
		length := int(args[2].Num)
		end := start + length
		if end > len(runes) {
			end = len(runes)
		}
		return value.NewString(string(runes[start:end])), nil
	}

	return value.NewString(string(runes[start:])), nil
}

func stringSplit(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Type != value.ValString || args[1].Type != value.ValString {
		return value.NewNone(), fmt.Errorf("String.Split() requires (String, delimiter String)")
	}
	parts := strings.Split(args[0].Str, args[1].Str)
	listVal := make([]value.Value, len(parts))
	for i, p := range parts {
		listVal[i] = value.NewString(p)
	}
	return value.Value{Type: value.ValList, Obj: listVal}, nil
}

func stringJoin(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Type != value.ValList || args[1].Type != value.ValString {
		return value.NewNone(), fmt.Errorf("String.Join() requires (List, delimiter String)")
	}
	list := args[0].Obj.([]value.Value)
	parts := make([]string, len(list))
	for i, item := range list {
		parts[i] = item.String()
	}
	return value.NewString(strings.Join(parts, args[1].Str)), nil
}

func stringReplace(args []value.Value) (value.Value, error) {
	if len(args) != 3 || args[0].Type != value.ValString || args[1].Type != value.ValString || args[2].Type != value.ValString {
		return value.NewNone(), fmt.Errorf("String.Replace() requires (String, target, replacement)")
	}
	res := strings.ReplaceAll(args[0].Str, args[1].Str, args[2].Str)
	return value.NewString(res), nil
}