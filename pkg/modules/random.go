package modules

import (
	"fmt"
	"math"
	"math/rand"
	"plain/pkg/value"
	"sync"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))
var rngMu sync.Mutex

func GetRandomModule() value.Value {
	randomModule := map[string]value.Value{
		"Seed":    value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(randomSeed)},
		"Int":     value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(randomInt)},
		"Float":   value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(randomFloat)},
		"Choice":  value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(randomChoice)},
		"Chance":  value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(randomChance)},
		"Shuffle": value.Value{Type: value.ValBuiltin, Obj: value.BuiltinFunction(randomShuffle)},
	}
	return value.Value{
		Type: value.ValStruct,
		Obj:  randomModule,
	}
}

func randomSeed(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Type != value.ValNumeric {
		return value.NewNone(), fmt.Errorf("Random.Seed() requires 1 Numeric argument")
	}
	if math.IsNaN(args[0].Num) || math.IsInf(args[0].Num, 0) || math.Trunc(args[0].Num) != args[0].Num ||
		args[0].Num < -math.Ldexp(1, 63) || args[0].Num >= math.Ldexp(1, 63) {
		return value.NewNone(), fmt.Errorf("Random.Seed() seed must be a finite signed 64-bit integer")
	}
	rngMu.Lock()
	rng = rand.New(rand.NewSource(int64(args[0].Num)))
	rngMu.Unlock()
	return value.NewNone(), nil
}

func randomInt(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Type != value.ValNumeric || args[1].Type != value.ValNumeric {
		return value.NewNone(), fmt.Errorf("Random.Int() requires (Numeric min, Numeric max)")
	}
	minValue, maxValue := args[0].Num, args[1].Num
	const int64Limit = 9223372036854775808.0
	if math.IsNaN(minValue) || math.IsInf(minValue, 0) || math.Trunc(minValue) != minValue ||
		math.IsNaN(maxValue) || math.IsInf(maxValue, 0) || math.Trunc(maxValue) != maxValue ||
		minValue < -int64Limit || minValue >= int64Limit || maxValue < -int64Limit || maxValue >= int64Limit {
		return value.NewNone(), fmt.Errorf("Random.Int() bounds must be finite signed 64-bit integers")
	}
	min := int64(minValue)
	max := int64(maxValue)
	if min > max {
		min, max = max, min
	}
	rangeSize := uint64(max) - uint64(min) + 1
	if rangeSize == 0 || rangeSize > math.MaxInt64 {
		return value.NewNone(), fmt.Errorf("Random.Int() range is too large")
	}
	rngMu.Lock()
	val := min + rng.Int63n(int64(rangeSize))
	rngMu.Unlock()
	return value.NewNumeric(float64(val)), nil
}

func randomFloat(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.NewNone(), fmt.Errorf("Random.Float() takes 0 arguments")
	}
	rngMu.Lock()
	val := rng.Float64()
	rngMu.Unlock()
	return value.NewNumeric(val), nil
}

func randomChoice(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Type != value.ValList {
		return value.NewNone(), fmt.Errorf("Random.Choice() requires 1 List argument")
	}
	list := args[0].Obj.([]value.Value)
	if len(list) == 0 {
		return value.NewNone(), nil
	}
	rngMu.Lock()
	idx := rng.Intn(len(list))
	rngMu.Unlock()
	return list[idx], nil
}

func randomChance(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Type != value.ValNumeric {
		return value.NewNone(), fmt.Errorf("Random.Chance() requires 1 Numeric percentage argument (0-100)")
	}
	pct := args[0].Num
	if math.IsNaN(pct) || pct < 0 || pct > 100 {
		return value.NewNone(), fmt.Errorf("Random.Chance() percentage must be between 0 and 100")
	}
	rngMu.Lock()
	roll := rng.Float64() * 100.0
	rngMu.Unlock()
	return value.NewBoolean(roll < pct), nil
}

func randomShuffle(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Type != value.ValList {
		return value.NewNone(), fmt.Errorf("Random.Shuffle() requires 1 List argument")
	}
	original := args[0].Obj.([]value.Value)
	shuffled := make([]value.Value, len(original))
	copy(shuffled, original)
	rngMu.Lock()
	rng.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	rngMu.Unlock()
	return value.Value{Type: value.ValList, Obj: shuffled}, nil
}
