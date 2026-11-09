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
