package compiler

import (
	"fmt"
	"plain/pkg/ast"
	"plain/pkg/bytecode"
	"plain/pkg/value"
)

type Compiler struct {
	chunk          *bytecode.Chunk
	symbolTable    *SymbolTable
	variableTypes  map[string]string // Tracks declared type per variable name
	globalVarCount byte
	localVarCount  byte
	inFunction     bool
}

func New() *Compiler {
	c := &Compiler{
		chunk:          bytecode.NewChunk(),
		symbolTable:    NewSymbolTable(),
		variableTypes:  make(map[string]string),
		globalVarCount: 0,
		localVarCount:  0,
		inFunction:     false,
	}

	// Pre-register standard modules in the global symbol table
	c.registerGlobalModule("Console")

	return c
}

func (c *Compiler) Bytecode() *bytecode.Chunk {
	// Append a clean Halt at the end of compilation
	c.chunk.WriteOpcode(bytecode.OpHalt)
	return c.chunk
}

func (c *Compiler) Compile(node ast.Node) error {
	switch n := node.(type) {
	case *ast.AssignmentStatement:
		switch left := n.Left.(type) {
		case *ast.Identifier:
			// 1. Compile RHS expression
			if err := c.Compile(n.Value); err != nil {
				return err
			}
			// 2. Lookup existing variable slot using Resolve
			symbol, exists := c.symbolTable.Resolve(left.Value)
			if !exists {
				return fmt.Errorf("undefined variable for assignment: %s", left.Value)
			}
			// 3. Emit store instruction and slot index based on scope
			if symbol.Scope == GlobalScope || !c.inFunction {
				c.chunk.WriteOpcode(bytecode.OpStoreVar)
				c.chunk.WriteByte(symbol.Index)
			} else {
				c.chunk.WriteOpcode(bytecode.OpStoreLocal)
				c.chunk.WriteByte(symbol.Index)
			}

		case *ast.PropertyAccessExpression:
			// 1. Compile object expression (pushes struct value)
			if err := c.Compile(left.Left); err != nil {
				return err
			}
			// 2. Compile RHS expression (pushes new value)
			if err := c.Compile(n.Value); err != nil {
				return err
			}
			// 3. Add property name as constant
			propNameVal := value.NewString(left.Property.Value)
			constIdx := c.chunk.AddConstant(propNameVal)
			c.chunk.WriteOpcode(bytecode.OpSetProperty)
			c.chunk.WriteByte(constIdx)

		default:
			return fmt.Errorf("invalid left-hand side for assignment: %T", n.Left)
		}

	case *ast.BooleanLiteral:
		val := value.NewBoolean(n.Value)
		constIdx := c.chunk.AddConstant(val)
		c.chunk.WriteOpcode(bytecode.OpPushConst)
		c.chunk.WriteByte(constIdx)

	case *ast.CallExpression:
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

	case *ast.ForStatement:
		// 1. Enter loop block scope (prevents loop iterator leakage)
		c.symbolTable = NewEnclosedSymbolTable(c.symbolTable)

		// 2. Evaluate Start expression
		if err := c.Compile(n.Start); err != nil {
			c.symbolTable = c.symbolTable.Parent // Ensure scope is restored on error
			return err
		}

		// Allocate or define loop variable in the current enclosed scope
		symbol, exists := c.symbolTable.Resolve(n.VarName.Value)
		if !exists {
			if !c.inFunction {
				symbol = Symbol{Name: n.VarName.Value, Scope: GlobalScope, Index: c.globalVarCount}
				c.symbolTable.Store[n.VarName.Value] = symbol
				c.globalVarCount++
			} else {
				symbol = c.symbolTable.Define(n.VarName.Value, c.localVarCount)
				c.localVarCount++
			}
		}

		// Store start value into loop variable
		if !c.inFunction {
			c.chunk.WriteOpcode(bytecode.OpStoreVar)
			c.chunk.WriteByte(symbol.Index)
		} else {
			c.chunk.WriteOpcode(bytecode.OpStoreLocal)
			c.chunk.WriteByte(symbol.Index)
		}

		// 3. Evaluate End expression and store in a temporary variable slot
		if err := c.Compile(n.End); err != nil {
			c.symbolTable = c.symbolTable.Parent
			return err
		}

		var endSymbol Symbol
		if !c.inFunction {
			endSymbol = Symbol{Name: "__end_temp", Scope: GlobalScope, Index: c.globalVarCount}
			c.symbolTable.Store["__end_temp"] = endSymbol
			c.globalVarCount++
			c.chunk.WriteOpcode(bytecode.OpStoreVar)
		} else {
			endSymbol = c.symbolTable.Define("__end_temp", c.localVarCount)
			c.localVarCount++
			c.chunk.WriteOpcode(bytecode.OpStoreLocal)
		}
		c.chunk.WriteByte(endSymbol.Index)

		// 4. Evaluate Step expression (or default to 1) and store in a temp slot
		if n.Step != nil {
			if err := c.Compile(n.Step); err != nil {
				c.symbolTable = c.symbolTable.Parent
				return err
			}
		} else {
			litVal := value.NewNumeric(1)
			constIdx := c.chunk.AddConstant(litVal)
			c.chunk.WriteOpcode(bytecode.OpPushConst)
			c.chunk.WriteByte(constIdx)
		}

		var stepSymbol Symbol
		if !c.inFunction {
			stepSymbol = Symbol{Name: "__step_temp", Scope: GlobalScope, Index: c.globalVarCount}
			c.symbolTable.Store["__step_temp"] = stepSymbol
			c.globalVarCount++
			c.chunk.WriteOpcode(bytecode.OpStoreVar)
		} else {
			stepSymbol = c.symbolTable.Define("__step_temp", c.localVarCount)
			c.localVarCount++
			c.chunk.WriteOpcode(bytecode.OpStoreLocal)
		}
		c.chunk.WriteByte(stepSymbol.Index)

		// --- LOOP START POSITION ---
		loopStartIndex := len(c.chunk.Code)

		// Condition Check: load loop var, load end temp, compare <=
		if !c.inFunction {
			c.chunk.WriteOpcode(bytecode.OpLoadVar)
		} else {
			c.chunk.WriteOpcode(bytecode.OpLoadLocal)
		}
		c.chunk.WriteByte(symbol.Index)

		if !c.inFunction {
			c.chunk.WriteOpcode(bytecode.OpLoadVar)
		} else {
			c.chunk.WriteOpcode(bytecode.OpLoadLocal)
		}
		c.chunk.WriteByte(endSymbol.Index)

		c.chunk.WriteOpcode(bytecode.OpLte)

		// Emit JumpIfFalse to exit the loop
		exitJumpPos := c.emit(bytecode.OpJumpIfFalse, 9999)

		// 5. Compile loop body statements
		for _, stmt := range n.Body {
			if err := c.Compile(stmt); err != nil {
				c.symbolTable = c.symbolTable.Parent
				return err
			}
		}

		// 6. Increment loop variable: i = i + step
		if !c.inFunction {
			c.chunk.WriteOpcode(bytecode.OpLoadVar)
		} else {
			c.chunk.WriteOpcode(bytecode.OpLoadLocal)
		}
		c.chunk.WriteByte(symbol.Index)

		if !c.inFunction {
			c.chunk.WriteOpcode(bytecode.OpLoadVar)
		} else {
			c.chunk.WriteOpcode(bytecode.OpLoadLocal)
		}
		c.chunk.WriteByte(stepSymbol.Index)

		c.chunk.WriteOpcode(bytecode.OpAdd)

		if !c.inFunction {
			c.chunk.WriteOpcode(bytecode.OpStoreVar)
		} else {
			c.chunk.WriteOpcode(bytecode.OpStoreLocal)
		}
		c.chunk.WriteByte(symbol.Index)

		// 7. Jump back to loop start
		c.emit(bytecode.OpJump, loopStartIndex)

		// Patch exit jump target
		endLoopIndex := len(c.chunk.Code)
		c.changeOperand(exitJumpPos, endLoopIndex)

		// 8. Exit loop block scope
		c.symbolTable = c.symbolTable.Parent

	case *ast.FunctionStatement:
		// 1. Pre-register function name in the current symbol table so it can reference itself recursively
		funcSymbol, exists := c.symbolTable.Resolve(n.Name.Value)
		if !exists {
			// Define function slot in the outer scope
			if !c.inFunction {
				funcSymbol = Symbol{Name: n.Name.Value, Scope: GlobalScope, Index: c.globalVarCount}
				c.symbolTable.Store[n.Name.Value] = funcSymbol
				c.globalVarCount++
			} else {
				funcSymbol = c.symbolTable.Define(n.Name.Value, c.localVarCount)
				c.localVarCount++ // <--- Ensure local variable slot count increments here!
			}
		}

		// 2. Emit a jump around the function body
		skipJumpPos := c.emit(bytecode.OpJump, 9999)

		// 3. Record function start address
		funcStartIndex := len(c.chunk.Code)

		// 4. Enter enclosed symbol table for the function scope
		oldSymbolTable := c.symbolTable
		c.symbolTable = NewEnclosedSymbolTable(oldSymbolTable)

		oldInFunction := c.inFunction
		c.inFunction = true

		oldLocalVarCount := c.localVarCount
		c.localVarCount = 0 // Function frames start local indices at 0

		// Define receiver in the function's symbol table if present
		if n.Receiver != nil {
			c.symbolTable.Define(n.Receiver.Name.Value, c.localVarCount)
			c.localVarCount++
		}

		// Define parameters in the function's symbol table with explicit indices
		for _, param := range n.Parameters {
			c.symbolTable.Define(param.Name.Value, c.localVarCount)
			c.localVarCount++
		}

		// 5. Compile function body statements
		for _, stmt := range n.Body {
			if err := c.Compile(stmt); err != nil {
				c.symbolTable = oldSymbolTable
				c.inFunction = oldInFunction
				c.localVarCount = oldLocalVarCount
				return err
			}
		}

		// Ensure every function ends with a return (defaulting to None)
		noneConstIdx := c.chunk.AddConstant(value.NewNone())
		c.chunk.WriteOpcode(bytecode.OpPushConst)
		c.chunk.WriteByte(noneConstIdx)
		c.chunk.WriteOpcode(bytecode.OpReturn)

		// Capture the inner symbol table and local var count before restoring the outer context!
		funcSymbolTable := c.symbolTable
		functionLocalVarCount := c.localVarCount

		// Restore outer context
		c.symbolTable = oldSymbolTable
		c.inFunction = oldInFunction
		c.localVarCount = oldLocalVarCount

		// Patch skip jump & store function object...
		endFuncIndex := len(c.chunk.Code)
		c.changeOperand(skipJumpPos, endFuncIndex)

		// 6. Check if the function captured any free variables using the inner symbol table
		freeSymbols := funcSymbolTable.FreeSymbols
		numFree := len(freeSymbols)

		// Create function object
		fnObj := &value.FunctionObject{
			Name:                  n.Name.Value,
			InstructionStartIndex: funcStartIndex,
			ParamCount:            len(n.Parameters),
			LocalVarCount:         int(functionLocalVarCount),
		}
		constIdx := c.chunk.AddConstant(value.Value{Type: value.ValFunction, Obj: fnObj})

		if numFree > 0 {
			// For each free symbol, emit instructions to load it from the enclosing scope
			for _, s := range freeSymbols {
				if s.Scope == LocalScope {
					c.chunk.WriteOpcode(bytecode.OpLoadLocal)
					c.chunk.WriteByte(s.Index)
				} else if s.Scope == FreeScope {
					c.chunk.WriteOpcode(bytecode.OpLoadFree)
					c.chunk.WriteByte(s.Index)
				}
			}

			// Emit OpClosure taking the function constant index and free count
			c.chunk.WriteOpcode(bytecode.OpClosure)
			c.chunk.WriteByte(constIdx)
			c.chunk.WriteByte(byte(numFree))
		} else {
			// Standard function (no free variables)
			c.chunk.WriteOpcode(bytecode.OpPushConst)
			c.chunk.WriteByte(constIdx)
		}

		// Store the resulting function/closure into its variable slot
		if funcSymbol.Scope == GlobalScope {
			c.chunk.WriteOpcode(bytecode.OpStoreVar)
			c.chunk.WriteByte(funcSymbol.Index)
		} else {
			c.chunk.WriteOpcode(bytecode.OpStoreLocal)
			c.chunk.WriteByte(funcSymbol.Index)
		}

	case *ast.Identifier:
		symbol, exists := c.symbolTable.Resolve(n.Value)
		if !exists {
			return fmt.Errorf("undefined variable: %s", n.Value)
		}

		if symbol.Scope == GlobalScope || !c.inFunction {
			c.chunk.WriteOpcode(bytecode.OpLoadVar)
			c.chunk.WriteByte(symbol.Index)
		} else if symbol.Scope == LocalScope {
			c.chunk.WriteOpcode(bytecode.OpLoadLocal)
			c.chunk.WriteByte(symbol.Index)
		} else if symbol.Scope == FreeScope {
			c.chunk.WriteOpcode(bytecode.OpLoadFree)
			c.chunk.WriteByte(symbol.Index)
		}
		return nil

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

	case *ast.IndexExpression:
		if err := c.Compile(n.Left); err != nil {
			return err
		}
		if err := c.Compile(n.Index); err != nil {
			return err
		}
		c.chunk.WriteOpcode(bytecode.OpIndex)

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
		case "/":
			c.chunk.WriteOpcode(bytecode.OpDiv)
		case "==":
			c.chunk.WriteOpcode(bytecode.OpEq)
		case ">":
			c.chunk.WriteOpcode(bytecode.OpGt)
		case ">=":
			c.chunk.WriteOpcode(bytecode.OpGte)
		case "<":
			c.chunk.WriteOpcode(bytecode.OpLt)
		case "<=":
			c.chunk.WriteOpcode(bytecode.OpLte)
		case "*":
			c.chunk.WriteOpcode(bytecode.OpMul)
		case "!=":
			c.chunk.WriteOpcode(bytecode.OpNeq)
		case "-":
			c.chunk.WriteOpcode(bytecode.OpSub)
		default:
			return fmt.Errorf("unsupported infix operator %s", n.Operator)
		}

	case *ast.IntegerLiteral:
		literalVal := value.NewNumeric(n.Value)
		constIdx := c.chunk.AddConstant(literalVal)
		c.chunk.WriteOpcode(bytecode.OpPushConst)
		c.chunk.WriteByte(constIdx)

	case *ast.ListLiteral:
		for _, el := range n.Elements {
			if err := c.Compile(el); err != nil {
				return err
			}
		}
		c.chunk.WriteOpcode(bytecode.OpBuildList)
		c.chunk.WriteByte(byte(len(n.Elements)))

	case *ast.NoneLiteral:
		val := value.NewNone()
		constIdx := c.chunk.AddConstant(val)
		c.chunk.WriteOpcode(bytecode.OpPushConst)
		c.chunk.WriteByte(constIdx)

	case *ast.Program:
		// 1. Pre-scan pass: register all top-level functions first so they can mutually recurse
		for _, s := range n.Statements {
			if fnStmt, ok := s.(*ast.FunctionStatement); ok {
				_, exists := c.symbolTable.Resolve(fnStmt.Name.Value)
				if !exists {
					if !c.inFunction {
						symbol := Symbol{
							Name:  fnStmt.Name.Value,
							Scope: GlobalScope,
							Index: c.globalVarCount,
						}
						c.symbolTable.Store[fnStmt.Name.Value] = symbol
						c.globalVarCount++
					}
				}
			}
		}

		// 2. Compilation pass: compile all statements normally now that all functions are known
		for _, s := range n.Statements {
			if err := c.Compile(s); err != nil {
				return err
			}
		}

	case *ast.PropertyAccessExpression:
		if err := c.Compile(n.Left); err != nil {
			return err
		}
		// Add property name as a constant in the constant pool
		propNameVal := value.NewString(n.Property.Value)
		constIdx := c.chunk.AddConstant(propNameVal)
		c.chunk.WriteOpcode(bytecode.OpGetProperty)
		c.chunk.WriteByte(constIdx)

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

	case *ast.StringLiteral:
		val := value.NewString(n.Value)
		constIdx := c.chunk.AddConstant(val)
		c.chunk.WriteOpcode(bytecode.OpPushConst)
		c.chunk.WriteByte(constIdx)

	case *ast.StructStatement:
		// Struct definitions act as types/schemas; for now we can treat them as declarations
		// that don't emit runtime executable instructions (or store struct schema metadata if needed).
		return nil

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

		// 2. Define variable in the current symbol table scope with explicit index
		var symbol Symbol
		if !c.inFunction {
			symbol, exists := c.symbolTable.Resolve(n.Name.Value)
			if !exists {
				symbol = Symbol{
					Name:  n.Name.Value,
					Scope: GlobalScope,
					Index: c.globalVarCount,
				}
				c.symbolTable.Store[n.Name.Value] = symbol
				c.globalVarCount++
			}
			c.chunk.WriteOpcode(bytecode.OpStoreVar)
			c.chunk.WriteByte(symbol.Index)
		} else {
			symbol = c.symbolTable.Define(n.Name.Value, c.localVarCount)
			c.localVarCount++
			c.chunk.WriteOpcode(bytecode.OpStoreLocal)
			c.chunk.WriteByte(symbol.Index)
		}

	default:
		return fmt.Errorf("unsupported AST node type: %T", n)
	}

	return nil
}

func (c *Compiler) changeOperand(opPos int, operand int) {
	// The operand starts immediately after the opcode
	c.chunk.Code[opPos+1] = byte(operand & 0xff)
	c.chunk.Code[opPos+2] = byte((operand >> 8) & 0xff)
}

func (c *Compiler) emit(op bytecode.Opcode, operand int) int {
	pos := len(c.chunk.Code)
	c.chunk.WriteOpcode(op)
	// Write 16-bit operand (low byte, high byte)
	c.chunk.WriteByte(byte(operand & 0xff))
	c.chunk.WriteByte(byte((operand >> 8) & 0xff))
	return pos
}

func (c *Compiler) inferExpressionType(expr ast.Expression) (string, error) {
	switch expr.(type) {
	case *ast.BooleanLiteral:
		return "Boolean", nil
	case *ast.IntegerLiteral:
		return "Numeric", nil
	case *ast.NoneLiteral:
		return "None", nil
	case *ast.StringLiteral:
		return "String", nil
	default:
		// For binary expressions, calls, or identifiers, we can expand or default to untyped for now
		return "", nil
	}
}

func (c *Compiler) registerGlobalModule(name string) {
	symbol := Symbol{
		Name:  name,
		Scope: GlobalScope,
		Index: c.globalVarCount,
	}
	c.symbolTable.Store[name] = symbol
	c.variableTypes[name] = "Struct" // Modules are struct namespaces
	c.globalVarCount++
}
