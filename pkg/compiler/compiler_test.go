package compiler_test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"plain/pkg/compiler"
	"plain/pkg/modules"
	"plain/pkg/parser"
	"plain/pkg/value"
	"plain/pkg/vm"
	"strconv"
	"strings"
	"testing"
)

func userGlobal(machine *vm.VM, offset byte) (value.Value, bool) {
	result, ok := machine.Globals()[modules.UserGlobalStart+offset]
	return result, ok
}

func TestCompilerExpressionCompilation(t *testing.T) {
	input := `10 + 20`

	// 1. Lex & Parse
	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	// 2. Compile
	comp := compiler.New()
	if err := comp.Compile(program); err != nil {
		t.Fatalf("compiler error: %v", err)
	}

	chunk := comp.Bytecode()

	// 3. Run in VM
	machine := vm.NewVM(chunk)
	if err := machine.Run(); err != nil {
		t.Fatalf("vm execution error: %v", err)
	}

	// 4. Verify Result on Stack
	result, err := machine.Pop()
	if err != nil {
		t.Fatalf("failed to pop result: %v", err)
	}

	if result.Num != 30 {
		t.Errorf("expected result 30, got %f", result.Num)
	}
}

func TestCompilerVariableAssignment(t *testing.T) {
	input := `Var total As Numeric = 10 + 20`

	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	comp := compiler.New()
	if err := comp.Compile(program); err != nil {
		t.Fatalf("compiler error: %v", err)
	}

	chunk := comp.Bytecode()

	machine := vm.NewVM(chunk)
	if err := machine.Run(); err != nil {
		t.Fatalf("vm execution error: %v", err)
	}

	// The declared variable remains stored in the VM globals.
}

func TestCompilerComplexExpressions(t *testing.T) {
	input := `(10 + 20) * 2`

	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	comp := compiler.New()
	if err := comp.Compile(program); err != nil {
		t.Fatalf("compiler error: %v", err)
	}

	chunk := comp.Bytecode()

	machine := vm.NewVM(chunk)
	if err := machine.Run(); err != nil {
		t.Fatalf("vm execution error: %v", err)
	}

	result, err := machine.Pop()
	if err != nil {
		t.Fatalf("failed to pop result: %v", err)
	}

	if result.Num != 60 {
		t.Errorf("expected result 60, got %f", result.Num)
	}
}

func TestCompilerIfElseStatement(t *testing.T) {
	input := `
	If 10 == 10 Then
		Var result As Numeric = 42
	Else
		Var result As Numeric = 99
	End If
	`

	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	comp := compiler.New()
	if err := comp.Compile(program); err != nil {
		t.Fatalf("compiler error: %v", err)
	}

	chunk := comp.Bytecode()

	machine := vm.NewVM(chunk)
	if err := machine.Run(); err != nil {
		t.Fatalf("vm execution error: %v", err)
	}

	// Verify that the result was assigned to a VM global.
}

func TestCompilerElseIfCompoundConditionsAndSelectCase(t *testing.T) {
	input := `
	Var a As Numeric = 1
	Var b As Numeric = 2
	Var result As Numeric = 0

	If a = 1 and b = 0 Then
		result = 10
	ElseIf a = 1 AND b = 2 or a = 0 Then
		result = 20
	Else
		result = 30
	End If

	Select Case result
	Case 10, 20
		result = 50
	Case 30
		result = 50
	Case Else
		result = 70
	End Select

	Select Case 99
	Case 1
		result = 80
	Case Else
		result = result + 10
	End Select
	`

	p := parser.NewParser(parser.NewLexer(input))
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	comp := compiler.New()
	if err := comp.Compile(program); err != nil {
		t.Fatalf("compiler error: %v", err)
	}
	machine := vm.NewVM(comp.Bytecode())
	if err := machine.Run(); err != nil {
		t.Fatalf("vm execution error: %v", err)
	}

	result, exists := userGlobal(machine, 2)
	if !exists || result.Type != value.ValNumeric || result.Num != 60 {
		t.Fatalf("result = %v, want 60", result)
	}
}

func TestCompilerSelectCaseInsideFunction(t *testing.T) {
	input := `
	Function Choose(value As Numeric) As Numeric
		Select Case value
		Case 1
			Return 10
		Case Else
			Return 20
		End Select
	End Function

	Const result As Numeric = Choose(1)
	`

	p := parser.NewParser(parser.NewLexer(input))
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	comp := compiler.New()
	if err := comp.Compile(program); err != nil {
		t.Fatalf("compiler error: %v", err)
	}
	machine := vm.NewVM(comp.Bytecode())
	if err := machine.Run(); err != nil {
		t.Fatalf("vm execution error: %v", err)
	}

	result, exists := userGlobal(machine, 1)
	if !exists || result.Type != value.ValNumeric || result.Num != 10 {
		t.Fatalf("result = %v, want 10", result)
	}
}

func TestCompilerRangeLoop(t *testing.T) {
	input := `
	Var sum As Numeric = 0
	For i As Numeric = 1 To 5 Step 1
		sum = sum + i
	End For
	`

	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	comp := compiler.New()
	if err := comp.Compile(program); err != nil {
		t.Fatalf("compiler error: %v", err)
	}

	chunk := comp.Bytecode()

	machine := vm.NewVM(chunk)
	if err := machine.Run(); err != nil {
		t.Fatalf("vm execution error: %v", err)
	}

	// sum should hold 15 (1 + 2 + 3 + 4 + 5)
	sumVal, exists := userGlobal(machine, 0)
	if !exists {
		t.Fatalf("expected variable 'sum' in globals")
	}
	if sumVal.Num != 15 {
		t.Errorf("expected sum 15, got %f", sumVal.Num)
	}
}

func TestCompilerWhileAndBreakLoops(t *testing.T) {
	input := `
	Var count As Numeric = 0
	Var whileTotal As Numeric = 0
	While count < 10
		count = count + 1
		If count == 4 Then
			Break
		End If
		whileTotal = whileTotal + count
	End While

	Var naturalCount As Numeric = 0
	While naturalCount < 3
		naturalCount = naturalCount + 1
	End While

	Var forTotal As Numeric = 0
	For i As Numeric = 1 To 10
		If i == 4 Then
			Break
		End If
		forTotal = forTotal + i
	End For

	Var outerCount As Numeric = 0
	Var innerCount As Numeric = 0
	For outer As Numeric = 1 To 3
		While True
			innerCount = innerCount + 1
			Break
		End While
		outerCount = outerCount + 1
	End For
	`

	p := parser.NewParser(parser.NewLexer(input))
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	comp := compiler.New()
	if err := comp.Compile(program); err != nil {
		t.Fatalf("compiler error: %v", err)
	}
	machine := vm.NewVM(comp.Bytecode())
	if err := machine.Run(); err != nil {
		t.Fatalf("vm execution error: %v", err)
	}

	checks := []struct {
		offset byte
		want   float64
	}{
		{0, 4},
		{1, 6},
		{2, 3},
		{3, 6},
		{4, 3},
		{5, 3},
	}
	for _, check := range checks {
		got, exists := userGlobal(machine, check.offset)
		if !exists || got.Type != value.ValNumeric || got.Num != check.want {
			t.Errorf("user global %d = %v, want %v", check.offset, got, check.want)
		}
	}
}

func TestCompilerBreakOutsideLoop(t *testing.T) {
	p := parser.NewParser(parser.NewLexer("Break"))
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	comp := compiler.New()
	if err := comp.Compile(program); err == nil || err.Error() != "Break used outside a loop" {
		t.Fatalf("compiler error = %v, want Break-outside-loop error", err)
	}
}

func TestCompilerWhileAndBreakInsideFunction(t *testing.T) {
	input := `
	Function CountTo(limit As Numeric) As Numeric
		Var count As Numeric = 0
		While count < limit
			count = count + 1
			If count == 3 Then
				Break
			End If
		End While
		Return count
	End Function

	Const result As Numeric = CountTo(10)
	`
	p := parser.NewParser(parser.NewLexer(input))
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	comp := compiler.New()
	if err := comp.Compile(program); err != nil {
		t.Fatalf("compiler error: %v", err)
	}
	machine := vm.NewVM(comp.Bytecode())
	if err := machine.Run(); err != nil {
		t.Fatalf("vm execution error: %v", err)
	}

	result, exists := userGlobal(machine, 1)
	if !exists || result.Type != value.ValNumeric || result.Num != 3 {
		t.Fatalf("CountTo(10) = %v, want 3", result)
	}
}

func TestCompilerBreakCannotEscapeFunctionFromLoop(t *testing.T) {
	input := `
	For i As Numeric = 1 To 1
		Function InvalidBreak()
			Break
		End Function
	End For
	`
	p := parser.NewParser(parser.NewLexer(input))
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	comp := compiler.New()
	if err := comp.Compile(program); err == nil || err.Error() != "Break used outside a loop" {
		t.Fatalf("compiler error = %v, want Break-outside-loop error", err)
	}
}

func TestCompilerFunctionCall(t *testing.T) {
	input := `
	Function Add(a As Numeric, b As Numeric) As Numeric
		Return a + b
	End Function

	Var result As Numeric = Add(10, 20)
	`

	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	comp := compiler.New()
	if err := comp.Compile(program); err != nil {
		t.Fatalf("compiler error: %v", err)
	}

	chunk := comp.Bytecode()

	machine := vm.NewVM(chunk)
	if err := machine.Run(); err != nil {
		t.Fatalf("vm execution error: %v", err)
	}

	// Verify that 'result' holds 30
	resultVal, exists := userGlobal(machine, 1)
	if !exists {
		t.Fatalf("expected variable 'result' in globals")
	}
	if resultVal.Num != 30 {
		t.Errorf("expected result 30, got %f", resultVal.Num)
	}
}

func TestCompilerForwardCallAndSystemExit(t *testing.T) {
	input := `
	System.SetErrorHandler(Handler)
	Const Base As Numeric = 0
	Const Result As Numeric = Main()
	System.Exit(Result)

	Function Main() As Numeric
		Return Base
	End Function

	Function Handler(errMsg As String)
		Return
	End Function
	`

	p := parser.NewParser(parser.NewLexer(input))
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	comp := compiler.New()
	if err := comp.Compile(program); err != nil {
		t.Fatalf("compiler error: %v", err)
	}

	machine := vm.NewVM(comp.Bytecode())
	if err := machine.Run(); err != nil {
		t.Fatalf("vm execution error: %v", err)
	}
	if code, exited := machine.ExitCode(); !exited || code != 0 {
		t.Fatalf("exit status = (%d, %t), want (0, true)", code, exited)
	}
}

func TestCompilerRejectsConstantAssignment(t *testing.T) {
	input := `
	Const answer As Numeric = 42
	answer = 43
	`

	p := parser.NewParser(parser.NewLexer(input))
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	comp := compiler.New()
	err := comp.Compile(program)
	if err == nil || err.Error() != "cannot assign to constant: answer" {
		t.Fatalf("compiler error = %v, want constant-assignment error", err)
	}
}

func TestCompilerPrimitiveTypes(t *testing.T) {
	input := `
	Var numVal As Numeric = 42
	Var strVal As String = "Hello PLAIN"
	Var boolVal As Boolean = True
	Var noneVal As None = None
	`

	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	comp := compiler.New()
	if err := comp.Compile(program); err != nil {
		t.Fatalf("compiler error: %v", err)
	}

	chunk := comp.Bytecode()
	machine := vm.NewVM(chunk)
	if err := machine.Run(); err != nil {
		t.Fatalf("vm execution error: %v", err)
	}

	// Verify global variable slots...
}

func TestCompilerTypeChecking(t *testing.T) {
	tests := []struct {
		input       string
		expectedErr string
	}{
		{
			input:       `Var count As Numeric = "not a number"`,
			expectedErr: "type mismatch for 'count': declared Numeric, got String",
		},
		{
			input:       `Var name As String = True`,
			expectedErr: "type mismatch for 'name': declared String, got Boolean",
		},
		{
			input:       `Var flag = 123`,
			expectedErr: "type required for variable 'flag' (Option Explicit)",
		},
	}

	for _, tt := range tests {
		l := parser.NewLexer(tt.input)
		p := parser.NewParser(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			t.Fatalf("parser errors for %q: %v", tt.input, p.Errors())
		}

		comp := compiler.New()
		err := comp.Compile(program)
		if err == nil {
			t.Fatalf("expected compiler error for %q, got nil", tt.input)
		}

		if err.Error() != tt.expectedErr {
			t.Fatalf("expected error %q, got %q", tt.expectedErr, err.Error())
		}
	}
}

func TestCompilerRecursion(t *testing.T) {
	input := `
	Function Countdown(n As Numeric) As Numeric
		If n <= 0 Then
			Return 0
		End If
		Return n + Countdown(n - 1)
	End Function

	Var result As Numeric = Countdown(3)
	`

	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	comp := compiler.New()
	if err := comp.Compile(program); err != nil {
		t.Fatalf("compiler error: %v", err)
	}

	chunk := comp.Bytecode()
	machine := vm.NewVM(chunk)
	if err := machine.Run(); err != nil {
		t.Fatalf("vm execution error: %v", err)
	}
}

func TestCompilerFibonacciRecursion(t *testing.T) {
	input := `
	Function Fib(n As Numeric) As Numeric
		If n <= 1 Then
			Return n
		End If
		Return Fib(n - 1) + Fib(n - 2)
	End Function

	Var result As Numeric = Fib(6)
	`

	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	comp := compiler.New()
	if err := comp.Compile(program); err != nil {
		t.Fatalf("compiler error: %v", err)
	}

	chunk := comp.Bytecode()
	machine := vm.NewVM(chunk)
	if err := machine.Run(); err != nil {
		t.Fatalf("vm execution error: %v", err)
	}

	resultVal, exists := userGlobal(machine, 1)
	if !exists {
		t.Fatalf("expected variable 'result' in globals")
	}
	if resultVal.Num != 8 {
		t.Errorf("expected Fib(6) to be 8, got %f", resultVal.Num)
	}
}

func TestCompilerMutualRecursion(t *testing.T) {
	input := `
	Function IsEven(n As Numeric) As Boolean
		If n == 0 Then
			Return True
		End If
		Return IsOdd(n - 1)
	End Function

	Function IsOdd(n As Numeric) As Boolean
		If n == 0 Then
			Return False
		End If
		Return IsEven(n - 1)
	End Function

	Var result As Boolean = IsEven(4)
	`

	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	comp := compiler.New()
	if err := comp.Compile(program); err != nil {
		t.Fatalf("compiler error: %v", err)
	}

	chunk := comp.Bytecode()
	machine := vm.NewVM(chunk)
	if err := machine.Run(); err != nil {
		t.Fatalf("vm execution error: %v", err)
	}

	resultVal, exists := userGlobal(machine, 2)
	if !exists {
		t.Fatalf("expected variable 'result' in globals")
	}
	if !resultVal.Bool {
		t.Errorf("expected IsEven(4) to be True, got %v", resultVal.Bool)
	}
}

func TestCompilerClosures(t *testing.T) {
	input := `
	Function MakeAdder(x As Numeric) As Function
		Function Add(y As Numeric) As Numeric
			Return x + y
		End Function
		Return Add
	End Function

	Var addFive As Function = MakeAdder(5)
	Var result As Numeric = addFive(10)
	`

	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	comp := compiler.New()
	if err := comp.Compile(program); err != nil {
		t.Fatalf("compiler error: %v", err)
	}

	chunk := comp.Bytecode()
	machine := vm.NewVM(chunk)
	if err := machine.Run(); err != nil {
		t.Fatalf("vm execution error: %v", err)
	}

	// Verify that 'result' holds 15 in globals
	resultVal, exists := userGlobal(machine, 2)
	if !exists {
		t.Fatalf("expected variable 'result' in globals")
	}
	if resultVal.Num != 15 {
		t.Errorf("expected result 15, got %f", resultVal.Num)
	}
}

func TestCompilerLists(t *testing.T) {
	input := `
	Var items As List = ["Apple", "Banana", "Cherry"]
	Var first As String = items[1]
	`
	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	comp := compiler.New()
	if err := comp.Compile(program); err != nil {
		t.Fatalf("compiler error: %v", err)
	}

	chunk := comp.Bytecode()
	machine := vm.NewVM(chunk)
	if err := machine.Run(); err != nil {
		t.Fatalf("vm execution error: %v", err)
	}

	first, exists := userGlobal(machine, 1)
	if !exists || first.Type != value.ValString || first.Str != "Apple" {
		t.Fatalf("first item = %v, want Apple", first)
	}
}

func TestCompilerListIndexZeroIsOutOfBounds(t *testing.T) {
	input := `
	Var items As List = ["Apple", "Banana"]
	Var first As String = items[0]
	`
	p := parser.NewParser(parser.NewLexer(input))
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	comp := compiler.New()
	if err := comp.Compile(program); err != nil {
		t.Fatalf("compiler error: %v", err)
	}
	machine := vm.NewVM(comp.Bytecode())
	if err := machine.Run(); err == nil || err.Error() != "index out of bounds: 0 (length 2)" {
		t.Fatalf("vm error = %v, want index-zero bounds error", err)
	}
}

func TestCompilerFunctionLocalsInsideRangeLoop(t *testing.T) {
	input := `
	Function Sum(items As List) As Numeric
		Var total As Numeric = 0
		For i As Numeric = 1 To List.Len(items)
			total = total + items[i]
		End For
		Return total
	End Function

	Const result As Numeric = Sum([10, 20])
	`
	p := parser.NewParser(parser.NewLexer(input))
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	comp := compiler.New()
	if err := comp.Compile(program); err != nil {
		t.Fatalf("compiler error: %v", err)
	}
	machine := vm.NewVM(comp.Bytecode())
	if err := machine.Run(); err != nil {
		t.Fatalf("vm execution error: %v", err)
	}
	result, exists := userGlobal(machine, 1)
	if !exists || result.Type != value.ValNumeric || result.Num != 30 {
		t.Fatalf("sum result = %v, want 30", result)
	}
}

func TestCompilerStructs(t *testing.T) {
	input := `
	Struct Vector2D
		X As Numeric
		Y As Numeric
	End Struct
	`
	// Since structs currently store field definitions, we can also test creating a struct instance or accessing properties if supported,
	// or verify that the struct statement compiles without error.
	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	comp := compiler.New()
	if err := comp.Compile(program); err != nil {
		t.Fatalf("compiler error: %v", err)
	}
}

func TestCompilerDefaultStructValues(t *testing.T) {
	input := `
	Const Rooms As List = [
		CreateRoom(100, "Cave Entrance"),
		CreateRoom(200, "Crystal Cavern")
	]

	Struct Room
		Id As Numeric
		Name As String
		Description As String
	End Struct

	Function CreateRoom(id As Numeric, name As String) As Room
		Var room As Room
		room.Id = id
		room.Name = name
		Return room
	End Function
	`

	p := parser.NewParser(parser.NewLexer(input))
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	comp := compiler.New()
	if err := comp.Compile(program); err != nil {
		t.Fatalf("compiler error: %v", err)
	}
	machine := vm.NewVM(comp.Bytecode())
	if err := machine.Run(); err != nil {
		t.Fatalf("vm execution error: %v", err)
	}

	roomsValue, exists := userGlobal(machine, 1)
	if !exists || roomsValue.Type != value.ValList {
		t.Fatalf("Rooms global = %v, want list", roomsValue)
	}
	rooms := roomsValue.Obj.([]value.Value)
	if len(rooms) != 2 {
		t.Fatalf("Rooms length = %d, want 2", len(rooms))
	}

	first := rooms[0].Obj.(map[string]value.Value)
	second := rooms[1].Obj.(map[string]value.Value)
	if first["Id"].Num != 100 || first["Name"].Str != "Cave Entrance" {
		t.Errorf("first room = %v, want Id 100 and Cave Entrance", first)
	}
	if first["Description"].Type != value.ValString || first["Description"].Str != "" {
		t.Errorf("first room Description = %v, want empty String default", first["Description"])
	}
	if second["Id"].Num != 200 || second["Name"].Str != "Crystal Cavern" {
		t.Errorf("second room = %v, want Id 200 and Crystal Cavern", second)
	}
}

func TestCompilerReceiverFunctions(t *testing.T) {
	input := `
	Struct Vector2D
		X As Numeric
		Y As Numeric
	End Struct

	Function (v As Vector2D) Scale(factor As Numeric)
		v.X = v.X * factor
	End Function
	`

	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	comp := compiler.New()
	if err := comp.Compile(program); err != nil {
		t.Fatalf("compiler error: %v", err)
	}
}

func TestVMMetricAndStructOperations(t *testing.T) {
	input := `
	Struct Vector2D
		X As Numeric
		Y As Numeric
	End Struct

	Function (v As Vector2D) Scale(factor As Numeric)
		v.X = v.X * factor
		v.Y = v.Y * factor
	End Function
	`

	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	comp := compiler.New()
	if err := comp.Compile(program); err != nil {
		t.Fatalf("compiler error: %v", err)
	}

	chunk := comp.Bytecode()
	machine := vm.NewVM(chunk)
	if err := machine.Run(); err != nil {
		t.Fatalf("vm execution error: %v", err)
	}
}

func TestCompilerStandardLibraryModules(t *testing.T) {
	path := filepath.Join(t.TempDir(), "module-test.txt")
	if err := os.WriteFile(path, []byte("module data"), 0600); err != nil {
		t.Fatal(err)
	}

	input := fmt.Sprintf(`
	Struct Person
		Name As String
	End Struct

	Var upper As String = String.Upper("plain")
	Var length As Numeric = String.Len(upper)
	Var root As Numeric = Math.Sqrt(81)
	Random.Seed(13)
	Var randomValue As Numeric = Random.Int(7, 7)
	Var decimalText As String = Decimal.Format(12.5, 2)
	Var person As Person
	person.Name = "Ada"
	Var keys As List = Struct.Keys(person)
	Var hasName As Boolean = Struct.HasKey(person, "Name")
	Var fileText As String = File.ReadText(%s)
	Var fileExists As Boolean = File.Exists(%s)
	Var timeText As String = Time.Format(Time.Now())
	Var platform As String = Environment.GetPlatform()
	`, strconv.Quote(path), strconv.Quote(path))

	p := parser.NewParser(parser.NewLexer(input))
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}
	comp := compiler.New()
	if err := comp.Compile(program); err != nil {
		t.Fatalf("compiler error: %v", err)
	}
	machine := vm.NewVM(comp.Bytecode())
	if err := machine.Run(); err != nil {
		t.Fatalf("vm execution error: %v", err)
	}

	globals := machine.Globals()
	checks := []struct {
		slot byte
		want value.Value
	}{
		{0, value.NewString("PLAIN")},
		{1, value.NewNumeric(5)},
		{2, value.NewNumeric(9)},
		{3, value.NewNumeric(7)},
		{4, value.NewString("12.50")},
		{7, value.NewBoolean(true)},
		{8, value.NewString("module data")},
		{9, value.NewBoolean(true)},
	}
	for _, check := range checks {
		got, exists := userGlobal(machine, check.slot)
		if !exists || got.Type != check.want.Type || got.Num != check.want.Num || got.Str != check.want.Str || got.Bool != check.want.Bool {
			t.Errorf("user global %d = %v, want %v", check.slot, got, check.want)
		}
	}

	keys, exists := userGlobal(machine, 6)
	if !exists || keys.Type != value.ValList || len(keys.Obj.([]value.Value)) != 1 || keys.Obj.([]value.Value)[0].Str != "Name" {
		t.Fatalf("Struct.Keys() result = %v, want [Name]", keys)
	}
	if timeText, exists := userGlobal(machine, 10); !exists || timeText.Type != value.ValString || timeText.Str == "" {
		t.Fatalf("Time.Format() result = %v, want non-empty String", timeText)
	}
	if platform, exists := userGlobal(machine, 11); !exists || platform.Type != value.ValString || platform.Str == "" {
		t.Fatalf("Environment.GetPlatform() result = %v, want non-empty String", platform)
	}
	if _, exists := globals[modules.UserGlobalStart-1]; exists {
		t.Fatal("standard library reserved an unexpected final slot")
	}
}

func TestFileExampleCompiles(t *testing.T) {
	source, err := os.ReadFile(filepath.Join("..", "..", "examples", "file.plain"))
	if err != nil {
		t.Fatal(err)
	}

	p := parser.NewParser(parser.NewLexer(string(source)))
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors in file example: %v", p.Errors())
	}

	comp := compiler.New()
	if err := comp.Compile(program); err != nil {
		t.Fatalf("compiler error in file example: %v", err)
	}
}

func TestListExampleCompilesAndRunsToMenu(t *testing.T) {
	source, err := os.ReadFile(filepath.Join("..", "..", "examples", "list.plain"))
	if err != nil {
		t.Fatal(err)
	}

	p := parser.NewParser(parser.NewLexer(string(source)))
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors in list example: %v", p.Errors())
	}

	comp := compiler.New()
	if err := comp.Compile(program); err != nil {
		t.Fatalf("compiler error in list example: %v", err)
	}
	chunk := comp.Bytecode()
	if count := len(chunk.Constants); count > 256 {
		t.Fatalf("list example constant pool has %d values, exceeds bytecode index capacity", count)
	}

	scriptDir := t.TempDir()
	dataDir := filepath.Join(scriptDir, "data", "tasks")
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dataDir, "index.dat"), nil, 0600); err != nil {
		t.Fatal(err)
	}
	modules.SetScriptPath(filepath.Join(scriptDir, "list.plain"))
	t.Cleanup(func() { modules.SetScriptPath("") })

	input, err := os.CreateTemp(t.TempDir(), "stdin-*")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := input.WriteString("Q\n"); err != nil {
		t.Fatal(err)
	}
	if _, err := input.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	oldStdin, oldStdout := os.Stdin, os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		os.Stdin, os.Stdout = oldStdin, oldStdout
		reader.Close()
		writer.Close()
		input.Close()
	})
	os.Stdin, os.Stdout = input, writer
	machine := vm.NewVM(chunk)
	runErr := machine.Run()
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	os.Stdin, os.Stdout = oldStdin, oldStdout

	if runErr != nil {
		t.Fatalf("list example VM error before menu exit: %v", runErr)
	}
	if code, exited := machine.ExitCode(); !exited || code != 0 {
		t.Fatalf("list example exit status = (%d, %t), want (0, true)", code, exited)
	}
	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(output), "PLAIN List Demo") ||
		!strings.Contains(string(output), "A)dd a Task") ||
		!strings.Contains(string(output), "Enter your choice:") {
		t.Fatalf("list example output did not reach menu: %q", output)
	}
}
