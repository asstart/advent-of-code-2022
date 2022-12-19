package adventofcode2022

import "fmt"



func Task6_1(ir InputReader, cnvrtInpt func(InputReader) (string, error)) (string, error) {
	data, err := cnvrtInpt(ir)
	if err != nil {
		return "", err
	}

	idx := getIdxFirstUniqueSubsEnds(data, 4)

	if idx == -1 {
		return "not found", nil
	} else {
		return fmt.Sprintf("Result: %v", idx), nil
	}
}

func Task6_2(ir InputReader, cnvrtInpt func(InputReader) (string, error)) (string, error) {
	data, err := cnvrtInpt(ir)
	if err != nil {
		return "", err
	}

	idx := getIdxFirstUniqueSubsEnds(data, 14)

	if idx == -1 {
		return "not found", nil
	} else {
		return fmt.Sprintf("Result: %v", idx), nil
	}
}

func getIdxFirstUniqueSubsEnds(str string, substrlen int) int {
	substrendsIdx := -1
	for i := 0; i < len(str)-substrlen; i++ {
		checker := 0
		found := true
		for j := i; j < i+substrlen; j++ {
			bitIdx := str[j] - 'a'
			if (checker & (1 << bitIdx)) > 0 {
				found = false
				break
			}
			checker |= (1 << bitIdx)
		}
		if found {
			substrendsIdx = i + substrlen
			break
		}
	}
	return substrendsIdx
}
