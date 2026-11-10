package vm

import (
	"fmt"
	"plain/pkg/bytecode"
	"plain/pkg/value"
)

type VM struct {
	chunk      *bytecode.Chunk
	ip         int
	stack      []value.Value
	globals    map[byte]value.Value // Variable storage slot
	frames     []Frame
	frameIndex int
}

func NewVM(chunk *bytecode.Chunk) *VM {
	vm := &VM{
		chunk:      chunk,
		ip:         0,
		stack:      make([]value.Value, 0, 256),
		globals:    make(map[byte]value.Value),
		frames:     make([]Frame, 1024),
		frameIndex: 1,
	}
	// Setup root frame for main script
	vm.frames[0] = Frame{
		fn:          nil,
		ip:          0,
		basePointer: 0,
	}
	return vm
}

func (vm *VM) Globals() map[byte]value.Value {
	return vm.globals
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

func (vm *VM) Push(val value.Value) {
	vm.stack = append(vm.stack, val)
}

func (vm *VM) Run() error {
	for vm.ip < len(vm.chunk.Code) {
		op := bytecode.Opcode(vm.chunk.Code[vm.ip])
		vm.ip++

		switch op {
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

		case bytecode.OpCall:
			numArgs := int(vm.chunk.Code[vm.ip])
			vm.ip++

			// Pop arguments in reverse order
			args := make([]value.Value, numArgs)
			for i := numArgs - 1; i >= 0; i-- {
				arg, err := vm.Pop()
				if err != nil {
					return err
				}
				args[i] = arg
			}

			// Now pop the function object itself
			callee, err := vm.Pop()
			if err != nil {
				return err
			}

			if callee.Type != value.ValFunction {
				return fmt.Errorf("type error: can only call functions, got %s", callee.Type)
			}

			fnObj := callee.Obj.(*value.FunctionObject)
			if numArgs != fnObj.ParamCount {
				return fmt.Errorf("wrong number of arguments: want %d, got %d", fnObj.ParamCount, numArgs)
			}

			// Save current frame's IP before switching
			vm.frames[vm.frameIndex-1].ip = vm.ip

			// Set base pointer for the new frame where arguments will reside on the stack
			basePtr := len(vm.stack)
			for _, arg := range args {
				vm.Push(arg)
			}

			// Push new call frame
			newFrame := Frame{
				fn:          fnObj,
				ip:          fnObj.InstructionStartIndex,
				basePointer: basePtr,
			}
			vm.frames[vm.frameIndex] = newFrame
			vm.frameIndex++

			// Jump VM instruction pointer to function start
			vm.ip = fnObj.InstructionStartIndex

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
			equal := (a.Type == b.Type && a.Num == b.Num && a.Str == b.Str)
			vm.Push(value.NewBoolean(equal))

		case bytecode.OpHalt:
			return nil

		case bytecode.OpJump:
			target := int(vm.chunk.Code[vm.ip]) | int(vm.chunk.Code[vm.ip+1])<<8
			vm.ip = target

		case bytecode.OpJumpIfFalse:
			target := int(vm.chunk.Code[vm.ip]) | int(vm.chunk.Code[vm.ip+1])<<8
			vm.ip += 2

			condition, err := vm.Pop()
			if err != nil {
				return err
			}

			isFalse := false
			if condition.Type == value.ValBoolean {
				isFalse = (condition.Num == 0) // assuming boolean stored 0/1 in Num
			} else if condition.Type == value.ValNumeric {
				isFalse = (condition.Num == 0)
			} else if condition.Type == value.ValNone {
				isFalse = true
			}

			if isFalse {
				vm.ip = target
			}

		case bytecode.OpLoadLocal:
			slotIndex := vm.chunk.Code[vm.ip]
			vm.ip++
			currentFrame := vm.frames[vm.frameIndex-1]
			val := vm.stack[currentFrame.basePointer+int(slotIndex)]
			vm.Push(val)

		case bytecode.OpLoadVar:
			dataIndex := vm.chunk.Code[vm.ip]
			vm.ip++
			val, exists := vm.globals[dataIndex]
			if !exists {
				return fmt.Errorf("undefined variable index %d", dataIndex)
			}
			vm.Push(val)

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
			vm.Push(value.NewBoolean(lte))

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

		case bytecode.OpPushConst:
			constIdx := vm.chunk.Code[vm.ip]
			vm.ip++
			if int(constIdx) >= len(vm.chunk.Constants) {
				return fmt.Errorf("constant index %d out of bounds", constIdx)
			}
			vm.Push(vm.chunk.Constants[constIdx])

		case bytecode.OpReturn:
			// Pop return value
			retVal, err := vm.Pop()
			if err != nil {
				return err
			}

			// Pop current frame
			vm.frameIndex--
			currentFrame := vm.frames[vm.frameIndex]

			// Reset stack pointer back to frame base pointer
			vm.stack = vm.stack[:currentFrame.basePointer]

			// Restore previous frame IP
			previousFrame := vm.frames[vm.frameIndex-1]
			vm.ip = previousFrame.ip

			// Push return value onto the stack
			vm.Push(retVal)

		case bytecode.OpStoreLocal:
			slotIndex := vm.chunk.Code[vm.ip]
			vm.ip++
			val, err := vm.Pop()
			if err != nil {
				return err
			}
			currentFrame := vm.frames[vm.frameIndex-1]
			vm.stack[currentFrame.basePointer+int(slotIndex)] = val

		case bytecode.OpStoreVar:
			dataIndex := vm.chunk.Code[vm.ip]
			vm.ip++
			val, err := vm.Pop()
			if err != nil {
				return err
			}
			vm.globals[dataIndex] = val

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

		default:
			return fmt.Errorf("unknown opcode %d at instruction pointer %d", op, vm.ip-1)
		}
	}
	return nil
}
