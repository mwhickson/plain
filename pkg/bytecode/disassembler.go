package bytecode

import (
	"fmt"
)

// Disassemble prints a human-readable stream of instructions in a Chunk.
func Disassemble(c *Chunk, name string) {
	fmt.Printf("=== %s ===\n", name)
	offset := 0
	for offset < len(c.Code) {
		offset = DisassembleInstruction(c, offset)
	}
}

// DisassembleInstruction prints a single instruction and returns the next instruction offset.
func DisassembleInstruction(c *Chunk, offset int) int {
	fmt.Printf("%04d  ", offset)
	op := Opcode(c.Code[offset])

	switch op {
	case OpPushConst:
		constIdx := c.Code[offset+1]
		val := c.Constants[constIdx]
		fmt.Printf("%-16s %4d ('%s')\n", op, constIdx, val)
		return offset + 2
	case OpStoreVar, OpLoadVar:
		varIdx := c.Code[offset+1]
		fmt.Printf("%-16s %4d\n", op, varIdx)
		return offset + 2
	case OpNop, OpPushNone, OpPushTrue, OpPushFalse, OpPop,
		OpAdd, OpSub, OpMul, OpDiv, OpMod, OpNegate,
		OpEq, OpNeq, OpLt, OpLte, OpGt, OpGte, OpReturn, OpHalt:
		fmt.Printf("%s\n", op)
		return offset + 1
	default:
		fmt.Printf("Unknown opcode %d\n", op)
		return offset + 1
	}
}