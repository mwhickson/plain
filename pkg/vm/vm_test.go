package vm_test

import (
	"plain/pkg/bytecode"
	"plain/pkg/value"
	"plain/pkg/vm"
	"testing"
)

func TestVMAddOpcode(t *testing.T) {
	chunk := bytecode.NewChunk()

	// 1. Load Constant 10
	idx1 := chunk.AddConstant(value.NewNumeric(10))
	chunk.WriteOpcode(bytecode.OpPushConst)
	chunk.WriteByte(idx1)

	// 2. Load Constant 20
	idx2 := chunk.AddConstant(value.NewNumeric(20))
	chunk.WriteOpcode(bytecode.OpPushConst)
	chunk.WriteByte(idx2)

	// 3. Add & Halt
	chunk.WriteOpcode(bytecode.OpAdd)
	chunk.WriteOpcode(bytecode.OpHalt)

	// Diagnostic Disassembly check (prints to stdout during go test -v)
	bytecode.Disassemble(chunk, "Test Addition")

	// Execute VM
	machine := vm.NewVM(chunk)
	if err := machine.Run(); err != nil {
		t.Fatalf("VM execution error: %v", err)
	}

	// Verify top of stack
	result, err := machine.Pop()
	if err != nil {
		t.Fatalf("Failed to pop result: %v", err)
	}

	if result.Type != value.ValNumeric || result.Num != 30 {
		t.Errorf("Expected 30, got %v", result)
	}
}