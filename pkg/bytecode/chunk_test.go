package bytecode_test

import (
	"plain/pkg/bytecode"
	"plain/pkg/value"
	"testing"
)

func TestChunkAddConstantDeduplicatesAndEnforcesByteIndexLimit(t *testing.T) {
	chunk := bytecode.NewChunk()
	first, err := chunk.AddConstant(value.NewString("shared"))
	if err != nil {
		t.Fatal(err)
	}
	again, err := chunk.AddConstant(value.NewString("shared"))
	if err != nil {
		t.Fatal(err)
	}
	if first != again || len(chunk.Constants) != 1 {
		t.Fatalf("duplicate string indices = (%d, %d), pool size = %d; want one shared entry", first, again, len(chunk.Constants))
	}

	for i := 1; i < 256; i++ {
		index, err := chunk.AddConstant(value.NewNumeric(float64(i)))
		if err != nil {
			t.Fatalf("adding constant %d failed: %v", i, err)
		}
		if int(index) != i {
			t.Fatalf("constant %d got index %d", i, index)
		}
	}
	if _, err := chunk.AddConstant(value.NewString("overflow")); err == nil {
		t.Fatal("constant pool accepted an index that cannot be represented in bytecode")
	}
}
