package modules

import (
	"fmt"
	"math"
	"plain/pkg/value"
	"time"
)

func GetTimeModule() value.Value {
	timeModule := map[string]value.Value{
		"Now":    value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(timeNow)},
		"Sleep":  value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(timeSleep)},
		"Format": value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(timeFormat)},
	}
	return value.Value{
		Type: value.ValStruct,
		Obj:  timeModule,
	}
}

func timeNow(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.NewNone(), fmt.Errorf("Time.Now() takes 0 arguments")
	}
	return value.NewNumeric(float64(time.Now().UnixMilli())), nil
}

func timeSleep(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Type != value.ValNumeric {
		return value.NewNone(), fmt.Errorf("Time.Sleep() requires 1 Numeric argument (milliseconds)")
	}
	ms := args[0].Num
	durationNanoseconds := ms * float64(time.Millisecond)
	if math.IsNaN(ms) || math.IsInf(ms, 0) || ms < 0 || durationNanoseconds >= math.Ldexp(1, 63) {
		return value.NewNone(), fmt.Errorf("Time.Sleep() milliseconds must be finite, non-negative, and within the supported duration")
	}
	time.Sleep(time.Duration(durationNanoseconds))
	return value.NewNone(), nil
}

func timeFormat(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Type != value.ValNumeric {
		return value.NewNone(), fmt.Errorf("Time.Format() requires 1 Numeric timestamp (milliseconds)")
	}
	timestamp := args[0].Num
	if math.IsNaN(timestamp) || math.IsInf(timestamp, 0) || timestamp < -math.Ldexp(1, 63) || timestamp >= math.Ldexp(1, 63) {
		return value.NewNone(), fmt.Errorf("Time.Format() timestamp must be a finite signed 64-bit integer")
	}
	t := time.UnixMilli(int64(timestamp))
	return value.NewString(t.Format("2006-01-02 15:04:05")), nil
}
