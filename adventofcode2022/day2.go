package adventofcode2022

import (
	"fmt"
	"strings"
)

// Rock  = 1 (A, X)
// Paper = 2 (B, Y)
// Scissors = 3 (C, Z)

// Lost = 0
// Draw = 3
// Win = 6

const (
	winScore  int = 6
	drawScore int = 3
	loseScor  int = 0
)

type RPS int

const (
	R RPS = iota
	P
	S
)

func (rps RPS) Score() int {
	switch rps {
	case R:
		return 1
	case P:
		return 2
	case S:
		return 3
	}
	panic("unexpected behaviour")
}

func parseABS(s string) RPS {
	switch s {
	case "A":
		return R
	case "B":
		return P
	case "C":
		return S
	default:
		panic(fmt.Errorf("unknown value %v in parseABS", s))
	}
}

func parseXYZ(s string) RPS {
	switch s {
	case "X":
		return R
	case "Y":
		return P
	case "Z":
		return S
	default:
		panic(fmt.Errorf("unknown value %v in parseXYZ", s))
	}
}

// RPS stands for: Rock, Paper, Scissors
type TupleRPS struct {
	l RPS
	r RPS
}

func ToTupleRPSArr(ir InputReader) ([]TupleRPS, error) {
	content, err := ir.GetInput()
	if err != nil {
		return nil, err
	}

	tuples := []TupleRPS{}
	for _, l := range content {
		l = strings.TrimSpace(l)
		if l != "" {
			splt := strings.Split(l, " ")
			if len(splt) != 2 {
				return nil, fmt.Errorf("can't split line: %v properly with space separator", l)
			}
			tuples = append(tuples, TupleRPS{parseABS(splt[0]), parseXYZ(splt[1])})
		}
	}
	return tuples, nil
}

// Standard rules
func Task2_1(ir InputReader, convertInput func(ir InputReader) ([]TupleRPS, error)) (string, error) {
	source, err := convertInput(ir)
	if err != nil {
		return "", err
	}

	score := 0

	for _, t := range source {
		score += ruleset1Score(t.l, t.r)
	}

	return fmt.Sprintf("result: %v", score), nil
}

func ruleset1Score(l RPS, r RPS) int {
	if l == r {
		return drawScore + r.Score()
	} else if l == R && r == S {
		return loseScor + r.Score()
	} else if l == P && r == R {
		return loseScor + r.Score()
	} else if l == S && r == P {
		return loseScor + r.Score()
	} else {
		return winScore + r.Score()
	}
}

func Task2_2(ir InputReader, convertInput func(ir InputReader) ([]TupleRPS, error)) (string, error) {
	source, err := convertInput(ir)
	if err != nil {
		return "", err
	}

	score := 0

	for _, t := range source {
		score += ruleset2Score(t.l, t.r)
	}

	return fmt.Sprintf("result: %v", score), nil
}

// X (R) means you need to lose
// Y (P) means you need to end the round in a draw
// Z (S) means you need to win
func ruleset2Score(l RPS, r RPS) int {
	switch r {
	case R:
		return loseTo(l).Score() + resultScore(r)
	case P:
		return drawWith(l).Score() + resultScore(r)
	case S:
		return winOf(l).Score() + resultScore(r)
	}
	panic("unexpected behaviour")
}

func resultScore(r RPS) int {
	switch r {
	case R:
		return loseScor
	case P:
		return drawScore
	case S:
		return winScore
	}
	panic("unexpected behaviour")
}

func winOf(sign RPS) RPS {
	switch sign {
	case R:
		return P
	case P:
		return S
	case S:
		return R
	}
	panic("unexpected behaviour")
}

func drawWith(sign RPS) RPS {
	return sign
}

func loseTo(sign RPS) RPS {
	switch sign {
	case R:
		return S
	case P:
		return R
	case S:
		return P
	}
	panic("unexpected behaviour")
}
