# PLAIN Virtual Machine Specification (v0.1)

This document defines the architecture, instruction set, execution model, and value representation for the PLAIN Virtual Machine (`plain-vm`).

## 1. Overview

The PLAIN Virtual Machine is a lightweight, stack-based runtime written in Go. It operates on compiled PLAIN Bytecode (`.pbc`) chunks, managing execution frames, global/local variable storage, closure upvalues, and memory for dynamic runtime values.

```text
       +------------------------------------+
       |          Bytecode Chunk            |
       |  - Code Stream: []byte             |
       |  - Constants Pool: []Value         |
       +------------------------------------+
                         |
                         v
+--------------------------------------------------+
|                    PLAIN VM                      |
|                                                  |
|  +--------------------+    +------------------+  |
|  |   Value Stack      |    |   Call Frames    |  |
|  |  [ Val1, Val2... ] |    |  - Closure Ptr   |  |
|  +--------------------+    |  - Inst Pointer  |  |
|                            |  - Base Pointer  |  |
|  +--------------------+    +------------------+  |
|  |   Global Slots     |                          |
|  |  0: System         |                          |
|  |  1: Console        |                          |
|  |  2..N: User Vars   |                          |
|  +--------------------+                          |
+--------------------------------------------------+
```

---

## 2. Value Representation

PLAIN utilizes dynamic runtime value wrappers (`value.Value`) to store scalar primitives, composite data structures, and executable objects on the evaluation stack and in variable tables.

### 2.1 Value Types
- `ValNone` — Represents an empty or default return value (`None`).
- `ValNumeric` — IEEE 754 64-bit floating-point number (`float64`).
- `ValString` — UTF-8 immutable string (`string`).
- `ValBoolean` — Boolean flag (`True` / `False`).
- `ValList` — Reference to a dynamic slice of `Value` objects (`[]value.Value`).
- `ValStruct` — Reference to a key-value map (`map[string]value.Value`), used for both user-defined structs and standard library module namespaces.
- `ValFunction` — Compiled function prototype (`*value.FunctionObject`) containing start IP, parameter count, and local variable count.
- `ValClosure` — Instantiated lexical closure (`*value.ClosureObject`) wrapping a `FunctionObject` and captured `Free` variables.
- `ValBuiltin` — Native Go function callback (`func(args []value.Value) (value.Value, error)`).

---

## 3. Bytecode Structure & Constants Pool

A compiled PLAIN program chunk consists of two primary structures:
1. **Constant Pool (`Constants`):** An indexed slice of static literal values (strings, numbers, function prototypes) referenced by 1-byte index.
2. **Instruction Stream (`Code`):** A sequential byte slice containing 1-byte Opcodes followed by zero or more 1-byte or 2-byte (little-endian `uint16`) operands.

---

## 4. Instruction Set Architecture (Opcodes)

All opcodes are 1-byte wide. Operands immediately follow their respective opcode in the instruction stream.

| Mnemonic | Operands | Stack Effect `(before -> after)` | Description |
| :--- | :--- | :--- | :--- |
| `OpPushConst` | `const_idx` (1 byte) | `[] -> [val]` | Pushes constant from `Constants[const_idx]` onto the stack. |
| `OpAdd` | None | `[a, b] -> [a + b]` | Pops `b` and `a`. Adds `Numeric` values or concatenates `String` values. |
| `OpSub` | None | `[a, b] -> [a - b]` | Pops `b` and `a`, pushes `a - b`. |
| `OpMul` | None | `[a, b] -> [a * b]` | Pops `b` and `a`, pushes `a * b`. |
| `OpDiv` | None | `[a, b] -> [a / b]` | Pops `b` and `a`, pushes `a / b`. Division by zero triggers runtime error. |
| `OpEq` | None | `[a, b] -> [bool]` | Pushes `True` if `a == b`, else `False`. |
| `OpNeq` | None | `[a, b] -> [bool]` | Pushes `True` if `a != b`, else `False`. |
| `OpLt` | None | `[a, b] -> [bool]` | Pushes `True` if `a < b`, else `False`. |
| `OpLte` | None | `[a, b] -> [bool]` | Pushes `True` if `a <= b`, else `False`. |
| `OpGt` | None | `[a, b] -> [bool]` | Pushes `True` if `a > b`, else `False`. |
| `OpGte` | None | `[a, b] -> [bool]` | Pushes `True` if `a >= b`, else `False`. |
| `OpStoreVar` | `var_idx` (1 byte) | `[val] -> []` | Pops top value and stores it into global slot `var_idx`. |
| `OpLoadVar` | `var_idx` (1 byte) | `[] -> [val]` | Pushes value from global slot `var_idx` onto the stack. |
| `OpStoreLocal` | `slot_idx` (1 byte) | `[val] -> []` | Pops top value and stores it at `frame.basePointer + slot_idx`. |
| `OpLoadLocal` | `slot_idx` (1 byte) | `[] -> [val]` | Pushes value from `frame.basePointer + slot_idx` onto the stack. |
| `OpLoadFree` | `free_idx` (1 byte) | `[] -> [val]` | Pushes captured closure variable from `frame.closure.Free[free_idx]`. |
| `OpBuildList` | `count` (1 byte) | `[e1..eN] -> [list]` | Pops `count` elements and constructs a `ValList`. |
| `OpBuildStruct` | `const_idx` (1 byte) | `[] -> [struct]` | Clones a struct default template into a fresh record value. |
| `OpIndex` | None | `[list, idx] -> [val]` | Indexes `list[idx]` and pushes the resulting element. |
| `OpGetProperty` | `const_idx` (1 byte) | `[struct] -> [val]` | Pops `struct`, looks up property named by `Constants[const_idx]`, pushes value. |
| `OpSetProperty` | `const_idx` (1 byte) | `[struct, val] -> []` | Pops `val` and `struct`, sets property named by `Constants[const_idx]` to `val`. |
| `OpJump` | `target` (2 bytes) | `[] -> []` | Unconditionally sets `IP` to 16-bit absolute target offset. |
| `OpJumpIfFalse` | `target` (2 bytes) | `[cond] -> []` | Pops condition; sets `IP` to `target` if falsy (`False`, `0`, or `None`). |
| `OpClosure` | `const_idx` (1B), `num_free` (1B) | `[f1..fN] -> [closure]` | Pops `num_free` captured values, wraps function at `const_idx` into `ValClosure`. |
| `OpCall` | `arg_count` (1 byte) | `[fn, a1..aN] -> [ret]` | Invokes `ValBuiltin`, `ValFunction`, or `ValClosure` with `arg_count` arguments. |
| `OpReturn` | None | `[ret] -> [ret]` | Pops active `Frame`, restores caller `IP` and stack pointer, pushes `ret`. |
| `OpHalt` | None | `[] -> []` | Cleanly terminates the virtual machine execution loop. |

---

## 5. Execution Model & Frame Management

```go
type Frame struct {
    closure     *value.ClosureObject
    ip          int
    basePointer int
}

type VM struct {
    chunk      *bytecode.Chunk
    ip         int
    stack      []value.Value
    globals    map[byte]value.Value
    frames     []Frame
    frameIndex int
}
```

### 5.1 Pre-Seeded Global Slots
Upon initialization (`NewVM`), the VM populates reserved global variable slots before script execution begins:
- **Slot `0`:** `System` module (`SetErrorHandler`, `Exit`)
- **Slot `1`:** `Console` module (`Print`, `PrintLine`, `Input`, `Cls`)

### 5.2 Runtime Error Interception (`handleError`)
If an instruction produces a runtime fault (such as division by zero, type mismatch, or out-of-bounds list indexing), the VM checks whether a handler has been registered via `System.SetErrorHandler`.
- If registered, the VM saves the active frame's `ip`, allocates a new `Frame` for the handler function/closure, passes the error message as a `ValString` argument, and resumes execution inside the handler.
- When the handler executes `OpReturn`, control resumes cleanly from the saved instruction pointer.

`System.Exit(code)` stops the VM without invoking the runtime error handler. The VM retains the requested code for the command-line runner, which exits with that status.