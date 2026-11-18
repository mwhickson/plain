package modules

import (
	"io"
	"math"
	"os"
	"path/filepath"
	"plain/pkg/value"
	"reflect"
	"testing"
)

func TestConsolePrintDoesNotAddNewline(t *testing.T) {
	oldStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = writer
	defer func() {
		os.Stdout = oldStdout
		writer.Close()
		reader.Close()
	}()

	module := GetConsoleModule()
	callModule(t, module, "Print", value.NewString("hello"))
	callModule(t, module, "PrintLine", value.NewString("world"))

	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	os.Stdout = oldStdout
	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(output), "helloworld\n"; got != want {
		t.Fatalf("console output = %q, want %q", got, want)
	}
}

func callModule(t *testing.T, module value.Value, name string, args ...value.Value) value.Value {
	t.Helper()
	exports, ok := module.Obj.(map[string]value.Value)
	if !ok {
		t.Fatal("module does not contain an export map")
	}
	export, ok := exports[name]
	if !ok {
		t.Fatalf("module does not export %q", name)
	}
	fn, ok := export.Obj.(value.BuiltinFunction)
	if !ok {
		t.Fatalf("module export %q is not a builtin", name)
	}
	result, err := fn(args)
	if err != nil {
		t.Fatalf("%s failed: %v", name, err)
	}
	return result
}

func TestStandardModulesFitReservedGlobalRange(t *testing.T) {
	modules := StandardModules()
	if len(modules) > int(UserGlobalStart) {
		t.Fatalf("%d standard modules exceed reserved global range ending at %d", len(modules), UserGlobalStart-1)
	}
	seen := make(map[string]bool, len(modules))
	for i, module := range modules {
		if module.Name == "" || seen[module.Name] {
			t.Fatalf("standard module %d has an empty or duplicate name %q", i, module.Name)
		}
		seen[module.Name] = true
		if module.Value.Type != value.ValStruct {
			t.Errorf("standard module %q has type %s, want Struct", module.Name, module.Value.Type)
		}
	}
}

func TestStringModuleUnicodeAndMidBounds(t *testing.T) {
	module := GetStringModule()
	input := value.NewString("héllo")

	if got := callModule(t, module, "Len", input); got.Num != 5 {
		t.Fatalf("String.Len() = %v, want 5", got.Num)
	}
	if got := callModule(t, module, "Left", input, value.NewNumeric(2)); got.Str != "hé" {
		t.Fatalf("String.Left() = %q, want hé", got.Str)
	}
	if got := callModule(t, module, "Right", input, value.NewNumeric(2)); got.Str != "lo" {
		t.Fatalf("String.Right() = %q, want lo", got.Str)
	}
	if got := callModule(t, module, "Mid", input, value.NewNumeric(2), value.NewNumeric(-1)); got.Str != "" {
		t.Fatalf("String.Mid() with negative length = %q, want empty string", got.Str)
	}
	if got := callModule(t, module, "ToNumeric", value.NewString(" 12.5 ")); got.Type != value.ValNumeric || got.Num != 12.5 {
		t.Fatalf("String.ToNumeric() = %v, want 12.5", got)
	}
	if got := callModule(t, module, "ToNumeric", value.NewString("not numeric")); got.Type != value.ValNumeric || got.Num != 0 {
		t.Fatalf("String.ToNumeric() for invalid input = %v, want 0", got)
	}
	if got := callModule(t, module, "FromNumeric", value.NewNumeric(-12.5)); got.Type != value.ValString || got.Str != "-12.5" {
		t.Fatalf("String.FromNumeric() = %v, want -12.5", got)
	}
}

func TestStructModuleKeysAreStable(t *testing.T) {
	module := GetStructModule()
	instance := value.Value{Type: value.ValStruct, Obj: map[string]value.Value{
		"z": value.NewNumeric(1),
		"a": value.NewNumeric(2),
	}}

	keys := callModule(t, module, "Keys", instance).Obj.([]value.Value)
	if !reflect.DeepEqual(keys, []value.Value{value.NewString("a"), value.NewString("z")}) {
		t.Fatalf("Struct.Keys() = %v, want [a z]", keys)
	}
	if got := callModule(t, module, "HasKey", instance, value.NewString("z")); !got.Bool {
		t.Fatal("Struct.HasKey() did not find existing key")
	}
}

func TestListModuleImmutableItemOperations(t *testing.T) {
	module := GetListModule()
	backing := make([]value.Value, 2, 3)
	backing[0] = value.NewString("first")
	backing[1] = value.NewString("second")
	original := value.Value{Type: value.ValList, Obj: backing}

	appended := callModule(t, module, "Append", original, value.NewString("third"))
	updated := callModule(t, module, "Update", original, value.NewNumeric(2), value.NewString("changed"))
	deleted := callModule(t, module, "Delete", original, value.NewNumeric(1))

	if got := appended.Obj.([]value.Value); !reflect.DeepEqual(got, []value.Value{
		value.NewString("first"), value.NewString("second"), value.NewString("third"),
	}) {
		t.Fatalf("List.Append() = %v, want [first second third]", got)
	}
	if got := updated.Obj.([]value.Value); !reflect.DeepEqual(got, []value.Value{
		value.NewString("first"), value.NewString("changed"),
	}) {
		t.Fatalf("List.Update() = %v, want [first changed]", got)
	}
	if got := deleted.Obj.([]value.Value); !reflect.DeepEqual(got, []value.Value{
		value.NewString("second"),
	}) {
		t.Fatalf("List.Delete() = %v, want [second]", got)
	}
	if got := original.Obj.([]value.Value); !reflect.DeepEqual(got, []value.Value{
		value.NewString("first"), value.NewString("second"),
	}) {
		t.Fatalf("list operation changed its input: %v", got)
	}
	if backing[:3][2].Type != value.ValNone {
		t.Fatalf("List.Append() modified the input list's spare capacity: %v", backing[:3][2])
	}
}

func TestListUpdateAndDeleteRejectInvalidIndices(t *testing.T) {
	module := GetListModule()
	exports := module.Obj.(map[string]value.Value)
	list := value.Value{Type: value.ValList, Obj: []value.Value{value.NewString("item")}}
	tests := []struct {
		name      string
		operation string
		args      []value.Value
	}{
		{"Delete zero", "Delete", []value.Value{list, value.NewNumeric(0)}},
		{"Delete fractional", "Delete", []value.Value{list, value.NewNumeric(1.5)}},
		{"Delete out of range", "Delete", []value.Value{list, value.NewNumeric(2)}},
		{"Update non-numeric", "Update", []value.Value{list, value.NewString("1"), value.NewString("new")}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fn := exports[test.operation].Obj.(value.BuiltinFunction)
			if _, err := fn(test.args); err == nil {
				t.Fatalf("List.%s() accepted invalid index arguments %v", test.operation, test.args)
			}
		})
	}
}

func TestMathAndDecimalModules(t *testing.T) {
	mathModule := GetMathModule()
	if got := callModule(t, mathModule, "Sqrt", value.NewNumeric(9)); got.Num != 3 {
		t.Fatalf("Math.Sqrt(9) = %v, want 3", got.Num)
	}
	if _, err := mathSqrt([]value.Value{value.NewNumeric(-1)}); err == nil {
		t.Fatal("Math.Sqrt() accepted a negative number")
	}

	decimalModule := GetDecimalModule()
	if got := callModule(t, decimalModule, "Round", value.NewNumeric(123.456), value.NewNumeric(2)); got.Num != 123.46 {
		t.Fatalf("Decimal.Round() = %v, want 123.46", got.Num)
	}
	if got := callModule(t, decimalModule, "Format", value.NewNumeric(12.5), value.NewNumeric(2)); got.Str != "12.50" {
		t.Fatalf("Decimal.Format() = %q, want 12.50", got.Str)
	}
	if _, err := decimalFormat([]value.Value{value.NewNumeric(1), value.NewNumeric(-1)}); err == nil {
		t.Fatal("Decimal.Format() accepted negative precision")
	}
}

func TestRandomModuleSeedAndBounds(t *testing.T) {
	module := GetRandomModule()
	callModule(t, module, "Seed", value.NewNumeric(42))
	first := callModule(t, module, "Float")
	callModule(t, module, "Seed", value.NewNumeric(42))
	second := callModule(t, module, "Float")
	if first.Num != second.Num {
		t.Fatalf("seeded Random.Float() values differ: %v and %v", first.Num, second.Num)
	}

	got := callModule(t, module, "Int", value.NewNumeric(5), value.NewNumeric(1))
	if got.Num < 1 || got.Num > 5 {
		t.Fatalf("Random.Int(5, 1) = %v, want a value from 1 through 5", got.Num)
	}
	if _, err := randomInt([]value.Value{value.NewNumeric(-4611686018427387904), value.NewNumeric(4611686018427387904)}); err == nil {
		t.Fatal("Random.Int() accepted a range too large for its generator")
	}
	if _, err := randomChance([]value.Value{value.NewNumeric(101)}); err == nil {
		t.Fatal("Random.Chance() accepted a percentage above 100")
	}
}

func TestFileModuleTextBytesAndValidation(t *testing.T) {
	module := GetFileModule()
	path := filepath.Join(t.TempDir(), "payload.bin")
	pathValue := value.NewString(path)

	callModule(t, module, "WriteText", pathValue, value.NewString("hello"))
	if got := callModule(t, module, "ReadText", pathValue); got.Str != "hello" {
		t.Fatalf("File.ReadText() = %q, want hello", got.Str)
	}
	callModule(t, module, "AppendText", pathValue, value.NewString(" world"))
	if got := callModule(t, module, "ReadText", pathValue); got.Str != "hello world" {
		t.Fatalf("File.AppendText() result = %q, want hello world", got.Str)
	}

	bytes := value.Value{Type: value.ValList, Obj: []value.Value{value.NewNumeric(0), value.NewNumeric(128), value.NewNumeric(255)}}
	callModule(t, module, "WriteBytes", pathValue, bytes)
	gotBytes := callModule(t, module, "ReadBytes", pathValue).Obj.([]value.Value)
	if !reflect.DeepEqual(gotBytes, []value.Value{value.NewNumeric(0), value.NewNumeric(128), value.NewNumeric(255)}) {
		t.Fatalf("File.ReadBytes() = %v, want byte values [0 128 255]", gotBytes)
	}
	if _, err := fileWriteBytes([]value.Value{pathValue, {Type: value.ValList, Obj: []value.Value{value.NewNumeric(256)}}}); err == nil {
		t.Fatal("File.WriteBytes() accepted a value above 255")
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("test payload disappeared: %v", err)
	}
	if got := callModule(t, module, "Exists", value.NewString(filepath.Join(filepath.Dir(path), "missing"))); got.Bool {
		t.Fatal("File.Exists() returned true for a missing file")
	}
}

func TestTimeModuleRejectsInvalidSleep(t *testing.T) {
	module := GetTimeModule()
	if got := callModule(t, module, "Now"); got.Type != value.ValNumeric {
		t.Fatalf("Time.Now() returned %s, want Numeric", got.Type)
	}
	if _, err := timeSleep([]value.Value{value.NewNumeric(-1)}); err == nil {
		t.Fatal("Time.Sleep() accepted a negative duration")
	}
	if _, err := timeFormat([]value.Value{value.NewNumeric(math.Inf(1))}); err == nil {
		t.Fatal("Time.Format() accepted a non-finite timestamp")
	}
	callModule(t, module, "Sleep", value.NewNumeric(0))
}
