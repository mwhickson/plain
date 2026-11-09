# PLAIN Language Specification (v0.1)

**PLAIN** (Programming Language Allowing Implementation Neatly) is an approachable, modern BASIC-inspired programming language built for fast, simple desktop applications and general-purpose tools.

This document serves as the official specification for the language syntax, type system, execution model, and core standard library design.

---

## 1. Core Principles

1. **One Obvious Way:** Reduce redundant language constructs. Prefer clean, unified abstractions over syntactic sugar.
2. **Predictable Semantics:** Pass-by-value semantics by default. No raw pointers, manual pointer arithmetic, or explicit reference parameters (`ByRef`).
3. **Batteries Included:** Rich standard library covering native GUI bindings, networking, persistent storage (SQLite), and standard data formats (JSON, CSV).
4. **Sandboxed Interoperability:** External system integration occurs via standard interfaces (OS signals, standard IO streams, network APIs, local files) rather than unsafe raw foreign function interface (FFI) memory access.

---

## 2. Type System

PLAIN features a small, opinionated set of primitive and composite types:

| Type Name | Description | Built-in Representation |
| :--- | :--- | :--- |
| `String` | UTF-8 encoded sequence of text characters | Native UTF-8 string |
| `Numeric` | Unified number type for integer and floating-point math | 64-bit float (`double`) / High-precision decimal |
| `Boolean` | Truth value (`True` or `False`) | 1-bit boolean |
| `None` | Represents the absence of a value | Unit / None singleton |
| `List` | Dynamic, ordered sequence of elements | Dynamic array |
| `Struct` | Heterogeneous key-value record and data container | Map / Named record |

### 2.1 Mutability and Declarations
Variables are explicitly bound using `Var` (mutable) or `Const` (immutable).

```plain
Const AppName As String = "PLAIN Desktop"
Var UserCount As Numeric = 42
```

---

## 3. Syntax Overview

### 3.1 Identifiers & Keywords
- Keywords are case-insensitive (e.g., `Function`, `function`, and `FUNCTION` are equivalent), though **PascalCase** is idiomatic.
- Statements are line-delimited (newline character). Multiple statements per line are not permitted.

### 3.2 Functions
Functions are the sole executable block unit (there are no separate `Sub` or `Procedure` keywords). 
- If a function has no explicit return statement, or returns without an expression, it returns `None`.

```plain
Function Add(a As Numeric, b As Numeric) As Numeric
    Return a + b
End Function

Function Greet(name As String)
    Print("Hello, " + name)
    ' Implicitly returns None
End Function
```

### 3.3 Method Binding on Structs
Functions can be bound to a target `Struct` as receiver methods (similar to Go):

```plain
Struct Vector2D
    X As Numeric
    Y As Numeric
End Struct

Function (v As Vector2D) Multiply(factor As Numeric)
    v.X = v.X * factor
    v.Y = v.Y * factor
End Function
```

### 3.4 Control Flow

#### If Statements
```plain
If count > 10 Then
    Print("High")
ElseIf count > 5 Then
    Print("Medium")
Else
    Print("Low")
End If
```

#### For Loops (Unified)
PLAIN features a single `For` construct covering both numeric ranges and list iterations:

```plain
' Range Loop
For i As Numeric = 1 To 10 Step 1
    Print(i)
End For

' Iteration Loop
Var items As List = ["Apple", "Banana", "Cherry"]
For item In items
    Print(item)
End For
```

#### Select Case Statements
Multi-branch condition matching uses structured `Select Case` blocks:

```plain
Select Case status
    Case "Pending"
        HandlePending()
    Case "Active", "Running"
        Execute()
    Case Else
        HandleUnknown()
End Select
```

---

## 4. Grammar Specification (EBNF)

Below is the formal Extended Backus-Naur Form (EBNF) grammar defining PLAIN syntax.

```ebnf
Program          = { Statement | FunctionDecl | StructDecl } ;

Statement        = VarDecl
                 | Assignment
                 | IfStmt
                 | ForStmt
                 | SelectStmt
                 | ReturnStmt
                 | ExpressionStmt ;

(* Declarations *)
VarDecl          = ( "Var" | "Const" ) Identifier [ "As" TypeName ] "=" Expression ;
StructDecl       = "Struct" Identifier { MemberDecl } "End" "Struct" ;
MemberDecl       = Identifier "As" TypeName ;

FunctionDecl     = "Function" [ MethodReceiver ] Identifier "(" [ ParameterList ] ")" [ "As" TypeName ] 
                   { Statement } 
                   "End" "Function" ;

MethodReceiver   = "(" Identifier "As" TypeName ")" ;
ParameterList    = Parameter { "," Parameter } ;
Parameter        = Identifier "As" TypeName ;

(* Control Flow *)
IfStmt           = "If" Expression "Then" { Statement }
                   { "ElseIf" Expression "Then" { Statement } }
                   [ "Else" { Statement } ]
                   "End" "If" ;

ForStmt          = "For" ( RangeLoop | IterLoop ) { Statement } "End" "For" ;
RangeLoop        = Identifier "As" TypeName "=" Expression "To" Expression [ "Step" Expression ] ;
IterLoop         = Identifier "In" Expression ;

SelectStmt       = "Select" "Case" Expression 
                   { CaseClause } 
                   [ "Case" "Else" { Statement } ] 
                   "End" "Select" ;
CaseClause       = "Case" ExpressionList { Statement } ;

ReturnStmt       = "Return" [ Expression ] ;
ExpressionStmt   = Expression ;

(* Types & Expressions *)
TypeName         = "String" | "Numeric" | "Boolean" | "List" | "Struct" | Identifier ;
ExpressionList   = Expression { "," Expression } ;
Expression       = LogicalOr ;
LogicalOr        = LogicalAnd { "Or" LogicalAnd } ;
LogicalAnd       = Equality { "And" Equality } ;
Equality         = Relational { ( "==" | "!=" ) Relational } ;
Relational       = Additive { ( "<" | "<=" | ">" | ">=" ) Additive } ;
Additive         = Multiplicative { ( "+" | "-" ) Multiplicative } ;
Multiplicative   = Unary { ( "*" | "/" | "%" ) Unary } ;
Unary            = [ "-" | "Not" ] Primary ;
Primary          = Literal 
                 | Identifier 
                 | CallOrAccess 
                 | "(" Expression ")" ;

CallOrAccess     = Primary ( "(" [ ExpressionList ] ")" | "." Identifier | "[" Expression "]" ) ;

Literal          = NumericLiteral | StringLiteral | BooleanLiteral | ListLiteral | "None" ;
NumericLiteral   = digit { digit } [ "." digit { digit } ] ;
StringLiteral    = '"' { string_char } '"' ;
BooleanLiteral   = "True" | "False" ;
ListLiteral      = "[" [ ExpressionList ] "]" ;

Identifier       = letter { letter | digit | "_" } ;
```

---

## 5. Execution Model

1. **Compilation Strategy:** Source code is parsed into an Abstract Syntax Tree (AST) and compiled into compact PLAIN Bytecode (`.pbc`).
2. **Virtual Machine (VM):** The bytecode is executed by a lightweight stack-based or register-based VM written in Go.
3. **Memory Management:** Memory allocation for composite objects (`List`, `Struct`) is handled via underlying host VM garbage collection. All user code is executed under strict pass-by-value guarantees.