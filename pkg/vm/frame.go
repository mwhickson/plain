package vm

import "plain/pkg/value"

type Frame struct {
	closure     *value.ClosureObject
	ip          int // Saved instruction pointer
	basePointer int // Base pointer in the evaluation stack for local variables
}
