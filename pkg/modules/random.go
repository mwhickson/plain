package modules

import (
	"fmt"
	"math/rand"
	"plain/pkg/value"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

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
	rng = rand.New(rand.NewSource(int64(args[0].Num)))
	return value.NewNone(), nil
}

func randomInt(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Type != value.ValNumeric || args[1].Type != value.ValNumeric {
		return value.NewNone(), fmt.Errorf("Random.Int() requires (Numeric min, Numeric max)")
	}
	min := int(args[0].Num)
	max := int(args[1].Num)
	if min > max {
		min, max = max, min
	}
	val := min + rng.Intn(max-min+1)
	return value.NewNumeric(float64(val)), nil
}

func randomFloat(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.NewNone(), fmt.Errorf("Random.Float() takes 0 arguments")
	}
	return value.NewNumeric(rng.Float64()), nil
}

func randomChoice(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Type != value.ValList {
		return value.NewNone(), fmt.Errorf("Random.Choice() requires 1 List argument")
	}
	list := args[0].Obj.([]value.Value)
	if len(list) == 0 {
		return value.NewNone(), nil
	}
	idx := rng.Intn(len(list))
	return list[idx], nil
}

func randomChance(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Type != value.ValNumeric {
		return value.NewNone(), fmt.Errorf("Random.Chance() requires 1 Numeric percentage argument (0-100)")
	}
	pct := args[0].Num
	roll := rng.Float64() * 100.0
	return value.NewBoolean(roll < pct), nil
}

func randomShuffle(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Type != value.ValList {
		return value.NewNone(), fmt.Errorf("Random.Shuffle() requires 1 List argument")
	}
	original := args[0].Obj.([]value.Value)
	shuffled := make([]value.Value, len(original))
	copy(shuffled, original)
	rng.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return value.Value{Type: value.ValList, Obj: shuffled}, nil
}