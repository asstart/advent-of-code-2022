package adventofcode2022

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type TupleString struct {
	_1 string
	_2 string
}

func ToArrTupleString(ir InputReader) ([]TupleString, error) {
	lines, err := ir.GetInput()
	if err != nil {
		return nil, err
	}

	res := []TupleString{}
	for i := 0; i < len(lines); i += 3 {
		res = append(res, TupleString{
			_1: lines[i],
			_2: lines[i+1],
		})
	}

	return res, nil
}

func isWrapped(line string) bool {
	return strings.HasPrefix(line, "[")
}

func unwrapArr(line string) string {
	return strings.TrimSuffix(strings.TrimPrefix(line, "["), "]")
}

func splitTopLevel(line string) []string {
	tokens := []string{}
	opened := 0
	bld := strings.Builder{}
	for _, item := range line {
		if string(item) == "[" {
			bld.WriteString(string(item))
			opened++
			continue
		} else if string(item) == "]" {
			bld.WriteString(string(item))
			opened--
			if opened == 0 {
				tokens = append(tokens, bld.String())
				bld = strings.Builder{}
			}
			continue
		}

		if opened > 0 {
			bld.WriteString(string(item))
		} else if opened == 0 {
			if string(item) == "," || string(item) == " " {
				if bld.Len() > 0 {
					tokens = append(tokens, bld.String())
					bld = strings.Builder{}
				}
			} else {
				bld.WriteString(string(item))
			}
		}
	}
	if bld.Len() > 0 {
		tokens = append(tokens, bld.String())
	}
	return tokens
}

func compareInt(l int, r int) int {
	if l > r {
		return 1
	} else if l == r {
		return 0
	} else {
		return -1
	}
}

func comparePacket(l string, r string) int {
	var inl bool = true
	var inr bool = true
	nl, err := strconv.Atoi(l)
	if err != nil {
		inl = false
	}
	nr, err := strconv.Atoi(r)
	if err != nil {
		inr = false
	}

	if inl && inr {
		return compareInt(nl, nr)
	}

	iwl := isWrapped(l)
	iwr := isWrapped(r)
	var tokens1 []string
	var tokens2 []string

	if iwl {
		tokens1 = splitTopLevel(unwrapArr(l))
	} else {
		tokens1 = splitTopLevel(l)
	}

	if iwr {
		tokens2 = splitTopLevel(unwrapArr(r))
	} else {
		tokens2 = splitTopLevel(r)
	}

	var cr int

	maxLen := Max(len(tokens1), len(tokens2))
	for i := 0; i < maxLen; i++ {
		if len(tokens1) == i {
			return -1
		} else if len(tokens2) == i {
			return 1
		}

		cr = comparePacket(tokens1[i], tokens2[i])
		if cr != 0 {
			return cr
		}
	}
	return cr
}

func Task13_1(ir InputReader, cnvrtInpt func(InputReader) ([]TupleString, error), debug bool) (string, error) {
	tuples, err := cnvrtInpt(ir)
	if err != nil {
		return "", err
	}

	res := make([]bool, len(tuples))
	for idx, t := range tuples {
		cr := comparePacket(t._1, t._2)
		res[idx] = cr < 0

	}

	sum := 0
	for idx, i := range res {
		if i {
			sum += idx + 1
		}
	}

	if debug {
		d13p1Debug(res)
	}

	return fmt.Sprintf("Result: %v", sum), nil
}

func Task13_2(ir InputReader, cnvrtInpt func(InputReader) ([]TupleString, error), debug bool) (string, error) {
	tuples, err := cnvrtInpt(ir)
	if err != nil {
		return "", err
	}

	tuples = append(tuples, TupleString{
		_1: "[[2]]",
		_2: "[[6]]",
	})

	packets := []string{}
	for _, t := range tuples {
		packets = append(packets, t._1)
		packets = append(packets, t._2)
	}

	sort.SliceStable(packets, func(i int, j int) bool {
		cr := comparePacket(packets[i], packets[j])
		return cr < 0
	})

	idx2, idx6 := 0, 0
	for idx, p := range packets {
		if p == "[[2]]" {
			idx2 = idx + 1
		} else if p == "[[6]]" {
			idx6 = idx + 1
		}
	}

	if debug {
		d13p2Debug(packets)
	}

	return fmt.Sprintf("Result: %v", idx2*idx6), nil
}

func d13p1Debug(comparison []bool) {
	for idx, v := range comparison {
		fmt.Printf("Idx: %v, is right order: %v\n", idx+1, v)
	}
}
func d13p2Debug(packates []string) {
	for _, p := range packates {
		fmt.Printf("%v\n", p)
	}
}
