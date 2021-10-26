package inspector

import (
	"fmt"
	"testing"
)

func TestByteSize(t *testing.T) {
	d := NewByteSize(`1000`, `KB`)
	if d.value != 1024000 {
		t.Error("Did not set byte value correctly")
	}
	second := NewByteSize(`0.9765625`, `MB`)
	if second.value != d.value {
		fmt.Println(second.value, d.value)
		t.Error("MB value not equivalent to KB value")
	}
	if second.format(`KB`) != 1000 {
		t.Error("Could not convert MB back to KB value")
	}
}
