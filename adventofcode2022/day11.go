package adventofcode2022

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Monkey struct {
	Items     []int
	Operation Operation
	Divisor   int
	IfTrue    int
	IfFalse   int
	Inspected int
}

type OpArgType int

const (
	Old OpArgType = iota
	Custom
)

type OpArg struct {
	Type  OpArgType
	Value int
}

type Operation struct {
	Arg1 OpArg
	Arg2 OpArg
	Func func(arg1 int, arg2 int) int
}

func (op Operation) Execute(value int) int {
	var arg1 int
	var arg2 int

	if op.Arg1.Type == Custom {
		arg1 = op.Arg1.Value
	} else {
		arg1 = value
	}

	if op.Arg2.Type == Custom {
		arg2 = op.Arg2.Value
	} else {
		arg2 = value
	}

	return op.Func(arg1, arg2)

}

func Addition(arg1 int, arg2 int) int {
	return arg1 + arg2
}

func Multiplaction(arg1 int, arg2 int) int {
	return arg1 * arg2
}

type Monkeys map[int]Monkey

func ToMonkeys(ir InputReader) (Monkeys, error) {
	lines, err := ir.GetInput()
	if err != nil {
		return nil, err
	}

	res := map[int]Monkey{}

	mnkIdxRe := regexp.MustCompile("Monkey (\\d):")
	trueCond := regexp.MustCompile("If true: throw to monkey (\\d+)$")
	falseCond := regexp.MustCompile("If false: throw to monkey (\\d+)")
	for i := 0; i < len(lines); i += 7 {
		if lines[i] == "" {
			continue
		}
		// parsing monkey idx
		idxs := mnkIdxRe.FindStringSubmatch(strings.TrimSpace(lines[i]))
		if len(idxs) != 2 {
			return nil, fmt.Errorf("expected format: [Monkey <N>:], got: [%v]", lines[i])
		}
		idx, err := strconv.Atoi(idxs[1])
		if err != nil {
			return nil, err
		}

		// parsing monkey items
		items := []int{}
		itemsArr := strings.Split(strings.ReplaceAll(lines[i+1], " ", ""), ":")
		if len(itemsArr) > 1 {
			itemsStr := strings.Split(itemsArr[1], ",")
			for _, i := range itemsStr {
				parsed, err := strconv.Atoi(i)
				if err != nil {
					return nil, err
				}
				items = append(items, parsed)
			}
		}

		// parsing operation

		opItms := strings.Split(lines[i+2], "=")
		if len(opItms) == 0 || len(opItms) > 2 {
			return nil, fmt.Errorf("expected format: [Operation: new = old * 19], got: [%v]", lines[i+2])
		}
		opItm := strings.TrimSpace(opItms[1])
		opItmSplitted := strings.Split(opItm, " ")
		if len(opItmSplitted) != 3 {
			return nil, fmt.Errorf("expected format: [old * 19], got: [%v]", opItm)
		}
		op := Operation{}
		switch opItmSplitted[0] {
		case "old":
			op.Arg1 = OpArg{Type: Old}
		default:
			i, err := strconv.Atoi(opItmSplitted[0])
			if err != nil {
				return nil, fmt.Errorf("expected format: [old * 19], got: [%v]", opItmSplitted[0])
			}
			op.Arg1 = OpArg{Type: Custom, Value: i}
		}
		switch opItmSplitted[1] {
		case "+":
			op.Func = Addition
		case "*":
			op.Func = Multiplaction
		default:
			return nil, fmt.Errorf("unsupported op: %v", opItmSplitted[1])
		}
		switch opItmSplitted[2] {
		case "old":
			op.Arg2 = OpArg{Type: Old}
		default:
			i, err := strconv.Atoi(opItmSplitted[2])
			if err != nil {
				return nil, fmt.Errorf("expected format: [old * 19], got: [%v]", opItmSplitted[2])
			}
			op.Arg2 = OpArg{Type: Custom, Value: i}
		}

		// parse condition

		condItms := strings.Split(strings.TrimSpace(lines[i+3]), " ")
		denominator, err := strconv.Atoi(condItms[len(condItms)-1])
		if err != nil {
			return nil, fmt.Errorf("expected condition line, got: %v", lines[i+3])
		}

		// if true

		ifTrueFound := trueCond.FindStringSubmatch(strings.TrimSpace(lines[i+4]))
		if len(ifTrueFound) != 2 {
			return nil, fmt.Errorf("expected: [If true: throw to monkey 1], got: %v", lines[i+4])
		}
		ifTrue, err := strconv.Atoi(ifTrueFound[1])
		if err != nil {
			return nil, err
		}

		// if false

		ifFalseFound := falseCond.FindStringSubmatch(strings.TrimSpace(lines[i+5]))
		if len(ifTrueFound) != 2 {
			return nil, fmt.Errorf("expected: [If false: throw to monkey 1], got: %v", lines[i+5])
		}
		ifFalse, err := strconv.Atoi(ifFalseFound[1])
		if err != nil {
			return nil, err
		}

		mnk := Monkey{
			Items:     items,
			Operation: op,
			Divisor:   denominator,
			IfTrue:    ifTrue,
			IfFalse:   ifFalse,
		}

		res[idx] = mnk

	}

	return res, nil
}

func getLcm(m Monkeys) int {
	lcm := 1
	for _, v := range m {
		lcm *= v.Divisor
	}
	return lcm
}

func inspect(monkeys Monkeys, lcm int, relief func(int) int) error {
	for i := 0; i < len(monkeys); i++ {
		curr := monkeys[i]
		for len(curr.Items) > 0 {
			curr.Inspected += 1
			tmpWorry := curr.Operation.Execute(curr.Items[0])
			tmpWorry %= lcm
			tmpWorry = relief(tmpWorry)
			if tmpWorry%curr.Divisor == 0 {
				throwTo, ok := monkeys[curr.IfTrue]
				if !ok {
					return fmt.Errorf("monkey: [%v] not exists", curr.IfTrue)
				}
				throwTo.Items = append(throwTo.Items, tmpWorry)
				monkeys[curr.IfTrue] = throwTo
			} else {
				throwTo, ok := monkeys[curr.IfFalse]
				if !ok {
					return fmt.Errorf("monkey: [%v] not exists", curr.IfFalse)
				}
				throwTo.Items = append(throwTo.Items, tmpWorry)
				monkeys[curr.IfFalse] = throwTo
			}
			curr.Items = curr.Items[1:]
		}
		monkeys[i] = curr
	}
	return nil
}

func Task11_1(ir InputReader, cnvrtInpt func(InputReader) (Monkeys, error), debug bool) (string, error) {
	monkeys, err := cnvrtInpt(ir)
	if err != nil {
		return "", err
	}

	lcm := getLcm(monkeys)
	rounds := 20
	for r := 0; r < rounds; r++ {
		err := inspect(monkeys, lcm, func(v int) int { return v / 3 })
		if err != nil {
			return "", err
		}
	}

	max1 := 0
	max2 := 0
	for _, v := range monkeys {
		if v.Inspected > max1 {
			if v.Inspected > max2 {
				tmp := max2
				max2 = v.Inspected
				max1 = tmp
			} else {
				max1 = v.Inspected
			}
		}
	}
	return fmt.Sprintf("Result: %v", max1*max2), nil
}

func Task11_2(ir InputReader, cnvrtInpt func(InputReader) (Monkeys, error), debug bool) (string, error) {
	monkeys, err := cnvrtInpt(ir)
	if err != nil {
		return "", err
	}

	lcm := getLcm(monkeys)
	rounds := 10000
	for r := 0; r < rounds; r++ {
		err := inspect(monkeys, lcm, func(v int) int { return v })
		if err != nil {
			return "", err
		}
	}

	max1 := int64(0)
	max2 := int64(0)
	for _, v := range monkeys {
		if int64(v.Inspected) > max1 {
			if int64(v.Inspected) > max2 {
				tmp := max2
				max2 = int64(v.Inspected)
				max1 = tmp
			} else {
				max1 = int64(v.Inspected)
			}
		}
	}
	return fmt.Sprintf("Result: %v", max1*max2), nil
}
