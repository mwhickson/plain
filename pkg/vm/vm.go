package vm

import (
	"fmt"
	"plain/pkg/bytecode"
	"plain/pkg/value"
)

type VM struct {
	chunk   *bytecode.Chunk
	ip      int
	stack   []value.Value
	globals map[byte]value.Value // Variable storage slot
}

func NewVM(chunk *bytecode.Chunk) *VM {
	return &VM{
		chunk:   chunk,
		ip:      0,
		stack:   make([]value.Value, 0, 256),
		globals: make(map[byte]value.Value),
	}
}

func (vm *VM) Push(val value.Value) {
	vm.stack = append(vm.stack, val)
}

func (vm *VM) Pop() (value.Value, error) {
	length := len(vm.stack)
	if length == 0 {
		return value.Value{}, fmt.Errorf("stack underflow")
	}
	val := vm.stack[length-1]
	vm.stack = vm.stack[:length-1]
	return val, nil
}

func (vm *VM) Run() error {
	for vm.ip < len(vm.chunk.Code) {
		op := bytecode.Opcode(vm.chunk.Code[vm.ip])
		vm.ip++

		switch op {
		case bytecode.OpStoreVar:
			dataIndex := vm.chunk.Code[vm.ip]
			vm.ip++
			val, err := vm.Pop()
			if err != nil {
				return err
			}
			vm.globals[dataIndex] = val

		case bytecode.OpLoadVar:
			dataIndex := vm.chunk.Code[vm.ip]
			vm.ip++
			val, exists := vm.globals[dataIndex]
			if !exists {
				return fmt.Errorf("undefined variable index %d", dataIndex)
			}
			vm.Push(val)

		case bytecode.OpPushConst:
			constIdx := vm.chunk.Code[vm.ip]
			vm.ip++
			if int(constIdx) >= len(vm.chunk.Constants) {
				return fmt.Errorf("constant index %d out of bounds", constIdx)
			}
			vm.Push(vm.chunk.Constants[constIdx])

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
				vm.Push(value.NewNumeric(a.Num + b.Num))
			} else {
				return fmt.Errorf("type error: '+' expects numeric operands, got %s and %s", a.Type, b.Type)
			}

		case bytecode.OpSub:
			b, err := vm.Pop()
			if err != nil {
				return err
			}
			a, err := vm.Pop()
			if err != nil {
				return err
			}
			if a.Type == value.ValNumeric && b.Type == value.ValNumeric {
				vm.Push(value.NewNumeric(a.Num - b.Num))
			} else {
				return fmt.Errorf("type error: '-' expects numeric operands")
			}

		case bytecode.OpMul:
			b, err := vm.Pop()
			if err != nil {
				return err
			}
			a, err := vm.Pop()
			if err != nil {
				return err
			}
			if a.Type == value.ValNumeric && b.Type == value.ValNumeric {
				vm.Push(value.NewNumeric(a.Num * b.Num))
			} else {
				return fmt.Errorf("type error: '*' expects numeric operands")
			}

		case bytecode.OpDiv:
			b, err := vm.Pop()
			if err != nil {
				return err
			}
			a, err := vm.Pop()
			if err != nil {
				return err
			}
			if a.Type == value.ValNumeric && b.Type == value.ValNumeric {
				if b.Num == 0 {
					return fmt.Errorf("division by zero")
				}
				vm.Push(value.NewNumeric(a.Num / b.Num))
			} else {
				return fmt.Errorf("type error: '/' expects numeric operands")
			}

		case bytecode.OpEq:
			b, err := vm.Pop()
			if err != nil {
				return err
			}
			a, err := vm.Pop()
			if err != nil {
				return err
			}
			// Simple equality evaluation
			equal := (a.Type == b.Type && a.Num == b.Num)
			var val float64
			if equal {
				val = 1
			} else {
				val = 0
			}
			vm.Push(value.NewNumeric(val))

		case bytecode.OpHalt:
			return nil

		case bytecode.OpJump:
			target := int(vm.chunk.Code[vm.ip]) | int(vm.chunk.Code[vm.ip+1])<<8
			vm.ip = target

		case bytecode.OpJumpIfFalse:
			target := int(vm.chunk.Code[vm.ip]) | int(vm.chunk.Code[vm.ip+1])<<8
			vm.ip += 2 // advance past the operand bytes

			condition, err := vm.Pop()
			if err != nil {
				return err
			}

			isFalse := (condition.Type == value.ValNumeric && condition.Num == 0)
			if isFalse {
				vm.ip = target
			}

		case bytecode.OpLte:
			b, err := vm.Pop()
			if err != nil {
				return err
			}
			a, err := vm.Pop()
			if err != nil {
				return err
			}
			lte := (a.Type == value.ValNumeric && b.Type == value.ValNumeric && a.Num <= b.Num)
			var val float64
			if lte {
				val = 1
			} else {
				val = 0
			}
			vm.Push(value.NewNumeric(val))

		default:
			return fmt.Errorf("unknown opcode %d at instruction pointer %d", op, vm.ip-1)
		}
	}
	return nil
}
