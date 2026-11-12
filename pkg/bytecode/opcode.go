package bytecode

// Opcode represents a single bytecode operation in the PLAIN VM.
type Opcode byte

const (
	OpNop Opcode = iota

	OpAdd
	OpCall
	OpClosure
	OpDiv
	OpEq
	OpGt
	OpGte
	OpJump
	OpJumpIfFalse
	OpLoadFree
	OpLoadLocal
	OpLoadVar
	OpLt
	OpLte
	OpMod
	OpMul
	OpNegate
	OpNeq
	OpPop
	OpPushConst
	OpPushFalse
	OpPushNone
	OpPushTrue
	OpReturn
	OpStoreLocal
	OpStoreVar
	OpSub

	// System / Execution
	OpHalt Opcode = 0xFF
)

// String returns the mnemonic name for the given Opcode.
func (o Opcode) String() string {
	switch o {
	case OpNop:
		return "OP_NOP"
	case OpAdd:
		return "OP_ADD"
	case OpCall:
		return "OP_CALL"
	case OpClosure:
		return "OP_CLOSURE"
	case OpDiv:
		return "OP_DIV"
	case OpEq:
		return "OP_EQ"
	case OpGt:
		return "OP_GT"
	case OpGte:
		return "OP_GTE"
	case OpJump:
		return "OP_JUMP"
	case OpJumpIfFalse:
		return "OP_JUMP_IF_FALSE"
	case OpLoadFree:
		return "OP_LOAD_FREE"
	case OpLoadLocal:
		return "OP_LOAD_LOCAL"
	case OpLoadVar:
		return "OP_LOAD_VAR"
	case OpLt:
		return "OP_LT"
	case OpLte:
		return "OP_LTE"
	case OpMod:
		return "OP_MOD"
	case OpMul:
		return "OP_MUL"
	case OpNegate:
		return "OP_NEGATE"
	case OpNeq:
		return "OP_NEQ"
	case OpPop:
		return "OP_POP"
	case OpPushConst:
		return "OP_PUSH_CONST"
	case OpPushFalse:
		return "OP_PUSH_FALSE"
	case OpPushNone:
		return "OP_PUSH_NONE"
	case OpPushTrue:
		return "OP_PUSH_TRUE"
	case OpReturn:
		return "OP_RETURN"
	case OpStoreLocal:
		return "OP_STORE_LOCAL"
	case OpStoreVar:
		return "OP_STORE_VAR"
	case OpSub:
		return "OP_SUB"
	case OpHalt:
		return "OP_HALT"
	default:
		return "OP_UNKNOWN"
	}
}
