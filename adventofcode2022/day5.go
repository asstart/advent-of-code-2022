package adventofcode2022

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Stacks map[int]Stack

type Stack struct {
	// 0 - top box in stack, len(boxes)-1 - last box in stack
	boxes []string
}

type Moves []Move

type Move struct {
	count int
	from  int
	to    int
}

func transpose(data [][]string) [][]string {
	res := make([][]string, len(data[0]))
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i]); j++ {
			res[j] = append(res[j], data[i][j])
		}
	}
	return res
}

func ToStacksAndMoves(ir InputReader) (Stacks, Moves, error) {
	content, err := ir.GetInput()
	if err != nil {
		return nil, nil, err
	}

	idx := 0
	line := "-1"
	horizontalBoxes := [][]string{}
	for {
		line = content[idx]
		finish, _ := regexp.Match("\\s{1,}\\d", []byte(line))
		if finish {
			break
		}
		boxes := []string{}
		ll := len(line)
		for i := 1; i < ll; i += 4 {
			boxes = append(boxes, string(line[i]))
		}
		horizontalBoxes = append(horizontalBoxes, boxes)
		idx++
	}

	vertBoxes := transpose(horizontalBoxes)
	stacks := Stacks{}
	for i := 0; i < len(vertBoxes); i++ {
		boxes := []string{}
		topStack := ""
		for _, b := range vertBoxes[i] {
			if b == " " {
				continue
			}
			if topStack == "" {
				topStack = b
			}
			boxes = append(boxes, b)
		}
		stacks[i] = Stack{boxes: boxes}
	}

	idx += 2 // to jump over empty line for future reading

	foundRe := regexp.MustCompile("move (\\d+) from (\\d+) to (\\d+)")
	moves := Moves{}
	for i := idx; i < len(content); i++ {
		line := content[i]
		found := foundRe.FindStringSubmatch(line)[1:]
		count, _ := strconv.Atoi(found[0])
		from, _ := strconv.Atoi(found[1])
		to, _ := strconv.Atoi(found[2])
		mv := Move{count: count, from: from - 1, to: to - 1}
		moves = append(moves, mv)
	}

	return stacks, moves, nil

}

func reverse(s interface{}) {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

func Task5_1(ir InputReader, cnvrtInpt func(InputReader) (Stacks, Moves, error)) (string, error) {
	stacks, moves, err := cnvrtInpt(ir)
	if err != nil {
		return "", err
	}

	for _, mv := range moves {
		from := stacks[mv.from]
		to := stacks[mv.to]

		tomove := make([]string, mv.count)
		copy(tomove, from.boxes[:mv.count])
		reverse(tomove)
		tomove = append(tomove, to.boxes...)
		to.boxes = tomove
		stacks[mv.to] = to

		if mv.count >= len(from.boxes) {
			from.boxes = []string{}
		} else {
			from.boxes = from.boxes[mv.count:]
		}
		stacks[mv.from] = from
	}

	res := strings.Builder{}
	for i := 0; i < len(stacks); i++ {
		res.WriteString(stacks[i].boxes[0])
	}

	return res.String(), nil
}

func Task5_2(ir InputReader, cnvrtInpt func(InputReader) (Stacks, Moves, error)) (string, error) {
	stacks, moves, err := cnvrtInpt(ir)
	if err != nil {
		return "", err
	}

	for _, mv := range moves {
		from := stacks[mv.from]
		to := stacks[mv.to]

		tomove := make([]string, mv.count)
		copy(tomove, from.boxes[:mv.count])
		tomove = append(tomove, to.boxes...)
		to.boxes = tomove
		stacks[mv.to] = to

		if mv.count >= len(from.boxes) {
			from.boxes = []string{}
		} else {
			from.boxes = from.boxes[mv.count:]
		}
		stacks[mv.from] = from
	}

	res := strings.Builder{}
	for i := 0; i < len(stacks); i++ {
		res.WriteString(stacks[i].boxes[0])
	}

	return res.String(), nil
}
