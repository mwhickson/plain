package modules

import (
	"fmt"
	"math"
	"plain/pkg/value"
	"strconv"
)

func GetDecimalModule() value.Value {
	decimalModule := map[string]value.Value{
		"Round":  value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(decimalRound)},
		"Format": value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(decimalFormat)},
	}
	return value.Value{
		Type: value.ValStruct,
		Obj:  decimalModule,
	}
}

func decimalRound(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Type != value.ValNumeric || args[1].Type != value.ValNumeric {
		return value.NewNone(), fmt.Errorf("Decimal.Round() requires (Numeric value, Numeric places)")
	}
	val := args[0].Num
	placesValue := args[1].Num
	if math.IsNaN(placesValue) || math.IsInf(placesValue, 0) || math.Trunc(placesValue) != placesValue || placesValue < -308 || placesValue > 308 {
		return value.NewNone(), fmt.Errorf("Decimal.Round() places must be a whole number from -308 to 308")
	}
	places := int(placesValue)
	shift := math.Pow(10, float64(places))
	rounded := math.Round(val*shift) / shift
	return value.NewNumeric(rounded), nil
}

func decimalFormat(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Type != value.ValNumeric || args[1].Type != value.ValNumeric {
		return value.NewNone(), fmt.Errorf("Decimal.Format() requires (Numeric value, Numeric places)")
	}
	val := args[0].Num
	placesValue := args[1].Num
	if math.IsNaN(placesValue) || math.IsInf(placesValue, 0) || math.Trunc(placesValue) != placesValue || placesValue < 0 || placesValue > 1000 {
		return value.NewNone(), fmt.Errorf("Decimal.Format() places must be a whole number from 0 to 1000")
	}
	places := int(placesValue)
	formatStr := "%." + strconv.Itoa(places) + "f"
	return value.NewString(fmt.Sprintf(formatStr, val)), nil
}
