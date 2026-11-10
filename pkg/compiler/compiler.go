package compiler

import (
	"fmt"
	"plain/pkg/ast"
	"plain/pkg/bytecode"
	"plain/pkg/value"
)

type Compiler struct {
	chunk       *bytecode.Chunk
	symbolTable map[string]byte
	varCount    byte
}

func New() *Compiler {
	return &Compiler{
		chunk:       bytecode.NewChunk(),
		symbolTable: make(map[string]byte),
		varCount:    0,
	}
}

func (c *Compiler) Compile(node ast.Node) error {
	switch n := node.(type) {
	case *ast.Program:
		for _, s := range n.Statements {
			if err := c.Compile(s); err != nil {
				return err
			}
		}

	case *ast.ExpressionStatement:
		if err := c.Compile(n.Expression); err != nil {
			return err
		}

	case *ast.VarStatement:
		// 1. Compile the RHS expression (evaluates and pushes result to stack)
		if err := c.Compile(n.Value); err != nil {
			return err
		}
		// 2. Assign or lookup variable slot
		varIdx, exists := c.symbolTable[n.Name.Value]
		if !exists {
			varIdx = c.varCount
			c.symbolTable[n.Name.Value] = varIdx
			c.varCount++
		}
		// 3. Emit store instruction and variable index
		c.chunk.WriteOpcode(bytecode.OpStoreVar)
		c.chunk.WriteByte(varIdx)

	case *ast.Identifier:
		varIdx, exists := c.symbolTable[n.Value]
		if !exists {
			return fmt.Errorf("undefined variable: %s", n.Value)
		}
		c.chunk.WriteOpcode(bytecode.OpLoadVar)
		c.chunk.WriteByte(varIdx)

	case *ast.IfStatement:
		// 1. Compile the condition expression
		if err := c.Compile(n.Condition); err != nil {
			return err
		}

		// 2. Emit OpJumpIfFalse with a placeholder operand (9999)
		jumpIfFalsePos := c.emit(bytecode.OpJumpIfFalse, 9999)

		// 3. Compile the consequence block
		for _, stmt := range n.Consequence {
			if err := c.Compile(stmt); err != nil {
				return err
			}
		}

		if len(n.Alternative) > 0 {
			// Emit OpJump to skip the alternative block after consequence runs
			jumpPos := c.emit(bytecode.OpJump, 9999)

			// Current position is where the alternative block starts
			alternativePos := len(c.chunk.Code)
			c.changeOperand(jumpIfFalsePos, alternativePos)

			// Compile alternative block
			for _, stmt := range n.Alternative {
				if err := c.Compile(stmt); err != nil {
					return err
				}
			}

			// Patch the jumpPos to point to the end
			endPos := len(c.chunk.Code)
			c.changeOperand(jumpPos, endPos)
		} else {
			// No alternative block; patch jumpIfFalsePos to jump straight to the end
			endPos := len(c.chunk.Code)
			c.changeOperand(jumpIfFalsePos, endPos)
		}

	case *ast.IntegerLiteral:
		literalVal := value.NewNumeric(n.Value)
		constIdx := c.chunk.AddConstant(literalVal)
		c.chunk.WriteOpcode(bytecode.OpPushConst)
		c.chunk.WriteByte(constIdx)

	case *ast.InfixExpression:
		// Compile left-hand side
		if err := c.Compile(n.Left); err != nil {
			return err
		}
		// Compile right-hand side
		if err := c.Compile(n.Right); err != nil {
			return err
		}

		// Emit operator instruction
		switch n.Operator {
		case "+":
			c.chunk.WriteOpcode(bytecode.OpAdd)
		case "-":
			c.chunk.WriteOpcode(bytecode.OpSub)
		case "*":
			c.chunk.WriteOpcode(bytecode.OpMul)
		case "/":
			c.chunk.WriteOpcode(bytecode.OpDiv)
		case "==":
			c.chunk.WriteOpcode(bytecode.OpEq)
		default:
			return fmt.Errorf("unsupported infix operator %s", n.Operator)
		}

	default:
		return fmt.Errorf("unsupported AST node type: %T", n)
	}

	return nil
}

func (c *Compiler) Bytecode() *bytecode.Chunk {
	// Append a clean Halt at the end of compilation
	c.chunk.WriteOpcode(bytecode.OpHalt)
	return c.chunk
}

func (c *Compiler) emit(op bytecode.Opcode, operand int) int {
	pos := len(c.chunk.Code)
	c.chunk.WriteOpcode(op)
	// Write 16-bit operand (low byte, high byte)
	c.chunk.WriteByte(byte(operand & 0xff))
	c.chunk.WriteByte(byte((operand >> 8) & 0xff))
	return pos
}

func (c *Compiler) changeOperand(opPos int, operand int) {
	// The operand starts immediately after the opcode
	c.chunk.Code[opPos+1] = byte(operand & 0xff)
	c.chunk.Code[opPos+2] = byte((operand >> 8) & 0xff)
}
