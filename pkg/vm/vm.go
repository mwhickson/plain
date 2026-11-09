package vm

import (
	"fmt"
	"plain/pkg/bytecode"
	"plain/pkg/value"
)

const StackMax = 1024

type VM struct {
	Chunk *bytecode.Chunk
	IP    int
	Stack [StackMax]value.Value
	SP    int
}

func NewVM(chunk *bytecode.Chunk) *VM {
	return &VM{
		Chunk: chunk,
		IP:    0,
		SP:    0,
	}
}

func (vm *VM) Push(val value.Value) error {
	if vm.SP >= StackMax {
		return fmt.Errorf("stack overflow")
	}
	vm.Stack[vm.SP] = val
	vm.SP++
	return nil
}

func (vm *VM) Pop() (value.Value, error) {
	if vm.SP == 0 {
		return value.Value{}, fmt.Errorf("stack underflow")
	}
	vm.SP--
	return vm.Stack[vm.SP], nil
}

func (vm *VM) Run() error {
	for vm.IP < len(vm.Chunk.Code) {
		op := bytecode.Opcode(vm.Chunk.Code[vm.IP])
		vm.IP++

		switch op {
		case bytecode.OpPushConst:
			constIdx := vm.Chunk.Code[vm.IP]
			vm.IP++
			if int(constIdx) >= len(vm.Chunk.Constants) {
				return fmt.Errorf("constant index %d out of bounds", constIdx)
			}
			if err := vm.Push(vm.Chunk.Constants[constIdx]); err != nil {
				return err
			}

		case bytecode.OpAdd:
			b, err := vm.Pop()
			if err != nil {
				return err
			}
			a, err := vm.Pop()
			if err != nil {
				return err
			}

			if a.Type == value.ValNumeric && b.Type == value.ValNumeric {
				if err := vm.Push(value.NewNumeric(a.Num + b.Num)); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("type error: '+' expects numeric operands, got %s and %s", a.Type, b.Type)
			}

		case bytecode.OpHalt:
			return nil

		default:
			return fmt.Errorf("unknown opcode %d at instruction pointer %d", op, vm.IP-1)
		}
	}
	return nil
}
