package bytecode

import "plain/pkg/value"

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

// AddConstant adds a literal Value to the constants pool and returns its index.
func (c *Chunk) AddConstant(val value.Value) byte {
	c.Constants = append(c.Constants, val)
	return byte(len(c.Constants) - 1)
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
