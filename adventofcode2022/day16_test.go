package adventofcode2022_test

import (
	"testing"

	"github.com/asstart/advent-of-code-2022/adventofcode2022"
	"github.com/stretchr/testify/assert"
)

func BenchmarkTask16_1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		adventofcode2022.Task16_1(
			&adventofcode2022.FileToStringsInputReader{Path: "../adventofcode2022/day16.data"},
			adventofcode2022.ToAdjacencyMatrix,
			false,
		)
	}
}

func TestTask16_1(t *testing.T) {
	res, err := adventofcode2022.Task16_1(
		&adventofcode2022.FileToStringsInputReader{Path: "../adventofcode2022/day16.data"},
		adventofcode2022.ToAdjacencyMatrix,
		false,
	)
	assert.Nil(t, err)
	assert.Equal(t, "1376", res)
}

func BenchmarkTask16_2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		adventofcode2022.Task16_2(
			&adventofcode2022.FileToStringsInputReader{Path: "../adventofcode2022/day16.data"},
			adventofcode2022.ToAdjacencyMatrix,
			false,
		)
	}
}

func TestTask16_2(t *testing.T) {
	res, err := adventofcode2022.Task16_2(
		&adventofcode2022.FileToStringsInputReader{Path: "../adventofcode2022/day16.data"},
		adventofcode2022.ToAdjacencyMatrix,
		false,
	)
	assert.Nil(t, err)
	assert.Equal(t, "1933", res)
}