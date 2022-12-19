package adventofcode2022

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"

	combinations "github.com/mxschmitt/golang-combinations"
)

func idx(s string) int {
	fst := []rune(s)[0]
	scnd := []rune(s)[1]
	return int(26*(fst-'A') + (scnd - 'A'))
}

type Day16Inpt struct {
	AdjacenyM      [][]int
	ValvesPressure map[int]int
}

func ToAdjacencyMatrix(ir InputReader) (Day16Inpt, error) {
	lines, err := ir.GetInput()
	if err != nil {
		return Day16Inpt{}, err
	}

	fromValveRe := regexp.MustCompile("Valve ([A-Z][A-Z]).*")
	toValvesRe := regexp.MustCompile("valves? (([A-Z][A-Z](, )?)+)")
	rateRe := regexp.MustCompile("rate=(\\d+)")
	valves := map[int]int{}
	m := make([][]int, 675)
	for i := 0; i < len(m); i++ {
		m[i] = make([]int, 675)
		for j := 0; j < len(m[i]); j++ {
			if i == j {
				m[i][j] = 0
			} else {
				m[i][j] = math.MaxInt32
			}
		}
	}
	for _, line := range lines {
		from := fromValveRe.FindStringSubmatch(line)
		if len(from) != 2 {
			return Day16Inpt{}, fmt.Errorf("expected format: [Valve AA], got: %v", line)
		}
		to := toValvesRe.FindStringSubmatch(line)
		if len(to) < 2 {
			return Day16Inpt{}, fmt.Errorf("expected format: [valves AA, BB], got: %v", line)
		}
		toList := strings.Split(to[1], ",")
		rateStr := rateRe.FindStringSubmatch(line)
		if len(rateStr) != 2 {
			return Day16Inpt{}, fmt.Errorf("expected format: [rate=5], got: %v", line)
		}
		rate, err := strconv.Atoi(rateStr[1])
		if err != nil {
			return Day16Inpt{}, err
		}
		idxFrom := idx(from[1])
		for _, t := range toList {
			idxTo := idx(strings.TrimSpace(t))
			m[idxFrom][idxTo] = 1
			m[idxTo][idxFrom] = 1
		}
		if rate != 0 {
			valves[idxFrom] = rate
		}
	}
	res := Day16Inpt{
		AdjacenyM:      m,
		ValvesPressure: valves,
	}

	return res, nil
}

// floyd-warshall to generate shortest distances between vertexes
func fw(dist [][]int) {
	for k := 0; k < 675; k++ {
		for i := 0; i < 675; i++ {
			for j := 0; j < 675; j++ {
				if dist[i][k]+dist[k][j] < dist[i][j] {
					dist[i][j] = dist[i][k] + dist[k][j]
				}
			}
		}
	}
}

type SearchState struct {
	MinLeft  int
	Pressure int
}

func conatins(v int, arr []int) bool {
	for _, i := range arr {
		if v == i {
			return true
		}
	}
	return false
}

func search(graph [][]int, curVertex int, path []int, prevKey string, storage map[string]SearchState, valves map[int]int, initMinutes int) {
	for valve, pres := range valves {
		if valve == curVertex || conatins(valve, path) {
			continue
		}
		prev, ok := storage[prevKey]
		if !ok {
			prev = SearchState{
				MinLeft:  initMinutes,
				Pressure: 0,
			}
		}
		
		nPath := make([]int, len(path)+1)
		copy(nPath, path)
		nPath[len(path)] = valve
		nMinLeft := prev.MinLeft - graph[valve][curVertex] - 1
		nPressure := prev.Pressure + nMinLeft*pres

		nKey := buildKey(nPath)

		nV, nOk := storage[nKey]
		if nOk && nV.Pressure > nPressure {
			continue
		}
		storage[nKey] = SearchState{
			MinLeft:  nMinLeft,
			Pressure: nPressure,
		}

		search(graph, valve, nPath, nKey, storage, valves, initMinutes)
	}
}

func buildKey(path []int) string {
	sort.Ints(path)
	str := strings.Builder{}
	for _, i := range path {
		str.WriteString(fmt.Sprintf("%v:", i))
	}
	return str.String()
}

func Task16_1(ir InputReader, cnvrtInpt func(InputReader) (Day16Inpt, error), debug bool) (string, error) {
	input, err := cnvrtInpt(ir)
	if err != nil {
		return "", err
	}

	fw(input.AdjacenyM)

	storage := map[string]SearchState{}
	path := []int{}

	search(input.AdjacenyM, 0, path, "", storage, input.ValvesPressure, 30)

	maxP := 0
	for _, v := range storage {
		if v.Pressure > maxP {
			maxP = v.Pressure
		}
	}

	return fmt.Sprintf("%v", maxP), nil
}

func Task16_2(ir InputReader, cnvrtInpt func(InputReader) (Day16Inpt, error), debug bool) (string, error) {
	input, err := cnvrtInpt(ir)
	if err != nil {
		return "", err
	}

	fw(input.AdjacenyM)

	storage := map[string]SearchState{}
	path := []int{}

	search(input.AdjacenyM, 0, path, "", storage, input.ValvesPressure, 26)

	maxP := 0
	for _, v := range storage {
		if v.Pressure > maxP {
			maxP = v.Pressure
		}
	}

	filteredV := []int{}
	for k, v := range input.ValvesPressure {
		if v != 0 {
			filteredV = append(filteredV, k)
		}
	}

	subsets := combinations.All(filteredV)
	MaxOf2 := 0
	for _, i := range subsets {
		ik := buildKey(i)
		for _, j := range subsets {
			if isDisjoint(i, j) {
				jk := buildKey(j)
				if storage[ik].Pressure+storage[jk].Pressure > MaxOf2 {
					MaxOf2 = storage[ik].Pressure + storage[jk].Pressure
				}
			}
		}
	}
	return fmt.Sprintf("%v", MaxOf2), nil
}

func isDisjoint(arr1 []int, arr2 []int) bool {
	for i := 0; i < len(arr1); i++ {
		for j := 0; j < len(arr2); j++ {
			if arr1[i] == arr2[j] {
				return false
			}
		}
	}
	return true
}
