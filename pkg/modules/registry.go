package modules

import "plain/pkg/value"

const UserGlobalStart byte = 30

type StandardModule struct {
	Name  string
	Value value.Value
}

func StandardModules() []StandardModule {
	return []StandardModule{
		{Name: "System", Value: GetSystemModule()},
		{Name: "Console", Value: GetConsoleModule()},
		{Name: "List", Value: GetListModule()},
		{Name: "String", Value: GetStringModule()},
		{Name: "Struct", Value: GetStructModule()},
		{Name: "Math", Value: GetMathModule()},
		{Name: "Random", Value: GetRandomModule()},
		{Name: "File", Value: GetFileModule()},
		{Name: "Time", Value: GetTimeModule()},
		{Name: "Decimal", Value: GetDecimalModule()},
		{Name: "Environment", Value: GetEnvironmentModule()},
	}
}
