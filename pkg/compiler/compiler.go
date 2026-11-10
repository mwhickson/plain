package compiler

import (
	"fmt"
	"plain/pkg/ast"
	"plain/pkg/bytecode"
	"plain/pkg/value"
)

type Compiler struct {
	chunk            *bytecode.Chunk
	symbolTable      map[string]byte
	variableTypes    map[string]string // Tracks declared type per variable name
	localSymbolTable map[string]byte
	varCount         byte
}

func New() *Compiler {
	return &Compiler{
		chunk:            bytecode.NewChunk(),
		symbolTable:      make(map[string]byte),
		localSymbolTable: nil,
		varCount:         0,
	}
}

func (c *Compiler) Compile(node ast.Node) error {
	switch n := node.(type) {
	case *ast.AssignmentStatement:
		// 1. Compile RHS expression
		if err := c.Compile(n.Value); err != nil {
			return err
		}
		// 2. Lookup existing variable slot
		varIdx, exists := c.symbolTable[n.Name.Value]
		if !exists {
			return fmt.Errorf("undefined variable for assignment: %s", n.Name.Value)
		}
		// 3. Emit store instruction and slot index
		c.chunk.WriteOpcode(bytecode.OpStoreVar)
		c.chunk.WriteByte(varIdx)

	case *ast.BooleanLiteral:
		val := value.NewBoolean(n.Value)
		constIdx := c.chunk.AddConstant(val)
		c.chunk.WriteOpcode(bytecode.OpPushConst)
		c.chunk.WriteByte(constIdx)

	case *ast.NoneLiteral:
		val := value.NewNone()
		constIdx := c.chunk.AddConstant(val)
		c.chunk.WriteOpcode(bytecode.OpPushConst)
		c.chunk.WriteByte(constIdx)

	case *ast.Program:
		for _, s := range n.Statements {
			if err := c.Compile(s); err != nil {
				return err
			}
		}

	case *ast.CallExpression: // or whatever your AST node for calls is named
		// 1. Compile the function expression (evaluates to the FunctionObject)
		if err := c.Compile(n.Function); err != nil {
			return err
		}
		// 2. Compile arguments
		for _, arg := range n.Arguments {
			if err := c.Compile(arg); err != nil {
				return err
			}
		}
		// 3. Emit OpCall with argument count
		c.chunk.WriteOpcode(bytecode.OpCall)
		c.chunk.WriteByte(byte(len(n.Arguments)))

	case *ast.ExpressionStatement:
		if err := c.Compile(n.Expression); err != nil {
			return err
		}

	case *ast.FunctionStatement:
		// 1. Emit a jump around the function body
		skipJumpPos := c.emit(bytecode.OpJump, 9999)

		// 2. Record function start address
		funcStartIndex := len(c.chunk.Code)

		// Save outer local symbol table state
		oldLocalSymbolTable := c.localSymbolTable

		// Initialize local symbol table for this function scope
		c.localSymbolTable = make(map[string]byte)
		for idx, param := range n.Parameters {
			c.localSymbolTable[param.Name.Value] = byte(idx)
		}

		// 3. Compile function body statements
		for _, stmt := range n.Body {
			if err := c.Compile(stmt); err != nil {
				return err
			}
		}

		// Ensure every function ends with a return (defaulting to None)
		noneConstIdx := c.chunk.AddConstant(value.NewNone())
		c.chunk.WriteOpcode(bytecode.OpPushConst)
		c.chunk.WriteByte(noneConstIdx)
		c.chunk.WriteOpcode(bytecode.OpReturn)

		// Restore outer local symbol table
		c.localSymbolTable = oldLocalSymbolTable

		// Patch skip jump
		endFuncIndex := len(c.chunk.Code)
		c.changeOperand(skipJumpPos, endFuncIndex)

		// 4. Create function object and store it as a global constant/variable
		fnObj := &value.FunctionObject{
			Name:                  n.Name.Value,
			InstructionStartIndex: funcStartIndex,
			ParamCount:            len(n.Parameters),
		}
		constIdx := c.chunk.AddConstant(value.Value{Type: value.ValFunction, Obj: fnObj})

		varIdx, exists := c.symbolTable[n.Name.Value]
		if !exists {
			varIdx = c.varCount
			c.symbolTable[n.Name.Value] = varIdx
			c.varCount++
		}
		c.chunk.WriteOpcode(bytecode.OpPushConst)
		c.chunk.WriteByte(constIdx)
		c.chunk.WriteOpcode(bytecode.OpStoreVar)
		c.chunk.WriteByte(varIdx)

	case *ast.ReturnStatement:
		if n.ReturnValue != nil {
			if err := c.Compile(n.ReturnValue); err != nil {
				return err
			}
		} else {
			// Return None by default
			noneConstIdx := c.chunk.AddConstant(value.NewNone())
			c.chunk.WriteOpcode(bytecode.OpPushConst)
			c.chunk.WriteByte(noneConstIdx)
		}
		c.chunk.WriteOpcode(bytecode.OpReturn)

	case *ast.VarStatement:
		if n.TypeName == "" {
			return fmt.Errorf("type required for variable '%s' (Option Explicit)", n.Name.Value)
		}

		// Infer type of RHS expression
		exprType, err := c.inferExpressionType(n.Value)
		if err != nil {
			return err
		}

		// If we can statically determine the RHS type, enforce a match
		if exprType != "" && exprType != n.TypeName {
			return fmt.Errorf("type mismatch for '%s': declared %s, got %s", n.Name.Value, n.TypeName, exprType)
		}

		// 1. Compile the RHS expression
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
		if c.localSymbolTable != nil {
			if localIdx, exists := c.localSymbolTable[n.Value]; exists {
				c.chunk.WriteOpcode(bytecode.OpLoadLocal)
				c.chunk.WriteByte(localIdx)
				return nil
			}
		}
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

	case *ast.StringLiteral:
		val := value.NewString(n.Value)
		constIdx := c.chunk.AddConstant(val)
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

func (c *Compiler) inferExpressionType(expr ast.Expression) (string, error) {
	switch expr.(type) {
	case *ast.IntegerLiteral:
		return "Numeric", nil
	case *ast.StringLiteral:
		return "String", nil
	case *ast.BooleanLiteral:
		return "Boolean", nil
	case *ast.NoneLiteral:
		return "None", nil
	default:
		// For binary expressions, calls, or identifiers, we can expand or default to untyped for now
		return "", nil
	}
}
