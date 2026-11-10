package vm

import "plain/pkg/value"

type Frame struct {
	fn          *value.FunctionObject
	ip          int // Saved instruction pointer
	basePointer int // Base pointer in the evaluation stack for local variables
}
