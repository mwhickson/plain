package modules

import (
	"fmt"
	"plain/pkg/value"
)

func GetStructModule() value.Value {
	structModule := map[string]value.Value{
		"Keys":   value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(structKeys)},
		"HasKey": value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(structHasKey)},
	}
	return value.Value{
		Type: value.ValStruct,
		Obj:  structModule,
	}
}

func structKeys(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Type != value.ValStruct {
		return value.NewNone(), fmt.Errorf("Struct.Keys() requires 1 Struct argument")
	}
	m := args[0].Obj.(map[string]value.Value)
	keys := make([]value.Value, 0, len(m))
	for k := range m {
		keys = append(keys, value.NewString(k))
	}
	return value.Value{Type: value.ValList, Obj: keys}, nil
}

func structHasKey(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Type != value.ValStruct || args[1].Type != value.ValString {
		return value.NewNone(), fmt.Errorf("Struct.HasKey() requires (Struct, String keyName)")
	}
	m := args[0].Obj.(map[string]value.Value)
	keyName := args[1].Str
	_, exists := m[keyName]
	return value.NewBoolean(exists), nil
}