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