package adventofcode2022

import (
	"fmt"
	"math"
)

type TupleIntArr struct {
	l []int
	r []int
}

func ToCharArr(ir InputReader) []TupleIntArr {
	content, err := ir.GetInput()
	if err != nil {
		panic(err)
	}

	lines, ok := content.([]string)
	if !ok {
		panic(fmt.Errorf("can't convert input to lines arr"))
	}

	converted := []TupleIntArr{}
	for _, line := range lines {
		// according to task condition, it should be even number of symbols in line
		middleIdx := len(line) / 2
		fst := line[:middleIdx]
		scnd := line[middleIdx:]
		converted = append(converted, TupleIntArr{
			l: stringToPriorityArr(fst),
			r: stringToPriorityArr(scnd),
		})
	}

	return converted
}

func stringToPriorityArr(s string) []int {
	res := []int{}
	for _, r := range s {
		res = append(res, getPriority(r))
	}
	return res
}

func getPriority(r rune) int {
	if r >= 97 {
		return int(math.Abs(float64(97-int(r)))) + 1
	} else {
		return int(math.Abs(float64(65-r))) + 1 + 26
	}
}

//Solution in this task is only for group of two sequences of items
//General solution can be found in Task3_2
func Task3_1(ir InputReader, convertInput func(ir InputReader) []TupleIntArr) (string, error) {
	data := convertInput(ir)
	commons := []int{}
	for _, t := range data {
		lefts := [53]int{}
		rights := [53]int{}
		found := false
		//populate temp storage, in case we found common element finish eagerly
		for i := 0; i < len(t.l); i++ {
			if rights[t.l[i]] == 1 {
				commons = append(commons, t.l[i])
				found = true
				break
			} else {
				lefts[t.l[i]] = 1
			}

			if lefts[t.r[i]] == 1 {
				commons = append(commons, t.r[i])
				found = true
				break
			} else {
				rights[t.r[i]] = 1
			}
		}

		if found {
			continue
		}
		//if has't finished eagerly, one more iteration
		for i := 0; i < len(t.l); i++ {
			if rights[t.l[i]] == 1 {
				commons = append(commons, t.l[i])
				break
			}
			if lefts[t.r[i]] == 1 {
				commons = append(commons, t.r[i])
				break
			}
		}
	}

	sum := 0
	for _, i := range commons {
		sum += i
	}

	return fmt.Sprintf("result: %v", sum), nil
}

func ToArrOfIntArr(ir InputReader, groups int) [][][]int {
	content, err := ir.GetInput()
	if err != nil {
		panic(err)
	}

	lines, ok := content.([]string)
	if !ok {
		panic(fmt.Errorf("can't convert input to lines arr"))
	}

	if len(lines)%groups != 0 {
		panic(fmt.Errorf("bad input, %v is not dividable by %v", len(lines), groups))
	}

	converted := [][][]int{}
	groupN := 0
	for i := 0; i < len(lines); i += groups {
		converted = append(converted, [][]int{})
		for j := 0; j < groups; j++ {
			converted[groupN] = append(converted[groupN], stringToPriorityArr(lines[i+j]))
		}
		groupN++
	}

	return converted
}

func storageIdx(bpNum int, item int) int {
	return bpNum*52 + item - 1
}

// data structure is follwing:
//
// [][][]int is:
// groups[
//	backpups_in_group[
//		items_in_backpack[
//			item_value
//			]
//		]
//	]
// ]
func Task3_2(ir InputReader, convertInput func(ir InputReader, groups int) [][][]int) (string, error) {
	groupSize := 3
	data := convertInput(ir, 3)

	commons := []int{}
	for _, group := range data {
		storage := map[int]bool{}
		for i, bp := range group {
			for _, item := range bp {
				storage[storageIdx(i, item)] = true
			}
		}

		//It should be enough to check only first backpack in group,
		//because it will contain common item for all groups
		check_bp := group[0]
		for _, item := range check_bp {
			found := true
			for i := 0; i < groupSize; i++ {
				if !found {
					break
				}
				_, ok := storage[storageIdx(i, item)]
				if !ok {
					found = false
					break
				}
			}
			if found {
				commons = append(commons, item)
				break
			}
		}
	}

	sum := 0
	for _, i := range commons {
		sum += i
	}
	return fmt.Sprintf("result: %v", sum), nil
}
