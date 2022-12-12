package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/asstart/advent-of-code-2022/adventofcode2022"
	"github.com/jessevdk/go-flags"
)

type opts struct {
	N string `short:"n"  description:"Number of task in format day_part, like 1_1, 1_2"`
	A bool   `short:"a"  descriotion:"Run all tasks"`
	D bool   `short:"d" description:"Debug mode"`
}

func main() {
	var o opts
	if _, err := flags.Parse(&o); err != nil {
		fmt.Printf("error while parsing flags: %s\n", err)
		os.Exit(1)
	}

	if o.N != "" && o.A {
		fmt.Printf("options (a, n) mustn't be used simultaneously, choose one!\n")
		os.Exit(1)
	}

	if o.N == "" && !o.A {
		fmt.Printf("at least one option (a, n) must be specified\n")
		os.Exit(1)
	}

	if o.A {
		runAll(o)
		os.Exit(0)
	}

	if o.N != "" {
		runTask(o)
		os.Exit(0)
	}
}

type RunFunc func(o opts) string

var tasks = map[string]RunFunc{
	"1_1":  t1_1,
	"1_2":  t1_2,
	"2_1":  t2_1,
	"2_2":  t2_2,
	"3_1":  t3_1,
	"3_2":  t3_2,
	"4_1":  t4_1,
	"4_2":  t4_2,
	"5_1":  t5_1,
	"5_2":  t5_2,
	"6_1":  t6_1,
	"6_2":  t6_2,
	"7_1":  t7_1,
	"7_2":  t7_2,
	"8_1":  t8_1,
	"8_2":  t8_2,
	"9_1":  t9_1,
	"9_2":  t9_2,
	"10_1": t10_1,
	"10_2": t10_2,
	"11_1": t11_1,
	"11_2": t11_2,
	"12_1": t12_1,
	"12_2": t12_2,
}

func runAll(o opts) {
	keys := make([]string, 0, len(tasks))
	for k := range tasks {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		f := tasks[k]
		r := f(o)
		fmt.Printf("Running task: %v\nResult      : %v\n", k, r)
	}
}

func runTask(o opts) {
	f, ok := tasks[o.N]
	if !ok {
		fmt.Printf("Task: %v not found\n", o.N)
		os.Exit(1)
	}
	r := f(o)
	fmt.Printf("Running task: %v\nResult      : %v\n", o.N, r)
}

func t1_1(o opts) string {
	res, err := adventofcode2022.Task1_1(
		&adventofcode2022.FileToStringsInputReader{Path: "adventofcode2022/day1.data"},
		adventofcode2022.ToIntOrSpaceArr,
	)
	if err != nil {
		res = err.Error()
	}
	return res
}

func t1_2(o opts) string {
	res, err := adventofcode2022.Task1_2(
		&adventofcode2022.FileToStringsInputReader{Path: "adventofcode2022/day1.data"},
		adventofcode2022.ToIntOrSpaceArr,
	)
	if err != nil {
		return err.Error()
	}
	return res
}

func t2_1(o opts) string {
	res, err := adventofcode2022.Task2_1(
		&adventofcode2022.FileToStringsInputReader{Path: "adventofcode2022/day2.data"},
		adventofcode2022.ToTupleRPSArr,
	)
	if err != nil {
		return err.Error()
	}
	return res
}

func t2_2(o opts) string {
	res, err := adventofcode2022.Task2_2(
		&adventofcode2022.FileToStringsInputReader{Path: "adventofcode2022/day2.data"},
		adventofcode2022.ToTupleRPSArr,
	)
	if err != nil {
		return err.Error()
	}
	return res
}

func t3_1(o opts) string {
	res, err := adventofcode2022.Task3_1(
		&adventofcode2022.FileToStringsInputReader{Path: "adventofcode2022/day3.data"},
		adventofcode2022.ToTupleIntArr,
	)
	if err != nil {
		return err.Error()
	}
	return res
}

func t3_2(o opts) string {
	res, err := adventofcode2022.Task3_2(
		&adventofcode2022.FileToStringsInputReader{Path: "adventofcode2022/day3.data"},
		adventofcode2022.To3DArray,
	)
	if err != nil {
		return err.Error()
	}
	return res
}

func t4_1(o opts) string {
	res, err := adventofcode2022.Task4_1(
		&adventofcode2022.FileToStringsInputReader{Path: "adventofcode2022/day4.data"},
		adventofcode2022.ToTupleSegment,
	)
	if err != nil {
		return err.Error()
	}
	return res
}

func t4_2(o opts) string {
	res, err := adventofcode2022.Task4_2(
		&adventofcode2022.FileToStringsInputReader{Path: "adventofcode2022/day4.data"},
		adventofcode2022.ToTupleSegment,
	)
	if err != nil {
		return err.Error()
	}
	return res
}

func t5_1(o opts) string {
	res, err := adventofcode2022.Task5_1(
		&adventofcode2022.FileToStringsInputReader{Path: "adventofcode2022/day5.data"},
		adventofcode2022.ToStacksAndMoves,
	)
	if err != nil {
		return err.Error()
	}
	return res
}

func t5_2(o opts) string {
	res, err := adventofcode2022.Task5_2(
		&adventofcode2022.FileToStringsInputReader{Path: "adventofcode2022/day5.data"},
		adventofcode2022.ToStacksAndMoves,
	)
	if err != nil {
		return err.Error()
	}
	return res
}

func t6_1(o opts) string {
	res, err := adventofcode2022.Task6_1(
		&adventofcode2022.FileToStringsInputReader{Path: "adventofcode2022/day6.data"},
		adventofcode2022.ToSingleLine,
	)
	if err != nil {
		return err.Error()
	}
	return res
}

func t6_2(o opts) string {
	res, err := adventofcode2022.Task6_2(
		&adventofcode2022.FileToStringsInputReader{Path: "adventofcode2022/day6.data"},
		adventofcode2022.ToSingleLine,
	)
	if err != nil {
		return err.Error()
	}
	return res
}

func t7_1(o opts) string {
	res, err := adventofcode2022.Task7_1(
		&adventofcode2022.FileToStringsInputReader{Path: "adventofcode2022/day7.data"},
		adventofcode2022.ToCmdQueue,
	)
	if err != nil {
		return err.Error()
	}
	return res
}

func t7_2(o opts) string {
	res, err := adventofcode2022.Task7_2(
		&adventofcode2022.FileToStringsInputReader{Path: "adventofcode2022/day7.data"},
		adventofcode2022.ToCmdQueue,
	)
	if err != nil {
		return err.Error()
	}
	return res
}

func t8_1(o opts) string {
	res, err := adventofcode2022.Task8_1(
		&adventofcode2022.FileToStringsInputReader{Path: "adventofcode2022/day8.data"},
		adventofcode2022.To2DTreeInfoArray,
		o.D,
	)
	if err != nil {
		return err.Error()
	}
	return res
}

func t8_2(o opts) string {
	res, err := adventofcode2022.Task8_2(
		&adventofcode2022.FileToStringsInputReader{Path: "adventofcode2022/day8.data"},
		adventofcode2022.To2DTreeInfoArray,
		o.D,
	)
	if err != nil {
		return err.Error()
	}
	return res
}

func t9_1(o opts) string {
	res, err := adventofcode2022.Task9_1(
		&adventofcode2022.FileToStringsInputReader{Path: "adventofcode2022/day9.data"},
		adventofcode2022.ToMoves,
		o.D,
	)
	if err != nil {
		return err.Error()
	}
	return res
}

func t9_2(o opts) string {
	res, err := adventofcode2022.Task9_2(
		&adventofcode2022.FileToStringsInputReader{Path: "adventofcode2022/day9.data"},
		adventofcode2022.ToMoves,
		o.D,
	)
	if err != nil {
		return err.Error()
	}
	return res
}

func t10_1(o opts) string {
	res, err := adventofcode2022.Task10_1(
		&adventofcode2022.FileToStringsInputReader{Path: "adventofcode2022/day10.data"},
		adventofcode2022.ToStatefulCmds,
		o.D,
	)
	if err != nil {
		return err.Error()
	}
	return res
}

func t10_2(o opts) string {
	res, err := adventofcode2022.Task10_2(
		&adventofcode2022.FileToStringsInputReader{Path: "adventofcode2022/day10.data"},
		adventofcode2022.ToStatefulCmds,
		o.D,
	)
	if err != nil {
		return err.Error()
	}
	return res
}

func t11_1(o opts) string {
	res, err := adventofcode2022.Task11_1(
		&adventofcode2022.FileToStringsInputReader{Path: "adventofcode2022/day11.data"},
		adventofcode2022.ToMonkeys,
		o.D,
	)
	if err != nil {
		return err.Error()
	}
	return res
}

func t11_2(o opts) string {
	res, err := adventofcode2022.Task11_2(
		&adventofcode2022.FileToStringsInputReader{Path: "adventofcode2022/day11.data"},
		adventofcode2022.ToMonkeys,
		o.D,
	)
	if err != nil {
		return err.Error()
	}
	return res
}

func t12_1(o opts) string {
	res, err := adventofcode2022.Task12_1(
		&adventofcode2022.FileToStringsInputReader{Path: "adventofcode2022/day12.data"},
		adventofcode2022.ToElevationMap,
		o.D,
	)
	if err != nil {
		return err.Error()
	}
	return res
}

func t12_2(o opts) string {
	res, err := adventofcode2022.Task12_2(
		&adventofcode2022.FileToStringsInputReader{Path: "adventofcode2022/day12.data"},
		adventofcode2022.ToElevationMap,
		o.D,
	)
	if err != nil {
		return err.Error()
	}
	return res
}
