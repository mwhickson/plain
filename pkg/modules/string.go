package modules

import (
	"fmt"
	"math"
	"plain/pkg/value"
	"strings"
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
		"Chr":     value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(stringChr)},
		"Trim":    value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(stringTrim)},
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
	countValue := args[1].Num
	if math.IsNaN(countValue) || math.IsInf(countValue, 0) {
		return value.NewNone(), fmt.Errorf("String.Left() count must be finite")
	}
	if countValue <= 0 {
		return value.NewString(""), nil
	}
	if countValue >= float64(len(runes)) {
		return value.NewString(string(runes)), nil
	}
	count := int(countValue)
	return value.NewString(string(runes[:count])), nil
}

func stringRight(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Type != value.ValString || args[1].Type != value.ValNumeric {
		return value.NewNone(), fmt.Errorf("String.Right() requires (String, Numeric count)")
	}
	runes := []rune(args[0].Str)
	countValue := args[1].Num
	if math.IsNaN(countValue) || math.IsInf(countValue, 0) {
		return value.NewNone(), fmt.Errorf("String.Right() count must be finite")
	}
	if countValue <= 0 {
		return value.NewString(""), nil
	}
	if countValue >= float64(len(runes)) {
		return value.NewString(string(runes)), nil
	}
	count := int(countValue)
	start := len(runes) - count
	return value.NewString(string(runes[start:])), nil
}

func stringMid(args []value.Value) (value.Value, error) {
	if len(args) < 2 || len(args) > 3 || args[0].Type != value.ValString || args[1].Type != value.ValNumeric {
		return value.NewNone(), fmt.Errorf("String.Mid() requires (String, 1-based start [, length])")
	}
	runes := []rune(args[0].Str)
	startValue := args[1].Num
	if math.IsNaN(startValue) || math.IsInf(startValue, 0) {
		return value.NewNone(), fmt.Errorf("String.Mid() start must be finite")
	}
	if startValue >= float64(len(runes)+1) {
		return value.NewString(""), nil
	}
	// 1-based start adjustment
	start := 0
	if startValue > 1 {
		start = int(startValue) - 1
	}
	if start >= len(runes) {
		return value.NewString(""), nil
	}

	if len(args) == 3 {
		if args[2].Type != value.ValNumeric {
			return value.NewNone(), fmt.Errorf("String.Mid() length must be Numeric")
		}
		lengthValue := args[2].Num
		if math.IsNaN(lengthValue) || math.IsInf(lengthValue, 0) {
			return value.NewNone(), fmt.Errorf("String.Mid() length must be finite")
		}
		if lengthValue <= 0 {
			return value.NewString(""), nil
		}
		remaining := len(runes) - start
		if lengthValue >= float64(remaining) {
			return value.NewString(string(runes[start:])), nil
		}
		end := start + int(lengthValue)
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

func stringChr(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Type != value.ValNumeric {
		return value.NewNone(), fmt.Errorf("String.Chr() requires 1 Numeric code argument")
	}
	code := int(args[0].Num)
	if code < 0 || code > 0x10FFFF {
		return value.NewNone(), fmt.Errorf("String.Chr() code out of Unicode range: %d", code)
	}
	return value.NewString(string(rune(code))), nil
}

func stringTrim(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Type != value.ValString {
		return value.NewNone(), fmt.Errorf("String.Trim() requires 1 String argument")
	}
	return value.NewString(strings.TrimSpace(args[0].Str)), nil
}
