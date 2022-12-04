package adventofcode2022

import (
	"fmt"
	"strconv"
	"strings"
)

type TupleSegment struct {
	_1 Segment
	_2 Segment
}
type Segment struct {
	l int
	r int
}

func ToTupleSegment(ir InputReader) []TupleSegment {
	content, err := ir.GetInput()
	if err != nil {
		panic(err)
	}

	lines, ok := content.([]string)
	if !ok {
		panic(fmt.Errorf("can't cast source to []string"))
	}

	converted := []TupleSegment{}

	for _, line := range lines {
		splitted := strings.Split(line, ",")
		if len(splitted) != 2 {
			panic(fmt.Sprintf("got: %v, expected format: 1-1,2-2", line))
		}
		s1 := strings.Split(splitted[0], "-")
		if len(s1) != 2 {
			panic(fmt.Sprintf("got: %v, expected format: 1-1", splitted[0]))
		}
		s2 := strings.Split(splitted[1], "-")
		if len(s2) != 2 {
			panic(fmt.Sprintf("got: %v, expected format: 1-1", splitted[1]))
		}

		s1L, err := strconv.Atoi(s1[0])
		if err != nil {
			panic(err)
		}
		s1R, err := strconv.Atoi(s1[1])
		if err != nil {
			panic(err)
		}
		s2L, err := strconv.Atoi(s2[0])
		if err != nil {
			panic(err)
		}
		s2R, err := strconv.Atoi(s2[1])
		if err != nil {
			panic(err)
		}
		converted = append(converted,
			TupleSegment{
				_1: Segment{l: s1L, r: s1R},
				_2: Segment{l: s2L, r: s2R},
			},
		)
	}

	return converted
}

func Task4_1(ir InputReader, convInput func(InputReader) []TupleSegment) (string, error) {
	data := convInput(ir)

	count := 0

	for _, s := range data {
		if s._1.l-s._2.l <= 0 && s._1.r-s._2.r >= 0 ||
			s._1.l-s._2.l >= 0 && s._1.r-s._2.r <= 0 {
			count++
		}
	}

	return fmt.Sprintf("result: %v", count), nil
}

func Task4_2(ir InputReader, convInput func(InputReader) []TupleSegment) (string, error) {
	data := convInput(ir)

	count := 0

	for _, s := range data {
		if s._1.l <= s._2.l && s._1.r >= s._2.l ||
			s._2.l <= s._1.l && s._2.r >= s._1.l {
			count++
		}
	}

	return fmt.Sprintf("result: %v", count), nil
}
