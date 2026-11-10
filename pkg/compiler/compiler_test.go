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
