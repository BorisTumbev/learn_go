package iteration

import (
	"fmt"
	"strings"
	"testing"
)

const testRepeat = 5

func TestRepeat(t *testing.T) {
	repeated := Repeat("a", testRepeat)
	expected := strings.Repeat("a", testRepeat)

	if repeated != expected {
		t.Errorf("expected %q but got %q", expected, repeated)
	}
}

func BenchmarkRepeat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Repeat("a", testRepeat)
	}
}

func ExampleRepeat() {
	rep := Repeat("a", 3)
	fmt.Print(rep)
	// Output: aaa
}
