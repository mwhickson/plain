# PLAIN Virtual Machine Specification (v0.1)

This document defines the architecture, instruction set, execution model, and value representation for the PLAIN Virtual Machine (`plain-vm`).

## 1. Overview

The PLAIN Virtual Machine is a lightweight, stack-based runtime written in Go. It operates on compiled PLAIN Bytecode (`.pbc`) chunks, managing execution frames, global/local variable storage, and memory for dynamic runtime values.

```
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
|  |   Value Stack      |    |    Call Frame    |  |
|  |  [ Val1, Val2... ] |    |  - Inst Pointer  |  |
|  +--------------------+    |  - Local Map     |  |
|                            +------------------+  |
+--------------------------------------------------+
```

---

## 2. Value Representation

PLAIN utilizes dynamic scalar types and explicit composite types. Values on the execution stack and in local/global storage are encapsulated in a tagged union structure:

### 2.1 Type Tags
- `0x00` — `ValNone`: Represents an empty or default return value.
- `0x01` — `ValNumeric`: IEEE 754 64-bit floating-point number.
- `0x02` — `ValString`: UTF-8 immutable string.
- `0x03` — `ValBoolean`: Boolean flag (`True` / `False`).
- `0x04` — `ValList`: Reference to a dynamic array of `Value` objects.
- `0x05` — `ValStruct`: Reference to a key-value map (`string -> Value`).

---

## 3. Bytecode Structure & Constants Pool

A compiled PLAIN program byte array consists of two primary tables:
1. **Constant Pool:** An indexed array of static literal values (strings, large numbers, predefined structures) referenced by byte index.
2. **Instruction Stream:** A sequential byte array containing Opcodes and their zero-or-more single-byte or multi-byte operands.

---

## 4. Instruction Set Architecture (Opcodes)

All opcodes are 1-byte wide (`uint8`). Operands immediately follow their respective opcode in the instruction stream.

| Hex Opcode | Mnemonic | Operands | Stack Diagram `(before -> after)` | Description |
| :--- | :--- | :--- | :--- | :--- |
| `0x00` | `OP_NOP` | None | `[] -> []` | No operation. |
| `0x01` | `OP_PUSH_CONST` | `const_idx` (1 byte) | `[] -> [val]` | Pushes constant from `Constants[const_idx]` to stack. |
| `0x02` | `OP_PUSH_NONE` | None | `[] -> [None]` | Pushes a `ValNone` onto stack. |
| `0x03` | `OP_PUSH_TRUE` | None | `[] -> [True]` | Pushes boolean `True`. |
| `0x04` | `OP_PUSH_FALSE` | None | `[] -> [False]` | Pushes boolean `False`. |
| `0x05` | `OP_POP` | None | `[val] -> []` | Removes top value from stack. |
| `0x10` | `OP_ADD` | None | `[b, a] -> [a + b]` | Pops `b` and `a`, pushes sum. Concatenates if strings. |
| `0x11` | `OP_SUB` | None | `[b, a] -> [a - b]` | Pops `b` and `a`, pushes `a - b`. |
| `0x12` | `OP_MUL` | None | `[b, a] -> [a * b]` | Pops `b` and `a`, pushes product. |
| `0x13` | `OP_DIV` | None | `[b, a] -> [a / b]` | Pops `b` and `a`, pushes quotient. Division by zero triggers runtime error. |
| `0x14` | `OP_MOD` | None | `[b, a] -> [a % b]` | Pops `b` and `a`, pushes modulo. |
| `0x15` | `OP_NEGATE` | None | `[a] -> [-a]` | Negates top numeric value on stack. |
| `0x20` | `OP_EQ` | None | `[b, a] -> [bool]` | Pushes `True` if `a == b`, else `False`. |
| `0x21` | `OP_NEQ` | None | `[b, a] -> [bool]` | Pushes `True` if `a != b`, else `False`. |
| `0x22` | `OP_LT` | None | `[b, a] -> [bool]` | Pushes `True` if `a < b`, else `False`. |
| `0x23` | `OP_LTE` | None | `[b, a] -> [bool]` | Pushes `True` if `a <= b`, else `False`. |
| `0x24` | `OP_GT` | None | `[b, a] -> [bool]` | Pushes `True` if `a > b`, else `False`. |
| `0x25` | `OP_GTE` | None | `[b, a] -> [bool]` | Pushes `True` if `a >= b`, else `False`. |
| `0x30` | `OP_STORE_VAR` | `var_idx` (1 byte) | `[val] -> []` | Stores top value into variable table slot `var_idx`. |
| `0x31` | `OP_LOAD_VAR` | `var_idx` (1 byte) | `[] -> [val]` | Loads value from variable table slot `var_idx` to stack. |
| `0x40` | `OP_JUMP` | `offset` (2 bytes) | `[] -> []` | Unconditionally jumps instruction pointer forward. |
| `0x41` | `OP_JUMP_IF_FALSE` | `offset` (2 bytes) | `[bool] -> []` | Pops stack value; jumps if value is `False`. |
| `0x50` | `OP_CALL` | `arg_count` (1 byte) | `[fn, arg1...]` | Creates new CallFrame and transfers execution to `fn`. |
| `0x51` | `OP_RETURN` | None | `[ret] -> [ret]` | Destroys current CallFrame, returning control to caller. |
| `0xFF` | `OP_HALT` | None | `[] -> []` | Terminate execution of virtual machine loop. |

---

## 5. Execution Loop Model

The core VM execution loop operates on an Instruction Pointer (`IP`) referencing the currently active `CallFrame`:

```go
type VM struct {
    frames     [64]CallFrame
    frameCount int
    stack      [1024]Value
    sp         int // Stack Pointer
}
```

- **Stack Overflow / Underflow Checks:** Evaluated dynamically during execution.
- **Pass-By-Value Handling:** Composite structures (`List`, `Struct`) passed as arguments or assignments trigger a copy-on-write or deep duplicate to guarantee pure value-semantics at runtime.