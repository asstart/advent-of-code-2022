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
		runAll()
		os.Exit(0)
	}

	if o.N != "" {
		runTask(o)
		os.Exit(0)
	}
}

type RunFunc func() string

var tasks = map[string]RunFunc{
	"1_1": t1_1,
	"1_2": t1_2,
	"2_1": t2_1,
	"2_2": t2_2,
}

func runAll() {
	keys := make([]string, 0, len(tasks))
	for k := range tasks {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		f := tasks[k]
		r := f()
		fmt.Printf("Running task: %v\nResult      : %v\n", k, r)
	}
}

func runTask(o opts) {
	f, ok := tasks[o.N]
	if !ok {
		fmt.Printf("Task: %v not found\n")
		os.Exit(1)
	}
	r := f()
	fmt.Printf("Running task: %v\nResult      : %v\n", o.N, r)
}

func t1_1() string {
	res, err := adventofcode2022.Task1_1(
		&adventofcode2022.FileToStringsInputReader{"adventofcode2022/day1.data"},
		adventofcode2022.ToLines,
	)
	if err != nil {
		res = err.Error()
	}
	return res
}

func t1_2() string {
	res, err := adventofcode2022.Task1_2(
		&adventofcode2022.FileToStringsInputReader{"adventofcode2022/day1.data"},
		adventofcode2022.ToLines,
	)
	if err != nil {
		return err.Error()
	}
	return res
}

func t2_1() string {
	res, err := adventofcode2022.Task2_1(
		&adventofcode2022.FileToStringsInputReader{"adventofcode2022/day2.data"},
		adventofcode2022.ToTuple,
	)
	if err != nil {
		return err.Error()
	}
	return res
}

func t2_2() string {
	res, err := adventofcode2022.Task2_2(
		&adventofcode2022.FileToStringsInputReader{"adventofcode2022/day2.data"},
		adventofcode2022.ToTuple,
	)
	if err != nil {
		return err.Error()
	}
	return res
}
