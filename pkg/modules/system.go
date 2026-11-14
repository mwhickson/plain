package modules

import (
	"fmt"
	"plain/pkg/value"
)

// Global reference for the registered runtime error handler function/closure
var GlobalErrorHandler value.Value

func GetSystemModule() value.Value {
	systemMap := map[string]value.Value{
		"SetErrorHandler": {
			Type: value.ValBuiltin,
			Obj: value.BuiltinFunction(func(args []value.Value) (value.Value, error) {
				if len(args) < 1 {
					return value.NewNone(), fmt.Errorf("System.SetErrorHandler expects a function argument")
				}
				fnArg := args[0]
				if fnArg.Type != value.ValFunction && fnArg.Type != value.ValClosure && fnArg.Type != value.ValBuiltin {
					return value.NewNone(), fmt.Errorf("type error: error handler must be a function, got %s", fnArg.Type)
				}
				GlobalErrorHandler = fnArg
				return value.NewNone(), nil
			}),
		},
		"Exit": {
			Type: value.ValBuiltin,
			Obj: value.BuiltinFunction(func(args []value.Value) (value.Value, error) {
				code := 0.0
				if len(args) > 0 && args[0].Type == value.ValNumeric {
					code = args[0].Num
				}
				// Can be handled or trigger exit
				return value.NewNone(), fmt.Errorf("EXIT_REQUESTED_%d", int(code))
			}),
		},
	}

	return value.Value{
		Type: value.ValStruct,
		Obj:  systemMap,
	}
}
