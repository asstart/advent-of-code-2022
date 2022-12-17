package adventofcode2022

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeInts(t *testing.T) {
	tt := []struct {
		name     string
		line     Line
		lines    []Line
		expected []Line
	}{
		// {"single, full before, no overlap", Line{1, 2}, []Line{Line{3, 4}}, []Line{{1, 2}, {3, 4}}},
		// {"single, from before", Line{1, 3}, []Line{Line{2, 4}}, []Line{{1, 4}}},
		{"single, from before same", Line{1, 2}, []Line{Line{2, 4}}, []Line{{1, 4}}},
		{"single, full inside", Line{2, 3}, []Line{Line{1, 4}}, []Line{{1, 4}}},
		{"single, full inside same", Line{1, 4}, []Line{Line{1, 4}}, []Line{{1, 4}}},
		{"single, from inside", Line{2, 4}, []Line{Line{1, 3}}, []Line{{1, 4}}},
		{"single, from inside same from", Line{1, 4}, []Line{Line{1, 3}}, []Line{{1, 4}}},
		{"single, from inside same to", Line{3, 4}, []Line{Line{1, 3}}, []Line{{1, 4}}},
		{"single, to after", Line{3, 4}, []Line{Line{1, 2}}, []Line{{1, 2}, {3, 4}}},
		{"single, to after same", Line{2, 3}, []Line{Line{1, 2}}, []Line{{1, 3}}},
		{"single, middle on overlap", Line{3, 4}, []Line{Line{1, 2}, {6, 7}}, []Line{{1, 2}, {3, 4}, {6, 7}}},
		{"single, middle on overlap border", Line{3, 4}, []Line{Line{1, 2}, {5, 6}}, []Line{{1, 2}, {3, 4}, {5, 6}}},
		{"single, middle on overlap same", Line{2, 5}, []Line{Line{1, 2}, {5, 7}}, []Line{{1, 7}}},

		{"mult, full before, no ovelap", Line{1, 2}, []Line{{3, 4}, {5, 6}}, []Line{{1, 2}, {3, 4}, {5, 6}}},
		{"mult, from before first", Line{1, 3}, []Line{{2, 4}, {5, 6}}, []Line{{1, 4}, {5, 6}}},
		{"mult, from before middle", Line{3, 5}, []Line{{1, 2}, {4, 6}, {7, 8}}, []Line{{1, 2}, {3, 6}, {7, 8}}},
		{"mult, from before last", Line{3, 5}, []Line{{1, 2}, {4, 6}}, []Line{{1, 2}, {3, 6}}},

		{"mult, inside first", Line{2, 3}, []Line{{1, 4}, {5, 6}}, []Line{{1, 4}, {5, 6}}},
		{"mult, inside middle", Line{4, 8}, []Line{{1, 2}, {4, 8}, {9, 10}}, []Line{{1, 2}, {4, 8}, {9, 10}}},
		{"mult, inside last", Line{6, 8}, []Line{{1, 2}, {3, 4}, {5, 10}}, []Line{{1, 2}, {3, 4}, {5, 10}}},

		{"mult, from before, ov fst, finish outside", Line{1, 5}, []Line{{2, 3}, {6, 7}}, []Line{{1, 5}, {6, 7}}},
		{"mult, from before, overlap full, finish outside", Line{1, 7}, []Line{{2, 3}, {4, 5}}, []Line{{1, 7}}},
		{"mult, from before, ov except last, finish outside", Line{1, 7}, []Line{{2, 3}, {4, 5}, {9, 10}}, []Line{{1, 7}, {9, 10}}},
		{"mult, from before, ov last, finish outside", Line{4, 7}, []Line{{1, 2}, {5, 6}}, []Line{{1, 2}, {4, 7}}},

		{"mult, from inside, ov fst, finish outside", Line{2, 5}, []Line{{1, 4}, {6, 7}}, []Line{{1, 5}, {6, 7}}},
		{"mult, from inside, overlap full, finish outside", Line{2, 7}, []Line{{1, 3}, {4, 5}}, []Line{{1, 7}}},
		{"mult, from inside, ov except last, finish outside", Line{2, 7}, []Line{{1, 3}, {4, 5}, {9, 10}}, []Line{{1, 7}, {9, 10}}},
		{"mult, from inside, ov last, finish outside", Line{5, 7}, []Line{{1, 2}, {4, 6}}, []Line{{1, 2}, {4, 7}}},

		{"mult, from before, ov fst, finish inside", Line{1, 4}, []Line{{2, 5}, {6, 7}}, []Line{{1, 5}, {6, 7}}},
		{"mult, from before, overlap full, finish inside", Line{1, 6}, []Line{{2, 3}, {4, 7}}, []Line{{1, 7}}},
		{"mult, from before, ov except last, finish inside", Line{1, 6}, []Line{{2, 3}, {4, 7}, {9, 10}}, []Line{{1, 7}, {9, 10}}},
		{"mult, from before, ov last, finish inside", Line{4, 6}, []Line{{1, 2}, {5, 7}}, []Line{{1, 2}, {4, 7}}},

		{"mult, full before, no ovelap", Line{14, 14}, []Line{{-1, 5}, {8, 12}}, []Line{{-1, 5}, {8, 12}, {14, 14}}},
		// {"mult, full before, no ovelap", Line{1,2}, []Line{{3,4},{5,6}}, []Line{{1,2},{3,4},{5,6}}},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			lines := mergeSegment(tc.line, tc.lines)
			assert.ElementsMatch(t, tc.expected, lines)
		})
	}
}
