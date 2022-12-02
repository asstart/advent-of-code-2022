package adventofcode2022

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func Task1_1(ir InputReader, convertInput func(ir InputReader) []string) (string, error) {

	lines := convertInput(ir)

	max := 0
	tmpSum := 0

	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			parsed, err := strconv.Atoi(l)
			if err != nil {
				return "", err
			}
			tmpSum += parsed
		} else {
			if tmpSum > max {
				max = tmpSum
			}
			tmpSum = 0
		}
	}

	// to handle the last tmpSum
	if tmpSum > max {
		max = tmpSum
	}

	return fmt.Sprintf("Day 1 Part 1 result: %v", max), nil
}

func Task1_2(ir InputReader, convertInput func(ir InputReader) []string) (string, error) {
	lines := convertInput(ir)

	max := [3]int{}
	tmpSum := 0

	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			parsed, err := strconv.Atoi(l)
			if err != nil {
				return "", err
			}
			tmpSum += parsed
		} else {
			minIdx, min := getMin(max[:])
			if tmpSum > min {
				max[minIdx] = tmpSum
			}
			tmpSum = 0
		}
	}

	// to handle the last tmpSum
	minIdx, min := getMin(max[:])
	if tmpSum > min {
		max[minIdx] = tmpSum
	}
	tmpSum = 0

	total := 0
	for _, i := range max {
		total += i
	}

	return fmt.Sprintf("Day 1 Part 2 items: %v, result: %v", max, total), nil
}

func getMin(arr []int) (int, int) {
	min := math.MaxInt
	i := -1
	for idx, itm := range arr {
		if itm < min {
			min = itm
			i = idx
		}
	}
	return i, min
}
