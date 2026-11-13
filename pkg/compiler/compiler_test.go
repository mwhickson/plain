package compiler_test

import (
	"plain/pkg/compiler"
	"plain/pkg/parser"
	"plain/pkg/vm"
	"testing"
)

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

	// For a variable declaration statement, the value remains stored in globals[0].
	// Let's test compiling a second statement that loads it: "total"
	// Or we can verify the VM globals directly if exposed, or append a load test.
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

	// Verify that globals[0] (or whichever slot 'result' took) got assigned 42!
	// Let's check how your symbol table assigns slots.
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
	sumVal, exists := machine.Globals()[0] // assuming sum is slot 0
	if !exists {
		t.Fatalf("expected variable 'sum' in globals")
	}
	if sumVal.Num != 15 {
		t.Errorf("expected sum 15, got %f", sumVal.Num)
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
	resultVal, exists := machine.Globals()[1] // Assuming result is slot 1 (slot 0 being the function object itself)
	if !exists {
		t.Fatalf("expected variable 'result' in globals")
	}
	if resultVal.Num != 30 {
		t.Errorf("expected result 30, got %f", resultVal.Num)
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

	resultVal, exists := machine.Globals()[1] // Assuming result slot
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

	resultVal, exists := machine.Globals()[2] // Assuming result slot
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
	// Slot layout: 0: MakeAdder, 1: addFive, 2: result
	resultVal, exists := machine.Globals()[2]
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
	Var first As String = items[0]
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
