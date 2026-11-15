package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"plain/pkg/value"
	"runtime"
)

var currentScriptPath string

func SetScriptPath(path string) {
	currentScriptPath = path
}

func GetEnvironmentModule() value.Value {
	envModule := map[string]value.Value{
		"GetPlatform":  value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(envGetPlatform)},
		"GetArch":      value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(envGetArch)},
		"NewLine":      value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(envNewLine)},
		"GetEnv":       value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(envGetEnv)},
		"SetEnv":       value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(envSetEnv)},
		"GetTempDir":   value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(envGetTempDir)},
		"GetHomeDir":   value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(envGetHomeDir)},
		"GetScriptDir": value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(envGetScriptDir)},
	}
	return value.Value{
		Type: value.ValStruct,
		Obj:  envModule,
	}
}

func envGetPlatform(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.NewNone(), fmt.Errorf("Environment.GetPlatform() takes 0 arguments")
	}
	return value.NewString(runtime.GOOS), nil
}

func envGetArch(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.NewNone(), fmt.Errorf("Environment.GetArch() takes 0 arguments")
	}
	return value.NewString(runtime.GOARCH), nil
}

func envNewLine(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.NewNone(), fmt.Errorf("Environment.NewLine() takes 0 arguments")
	}
	nl := "\n"
	if runtime.GOOS == "windows" {
		nl = "\r\n"
	}
	return value.NewString(nl), nil
}

func envGetEnv(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Type != value.ValString {
		return value.NewNone(), fmt.Errorf("Environment.GetEnv() requires 1 String name argument")
	}
	val := os.Getenv(args[0].Str)
	return value.NewString(val), nil
}

func envSetEnv(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Type != value.ValString || args[1].Type != value.ValString {
		return value.NewNone(), fmt.Errorf("Environment.GetEnv() requires (String name, String value)")
	}
	err := os.Setenv(args[0].Str, args[1].Str)
	if err != nil {
		return value.NewNone(), err
	}
	return value.NewNone(), nil
}

func envGetTempDir(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.NewNone(), fmt.Errorf("Environment.GetTempDir() takes 0 arguments")
	}
	return value.NewString(os.TempDir()), nil
}

func envGetHomeDir(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.NewNone(), fmt.Errorf("Environment.GetHomeDir() takes 0 arguments")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return value.NewString(""), nil
	}
	return value.NewString(home), nil
}

func envGetScriptDir(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.NewNone(), fmt.Errorf("Environment.GetScriptDir() takes 0 arguments")
	}
	if currentScriptPath == "" {
		wd, err := os.Getwd()
		if err != nil {
			return value.NewString("."), nil
		}
		return value.NewString(wd), nil
	}
	dir := filepath.Dir(currentScriptPath)
	return value.NewString(dir), nil
}
