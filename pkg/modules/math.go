package modules

import (
	"fmt"
	"math"
	"plain/pkg/value"
)

func GetMathModule() value.Value {
	mathModule := map[string]value.Value{
		"Abs":   value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(mathAbs)},
		"Floor": value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(mathFloor)},
		"Ceil":  value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(mathCeil)},
		"Round": value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(mathRound)},
		"Sqrt":  value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(mathSqrt)},
		"Min":   value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(mathMin)},
		"Max":   value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(mathMax)},
		"Pow":   value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(mathPow)},
	}
	return value.Value{
		Type: value.ValStruct,
		Obj:  mathModule,
	}
}

func mathAbs(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Type != value.ValNumeric {
		return value.NewNone(), fmt.Errorf("Math.Abs() requires 1 Numeric argument")
	}
	return value.NewNumeric(math.Abs(args[0].Num)), nil
}

func mathFloor(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Type != value.ValNumeric {
		return value.NewNone(), fmt.Errorf("Math.Floor() requires 1 Numeric argument")
	}
	return value.NewNumeric(math.Floor(args[0].Num)), nil
}

func mathCeil(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Type != value.ValNumeric {
		return value.NewNone(), fmt.Errorf("Math.Ceil() requires 1 Numeric argument")
	}
	return value.NewNumeric(math.Ceil(args[0].Num)), nil
}

func mathRound(args []value.Value) (value.Value, error) {
	if len(args) < 1 || len(args) > 2 || args[0].Type != value.ValNumeric {
		return value.NewNone(), fmt.Errorf("Math.Round() requires (Numeric [, decimals])")
	}
	val := args[0].Num
	if len(args) == 2 {
		if args[1].Type != value.ValNumeric {
			return value.NewNone(), fmt.Errorf("Math.Round() decimals must be Numeric")
		}
		places := int(args[1].Num)
		shift := math.Pow(10, float64(places))
		val = math.Round(val*shift) / shift
		return value.NewNumeric(val), nil
	}
	return value.NewNumeric(math.Round(val)), nil
}

func mathSqrt(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Type != value.ValNumeric {
		return value.NewNone(), fmt.Errorf("Math.Sqrt() requires 1 Numeric argument")
	}
	if args[0].Num < 0 {
		return value.NewNone(), fmt.Errorf("Math.Sqrt() cannot calculate square root of a negative number")
	}
	return value.NewNumeric(math.Sqrt(args[0].Num)), nil
}

func mathMin(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Type != value.ValNumeric || args[1].Type != value.ValNumeric {
		return value.NewNone(), fmt.Errorf("Math.Min() requires (Numeric, Numeric)")
	}
	return value.NewNumeric(math.Min(args[0].Num, args[1].Num)), nil
}

func mathMax(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Type != value.ValNumeric || args[1].Type != value.ValNumeric {
		return value.NewNone(), fmt.Errorf("Math.Max() requires (Numeric, Numeric)")
	}
	return value.NewNumeric(math.Max(args[0].Num, args[1].Num)), nil
}

func mathPow(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Type != value.ValNumeric || args[1].Type != value.ValNumeric {
		return value.NewNone(), fmt.Errorf("Math.Pow() requires (Numeric base, Numeric exp)")
	}
	return value.NewNumeric(math.Pow(args[0].Num, args[1].Num)), nil
}