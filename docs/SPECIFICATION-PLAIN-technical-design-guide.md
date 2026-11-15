# PLAIN Technical Design Document: Compiler and Virtual Machine Architecture

**Target Audience:** Compiler Engineers, VM Implementers, and Tooling/Grammar Authors (TextMate, Tree-sitter)  
**Status:** Architecture Reference (v0.1)

---

## 1. Executive Summary

PLAIN (Programming Language Allowing Implementation Neatly) is a statically-typed, BASIC-inspired programming language engineered around a clean, single-pass Abstract Syntax Tree (AST) compiler targeting a custom stack-based Virtual Machine in Go.

The architecture prioritizes **transparency, low implementation overhead, and modular extensibility**. Rather than relying on complex foreign function interfaces (FFIs) or heavy runtime bytecode verifiers, PLAIN exposes native capabilities through pre-seeded standard-library `Struct` namespaces backed by Go callback functions (`ValBuiltin`).

```text
[ Source Code (.plain) ]
           │
           ▼
    ┌─────────────┐
    │    Lexer    │  ──► Token Stream (strips ' comments, normalizes keywords)
    └─────────────┘
           │
           ▼
    ┌─────────────┐
    │Pratt Parser │  ──► Abstract Syntax Tree (AST)
    └─────────────┘
           │
           ▼
    ┌─────────────┐
    │  Compiler   │  ──► Bytecode Chunk (.pbc) + Constant Pool + Symbol Resolution
    └─────────────┘
           │
           ▼
    ┌─────────────┐
    │ Stack VM    │  ──► Call Frames, Operand Stack, Global Module Slots, Error Hooks
    └─────────────┘
```

---

## 2. Frontend Pipeline: Lexing & Parsing

### 2.1 Lexical Analysis (`pkg/parser/lexer.go`, `pkg/token/token.go`)
The lexer performs single-pass UTF-8 scanning to produce tokens with line/column metadata:
- **Case-Insensitive Keywords:** Keywords (`Var`, `As`, `Function`, `If`, `Then`, `Else`, `End`, `For`, `To`, `Step`, `Return`, `Struct`, `True`, `False`, `None`) are matched case-insensitively during lookup, while PascalCase is preserved as the canonical style.
- **Line Comments:** Single-quote (`'`) characters initiate line comments; the lexer consumes all characters until the next newline (`\n`) or `EOF`.
- **Identifier Flexibility:** Standard library module names and method names are lexed as standard `IDENT` tokens rather than reserved keywords, allowing clean dot-notation property resolution without grammar pollution.

### 2.2 Pratt Parser (`pkg/parser/parser.go`)
PLAIN uses a Top-Down Operator Precedence (Pratt) parser to resolve expressions cleanly without left-recursion:
- **Line-Delimited Statements:** Newlines act as statement boundaries. Multi-line blocks (`If...End If`, `For...End For`, `Function...End Function`, `Struct...End Struct`) parse nested statement slices until their corresponding `End <Construct>` token sequence is encountered.
- **Postfix & Infix Dispatch:**
  - `(` triggers `CallExpression` parsing.
  - `.` triggers `PropertyAccessExpression` parsing (used for both struct fields and module methods).
  - `[` triggers `IndexExpression` parsing for lists.
- **Method Receivers:** Function declarations optionally parse a receiver clause between `Function` and the identifier: `Function (receiver As Type) MethodName(params...)`.

---

## 3. Compiler & Symbol Resolution (`pkg/compiler`)

The compiler walks the AST in a single pass, emitting 1-byte opcodes and little-endian 16-bit jump offsets into a `bytecode.Chunk`.

### 3.1 Two-Pass Top-Level Function Pre-Scan
When compiling an `*ast.Program` node, the compiler performs a lightweight pre-scan of all top-level `*ast.FunctionStatement` nodes before emitting code. This registers every global function symbol in the root `SymbolTable` up front, enabling **mutual recursion** and forward function calls without requiring header declarations.

### 3.2 Hierarchical Symbol Tables & Closures (`symbol.go`)
Scoping is managed via linked `SymbolTable` instances:
1. **Global Scope (`GlobalScope`):** Top-level variables, functions, and pre-registered standard modules (`0..255`). Accessed via `OpLoadVar` / `OpStoreVar`.
2. **Local Scope (`LocalScope`):** Function receivers, parameters, and local variables indexed relative to the active frame's `basePointer`. Accessed via `OpLoadLocal` / `OpStoreLocal`.
3. **Enclosed / Free Scope (`FreeScope`):** When an inner function references a variable defined in an outer function's `LocalScope`, `Resolve()` automatically promotes the symbol into the inner table's `FreeSymbols` list and emits `OpClosure` with upvalue capture instructions (`OpLoadFree`).

### 3.3 Static Type Verification
Variable declarations (`Var x As Type = expr`) enforce explicit type annotations (`Option Explicit` semantics). The compiler performs static type inference on literal and known expressions (`inferExpressionType`) to reject mismatches at compile time.

---

## 4. Standard Library Module Architecture (`pkg/modules`)

To avoid complex FFI bindings or dynamic linker overhead, PLAIN implements standard library modules as native `ValStruct` maps populated with `ValBuiltin` functions:

1. **Compiler Pre-Registration:** During `compiler.New()`, module names are registered into the global symbol table before user code is compiled, using the ordered `modules.StandardModules()` registry. The current modules occupy slots `0..10`; slots below `30` remain reserved for future standard-library modules.
2. **VM Slot Seeding:** During `vm.NewVM()`, the VM instantiates the same registry and binds each module to its corresponding index in `vm.globals`. User globals start at slot `30`.
3. **Uniform Invocation:** Because modules are standard `ValStruct` values holding `ValBuiltin` entries, calling `Console.Print("Hi")` compiles to standard `OpLoadVar 1` $\rightarrow$ `OpGetProperty "Print"` $\rightarrow$ `OpPushConst "Hi"` $\rightarrow$ `OpCall 1`. No special bytecode instructions are needed for standard library calls.

---

## 5. Virtual Machine & Error Interception (`pkg/vm`)

### 5.1 Stack & Call Frame Mechanics
The VM maintains an operand stack (`[]value.Value`) and a fixed-capacity call frame stack (`frames[1024]`):
- **Root Frame (`frames[0]`):** Represents top-level script execution.
- **Function Invocation (`OpCall`):**
  - If the callee is `ValBuiltin`, arguments are sliced directly from the top of the operand stack, passed to the Go function, and replaced with the return value without allocating a VM `Frame`.
  - If the callee is `ValFunction` or `ValClosure`, the VM saves `vm.ip` into `frames[frameIndex-1].ip`, pushes a new `Frame` with `basePointer` pointing to the first argument on the stack, allocates `None` slots for remaining local variables, and jumps `vm.ip` to `InstructionStartIndex`.
- **Function Return (`OpReturn`):** Pops the return value, truncates `vm.stack` back to `currentFrame.basePointer`, decrements `frameIndex`, restores `vm.ip` from the caller's frame, and pushes the return value.

### 5.2 Global Runtime Error Handler Hook
The `System` module exposes `System.SetErrorHandler(fn)`, storing a reference to a user-provided `ValFunction`, `ValClosure`, or `ValBuiltin` in `modules.GlobalErrorHandler`.

When the VM encounters a runtime error (e.g., division by zero, out-of-bounds index, type error):
1. `vm.handleError(err)` intercepts the error before `Run()` terminates.
2. If a handler is registered, the VM wraps `err.Error()` in a `ValString`, saves the current instruction pointer into the active frame (`vm.frames[vm.frameIndex-1].ip = vm.ip`), pushes a new call frame for the handler with the error string as its first argument, and continues execution.
3. Saving the active frame's `ip` prior to entering the handler is critical: when the handler executes `OpReturn`, the VM resumes at the instruction immediately following the fault rather than resetting `ip` to `0`.

---

## 6. Interactive REPL Design (`cmd/plain/main.go`)

The REPL is modeled after classic line-oriented BASIC environments while adapting to PLAIN's line-number-free syntax via two distinct input modes:

1. **Program Mode (Default):** Any standard input line is first passed through a silent parser check (`parser.NewParser`). If valid, the line is appended to the in-memory session buffer (`[]string`) without executing the program. Users inspect and manage the buffer via `list`, `edit <line>`, `new`, `save <file>`, `load <file>`, and `run`.
2. **Direct Mode (`!` prefix):** Lines prefixed with `!` (e.g., `!10 + 20` or `!Console.Cls()`) strip the `!` prefix, compile the snippet in isolation, execute it immediately on a fresh VM instance, and print the resulting top-of-stack value without altering the buffered program.

---

## 7. Notes for Grammar & Editor Tooling Authors

When building **TextMate** grammars or **Tree-sitter** parsers for PLAIN, note the following structural anchors:
- **Comments:** `'` starts a line comment (`comment.line.single-quote.plain`).
- **Block Delimiters:** All compound structures open with a keyword and close with a two-word `End <Keyword>` pair (`If`/`End If`, `For`/`End For`, `Select`/`End Select`, `Function`/`End Function`, `Struct`/`End Struct`).
- **Type Annotations:** Always preceded by the `As` keyword (`Identifier As TypeName`).
- **Built-in Namespaces:** Standard library module names act as root identifiers (`support.class.builtin.plain`).
- **String Literals:** Double-quoted (`"..."`) strings; `+` is overloaded for concatenation.