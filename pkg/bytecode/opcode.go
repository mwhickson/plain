package bytecode

// Opcode represents a single bytecode operation in the PLAIN VM.
type Opcode byte

const (
	OpNop Opcode = iota
	OpPushConst
	OpPushNone
	OpPushTrue
	OpPushFalse
	OpPop

	// Arithmetic Operations
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpMod
	OpNegate

	// Comparison Operations
	OpEq
	OpNeq
	OpLt
	OpLte
	OpGt
	OpGte

	// Variables
	OpStoreVar
	OpLoadVar

	// Control Flow
	OpJump
	OpJumpIfFalse

	// Functions & Calls
	OpCall
	OpReturn
	OpLoadLocal
	OpStoreLocal

	// System / Execution
	OpHalt Opcode = 0xFF
)

// String returns the mnemonic name for the given Opcode.
func (o Opcode) String() string {
	switch o {
	case OpNop:
		return "OP_NOP"
	case OpPushConst:
		return "OP_PUSH_CONST"
	case OpPushNone:
		return "OP_PUSH_NONE"
	case OpPushTrue:
		return "OP_PUSH_TRUE"
	case OpPushFalse:
		return "OP_PUSH_FALSE"
	case OpPop:
		return "OP_POP"
	case OpAdd:
		return "OP_ADD"
	case OpSub:
		return "OP_SUB"
	case OpMul:
		return "OP_MUL"
	case OpDiv:
		return "OP_DIV"
	case OpMod:
		return "OP_MOD"
	case OpNegate:
		return "OP_NEGATE"
	case OpEq:
		return "OP_EQ"
	case OpNeq:
		return "OP_NEQ"
	case OpLt:
		return "OP_LT"
	case OpLte:
		return "OP_LTE"
	case OpGt:
		return "OP_GT"
	case OpGte:
		return "OP_GTE"
	case OpStoreVar:
		return "OP_STORE_VAR"
	case OpLoadVar:
		return "OP_LOAD_VAR"
	case OpJump:
		return "OP_JUMP"
	case OpJumpIfFalse:
		return "OP_JUMP_IF_FALSE"
	case OpCall:
		return "OP_CALL"
	case OpReturn:
		return "OP_RETURN"
	case OpHalt:
		return "OP_HALT"
	default:
		return "OP_UNKNOWN"
	}
}
