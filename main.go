package main

import (
	"fmt"
	"log"
	"os"

	"github.com/asstart/advent-of-code-2022/adventofcode2022"
	"github.com/jessevdk/go-flags"
)

type opts struct {
	N string `short:"n" required:"true" description:"Number of task in format day_part, like 1_1, 1_2"`
}

func main() {
	var o opts
	if _, err := flags.Parse(&o); err != nil {
		log.Printf("error while parsing flags: %s", err)
		os.Exit(1)
	}

	fmt.Printf("Running task: %v\nResult      : %v\n", o.N, runTask(o))
}

func runTask(o opts) string {
	switch o.N {
	case "1_1":
		res, err := adventofcode2022.Part1(
			&adventofcode2022.FileToStringsInputReader{"adventofcode2022/day1_pt1.data"},
			adventofcode2022.ToLines,
		)
		if err != nil {
			return err.Error()
		}
		return res
	case "1_2":
		res, err := adventofcode2022.Part2(
			&adventofcode2022.FileToStringsInputReader{"adventofcode2022/day1_pt2.data"},
			adventofcode2022.ToLines,
		)
		if err != nil {
			return err.Error()
		}
		return res
	default:
		return fmt.Errorf("function for task not found").Error()
	}
}
