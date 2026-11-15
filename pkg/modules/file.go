package modules

import (
	"fmt"
	"os"
	"plain/pkg/value"
)

func GetFileModule() value.Value {
	fileModule := map[string]value.Value{
		"ReadText":   value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(fileReadText)},
		"WriteText":  value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(fileWriteText)},
		"AppendText": value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(fileAppendText)},
		"Exists":     value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(fileExists)},
		"ReadBytes":  value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(fileReadBytes)},
		"WriteBytes": value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(fileWriteBytes)},
	}
	return value.Value{
		Type: value.ValStruct,
		Obj:  fileModule,
	}
}

func fileReadText(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Type != value.ValString {
		return value.NewNone(), fmt.Errorf("File.ReadText() requires 1 String path argument")
	}
	data, err := os.ReadFile(args[0].Str)
	if err != nil {
		return value.NewNone(), err
	}
	return value.NewString(string(data)), nil
}

func fileWriteText(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Type != value.ValString || args[1].Type != value.ValString {
		return value.NewNone(), fmt.Errorf("File.WriteText() requires (String path, String content)")
	}
	err := os.WriteFile(args[0].Str, []byte(args[1].Str), 0644)
	if err != nil {
		return value.NewNone(), err
	}
	return value.NewNone(), nil
}

func fileAppendText(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Type != value.ValString || args[1].Type != value.ValString {
		return value.NewNone(), fmt.Errorf("File.AppendText() requires (String path, String content)")
	}
	f, err := os.OpenFile(args[0].Str, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return value.NewNone(), err
	}
	defer f.Close()

	_, err = f.WriteString(args[1].Str)
	if err != nil {
		return value.NewNone(), err
	}
	return value.NewNone(), nil
}

func fileExists(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Type != value.ValString {
		return value.NewNone(), fmt.Errorf("File.Exists() requires 1 String path argument")
	}
	_, err := os.Stat(args[0].Str)
	return value.NewBoolean(err == nil), nil
}

func fileReadBytes(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Type != value.ValString {
		return value.NewNone(), fmt.Errorf("File.ReadBytes() requires 1 String path argument")
	}
	data, err := os.ReadFile(args[0].Str)
	if err != nil {
		return value.NewNone(), err
	}
	list := make([]value.Value, len(data))
	for i, b := range data {
		list[i] = value.NewNumeric(float64(b))
	}
	return value.Value{Type: value.ValList, Obj: list}, nil
}

func fileWriteBytes(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Type != value.ValString || args[1].Type != value.ValList {
		return value.NewNone(), fmt.Errorf("File.WriteBytes() requires (String path, List byteValues)")
	}
	list := args[1].Obj.([]value.Value)
	data := make([]byte, len(list))
	for i, item := range list {
		if item.Type != value.ValNumeric {
			return value.NewNone(), fmt.Errorf("File.WriteBytes() list elements must be Numeric bytes")
		}
		data[i] = byte(item.Num)
	}
	err := os.WriteFile(args[0].Str, data, 0644)
	if err != nil {
		return value.NewNone(), err
	}
	return value.NewNone(), nil
}