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

	case *ast.ForStatement:
		// 1. Evaluate Start expression
		if err := c.Compile(n.Start); err != nil {
			return err
		}
		// Allocate or lookup loop variable slot
		varIdx, exists := c.symbolTable[n.VarName.Value]
		if !exists {
			varIdx = c.varCount
			c.symbolTable[n.VarName.Value] = varIdx
			c.varCount++
		}
		// Store start value into loop variable
		c.chunk.WriteOpcode(bytecode.OpStoreVar)
		c.chunk.WriteByte(varIdx)

		// 2. Evaluate End expression and store in a hidden temporary variable slot
		if err := c.Compile(n.End); err != nil {
			return err
		}
		endVarIdx := c.varCount
		c.varCount++
		c.chunk.WriteOpcode(bytecode.OpStoreVar)
		c.chunk.WriteByte(endVarIdx)

		// 3. Evaluate Step expression (or default to 1) and store in a temp slot
		stepVarIdx := c.varCount
		c.varCount++
		if n.Step != nil {
			if err := c.Compile(n.Step); err != nil {
				return err
			}
		} else {
			// Push constant 1
			litVal := value.NewNumeric(1)
			constIdx := c.chunk.AddConstant(litVal)
			c.chunk.WriteOpcode(bytecode.OpPushConst)
			c.chunk.WriteByte(constIdx)
		}
		c.chunk.WriteOpcode(bytecode.OpStoreVar)
		c.chunk.WriteByte(stepVarIdx)

		// --- LOOP START POSITION ---
		loopStartIndex := len(c.chunk.Code)

		// Condition Check: (i <= end) -> load i, load end, compare <=
		c.chunk.WriteOpcode(bytecode.OpLoadVar)
		c.chunk.WriteByte(varIdx)
		c.chunk.WriteOpcode(bytecode.OpLoadVar)
		c.chunk.WriteByte(endVarIdx)

		// Let's ensure we have a Less-Than-Or-Equal opcode (`OpLte`), or use comparison.
		// (Let's check if OpLte is in your bytecode opcode definitions!)
		c.chunk.WriteOpcode(bytecode.OpLte)

		// Emit JumpIfFalse to exit the loop (placeholder target)
		exitJumpPos := c.emit(bytecode.OpJumpIfFalse, 9999)

		// 4. Compile loop body statements
		for _, stmt := range n.Body {
			if err := c.Compile(stmt); err != nil {
				return err
			}
		}

		// 5. Increment loop variable: i = i + step
		c.chunk.WriteOpcode(bytecode.OpLoadVar)
		c.chunk.WriteByte(varIdx)
		c.chunk.WriteOpcode(bytecode.OpLoadVar)
		c.chunk.WriteByte(stepVarIdx)
		c.chunk.WriteOpcode(bytecode.OpAdd)
		c.chunk.WriteOpcode(bytecode.OpStoreVar)
		c.chunk.WriteByte(varIdx)

		// 6. Jump back to loop start
		c.emit(bytecode.OpJump, loopStartIndex)

		// Patch exit jump target
		endLoopIndex := len(c.chunk.Code)
		c.changeOperand(exitJumpPos, endLoopIndex)

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
