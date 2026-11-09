# PLAIN Programming Language

PLAIN (Programming Language Allowing Implementation Neatly) is a BASIC-like/inspired-by-BASIC programming language.

PLAIN, and its associated tooling, is intended to fill software development niches historically attributed to dBase, Delphi, FoxPro, HyperCard, MS Access, PowerBuilder, Xojo, and Visual Basic.

> NOTE: PLAIN is a working-title backronym.

## Package Structure

```
plain/
├── cmd/
│   └── plain/
│       └── main.go          # CLI entry point (repl, run file, disassemble)
├── pkg/
│   ├── ast/
│   │   └── ast.go           # Abstract Syntax Tree node definitions
│   ├── bytecode/
│   │   ├── opcode.go        # Opcode definitions and string representations
│   │   ├── chunk.go         # Bytecode chunk & constant pool containers
│   │   └── disassembler.go  # Helper to print bytecode human-readably (for debugging!)
│   ├── compiler/
│   │   └── compiler.go      # AST -> Bytecode emitter
│   ├── parser/
│   │   ├── lexer.go         # Tokenizer (source text -> tokens)
│   │   └── parser.go        # Syntax analyzer (tokens -> AST)
│   ├── value/
│   │   └── value.go         # PLAIN runtime values (Numeric, String, List, Struct)
│   └── vm/
│       ├── frame.go         # Call frames and scope management
│       └── vm.go            # Core bytecode execution loop
├── go.mod
└── README.md
```

## Disclaimer

This application was produced using AI (Gemini).
