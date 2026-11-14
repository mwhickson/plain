package modules

import (
	"bufio"
	"fmt"
	"os"
	"plain/pkg/value"
	"strings"
)

func GetConsoleModule() value.Value {
	consoleMap := map[string]value.Value{
		"Print": value.Value{
			Type: value.ValBuiltin,
			Obj: value.BuiltinFunction(func(args []value.Value) (value.Value, error) {
				var parts []string
				for _, arg := range args {
					parts = append(parts, arg.String())
				}
				fmt.Println(strings.Join(parts, " "))
				return value.NewNone(), nil
			}),
		},
		"PrintLine": value.Value{
			Type: value.ValBuiltin,
			Obj: value.BuiltinFunction(func(args []value.Value) (value.Value, error) {
				var parts []string
				for _, arg := range args {
					parts = append(parts, arg.String())
				}
				fmt.Println(strings.Join(parts, " "))
				return value.NewNone(), nil
			}),
		},
		"Input": value.Value{
			Type: value.ValBuiltin,
			Obj: value.BuiltinFunction(func(args []value.Value) (value.Value, error) {
				if len(args) > 0 {
					fmt.Print(args[0].String())
				}
				reader := bufio.NewReader(os.Stdin)
				text, err := reader.ReadString('\n')
				if err != nil {
					return value.NewString(""), err
				}
				text = strings.TrimRight(text, "\r\n")
				return value.NewString(text), nil
			}),
		},
		"Cls": value.Value{
			Type: value.ValBuiltin,
			Obj: value.BuiltinFunction(func(args []value.Value) (value.Value, error) {
				// Clear screen ANSI escape sequence
				fmt.Print("\033[H\033[2J")
				return value.NewNone(), nil
			}),
		},
	}

	return value.Value{
		Type: value.ValStruct,
		Obj:  consoleMap,
	}
}
