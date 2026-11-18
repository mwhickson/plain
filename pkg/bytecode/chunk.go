package bytecode

import (
	"fmt"
	"plain/pkg/value"
)

// Chunk represents a contiguous block of compiled bytecode instructions and constants.
type Chunk struct {
	Code      []byte
	Constants []value.Value
}

func NewChunk() *Chunk {
	return &Chunk{
		Code:      make([]byte, 0),
		Constants: make([]value.Value, 0),
	}
}

// AddConstant adds a value to the constants pool and returns its 1-byte index.
func (c *Chunk) AddConstant(val value.Value) (byte, error) {
	for i, existing := range c.Constants {
		if sameConstant(existing, val) {
			return byte(i), nil
		}
	}
	if len(c.Constants) > 255 {
		return 0, fmt.Errorf("constant pool exceeds 256 entries")
	}
	c.Constants = append(c.Constants, val)
	return byte(len(c.Constants) - 1), nil
}

func sameConstant(a, b value.Value) bool {
	if a.Type != b.Type {
		return false
	}
	switch a.Type {
	case value.ValNone:
		return true
	case value.ValBoolean:
		return a.Bool == b.Bool
	case value.ValNumeric:
		return a.Num == b.Num
	case value.ValString:
		return a.Str == b.Str
	default:
		return false
	}
}

// WriteByte appends an instruction byte to the chunk.
func (c *Chunk) WriteByte(b byte) error {
	c.Code = append(c.Code, b)
	return nil
}

// WriteOpcode appends an opcode byte to the chunk.
func (c *Chunk) WriteOpcode(op Opcode) {
	c.Code = append(c.Code, byte(op))
}
