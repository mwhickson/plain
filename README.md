# PLAIN Programming Language

PLAIN (Programming Language Allowing Implementation Neatly) is a statically-typed, BASIC-inspired programming language compiled to custom bytecode and executed on a stack-based virtual machine.

PLAIN, and its associated tooling, is intended to fill software development niches historically attributed to dBase, Delphi, FoxPro, HyperCard, MS Access, PowerBuilder, Xojo, and Visual Basic.

> NOTE: PLAIN is a working-title backronym.

## Project Goals

- **Approachable Syntax:** Clear, English-like keywords and structured blocks (`If...End If`, `For...End For`, `Function...End Function`) inspired by classic rapid application development environments.
- **Explicit Typing:** Enforced variable declarations (`Var x As Numeric`) to catch type mismatches at compile time.
- **Clean Architecture:** A decoupled pipeline spanning lexing, Pratt parsing, AST construction, bytecode compilation, and virtual machine execution.
- **Interactive Tooling:** A built-in CLI runner and a classic line-oriented REPL for rapid prototyping and debugging.

## Implemented Features

- **Primitive Types:** `Numeric`, `String`, `Boolean`, and `None`.
- **Control Flow:** Conditional branching (`If...Then...Else...End If`) and range-based loops (`For...To...Step...End For`).
- **Functions & Closures:** First-class functions, recursion, mutual recursion, and lexical closures capturing enclosing scope variables.
- **Data Structures:** Heterogeneous lists (`[1, 2, 3]`) with zero-based indexing, and user-defined `Struct` types with dot-notation property access and assignment.
- **Receiver Functions:** Methods bound to struct types (e.g., `Function (v As Vector2D) Scale(...)`).
- **Bytecode Disassembler:** Built-in instruction disassembly for inspecting compiled chunks.

## Language Examples

### Variables and Control Flow
```plain
Var total As Numeric = 0

For i As Numeric = 1 To 10 Step 1
    If i > 5 Then
        total = total + i
    End If
End For
```

### Functions and Closures
```plain
Function MakeAdder(x As Numeric) As Function
    Function Add(y As Numeric) As Numeric
        Return x + y
    End Function
    Return Add
End Function

Var addFive As Function = MakeAdder(5)
Var result As Numeric = addFive(10)
```

### Structs and Receiver Functions
```plain
Struct Vector2D
    X As Numeric
    Y As Numeric
End Struct

Function (v As Vector2D) Scale(factor As Numeric)
    v.X = v.X * factor
    v.Y = v.Y * factor
End Function
```

## Usage

### Running a Source File
Execute a `.plain` file directly from the command line:
```bash
go run cmd/plain/main.go path/to/script.plain
```

### Interactive REPL Mode
Launch the CLI without arguments to start the interactive REPL:
```bash
go run cmd/plain/main.go
```

#### REPL Commands
- `help` - Display available REPL commands
- `new` - Clear the current session buffer
- `list` - Display numbered lines currently stored in memory
- `run` - Compile and execute the code in memory
- `edit <line>` - Modify a specific line number (submit a blank line to delete it)
- `save <file>` - Save the current session lines to disk
- `load <file>` - Load a file into the session buffer and execute it
- `exit` / `quit` - Exit the REPL

## Package Structure

```text
plain/
├── cmd/
│   └── plain/
│       └── main.go          # CLI entry point (file execution and interactive REPL)
├── pkg/
│   ├── ast/
│   │   └── ast.go           # Abstract Syntax Tree node definitions
│   ├── bytecode/
│   │   ├── chunk.go         # Bytecode chunk and constant pool containers
│   │   ├── disassembler.go  # Human-readable bytecode printer
│   │   └── opcode.go        # Opcode definitions and string representations
│   ├── compiler/
│   │   ├── compiler.go      # AST to bytecode compiler and type checker
│   │   └── symbol.go        # Lexical scoping and symbol resolution
│   ├── parser/
│   │   ├── lexer.go         # Tokenizer (source text to tokens)
│   │   └── parser.go        # Pratt parser (tokens to AST)
│   ├── token/
│   │   └── token.go         # Lexical token types and keyword lookups
│   ├── value/
│   │   └── value.go         # Runtime value representations
│   └── vm/
│       ├── frame.go         # Call frames and closure state
│       └── vm.go            # Stack-based bytecode execution loop
├── go.mod
└── README.md
```

## Disclaimer

This application was produced using AI (Gemini).
