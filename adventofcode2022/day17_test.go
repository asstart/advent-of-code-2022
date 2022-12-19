package adventofcode2022_test

import (
	"testing"

	"github.com/asstart/advent-of-code-2022/adventofcode2022"
)

func BenchmarkTask17_2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		adventofcode2022.Task17_1(
			&adventofcode2022.FileToStringsInputReader{Path: "../adventofcode2022/day17.data"},
			adventofcode2022.ToDirections,
			false,
		)
	}
}
