package adventofcode2022

import (
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

func DirectionOf(s string) (Direction, error) {
	switch s {
	case "U":
		return UP, nil
	case "D":
		return DOWN, nil
	case "R":
		return RIGHT, nil
	case "L":
		return LEFT, nil
	default:
		return UP, fmt.Errorf("unexpected direction: %v", s)
	}
}

type KnotMove struct {
	Direction Direction
	Count     int
}

func ToMoves(ir InputReader) ([]KnotMove, error) {
	lines, err := ir.GetInput()
	if err != nil {
		return nil, err
	}

	moves := []KnotMove{}
	for _, line := range lines {
		splitted := strings.Split(line, " ")
		if len(splitted) != 2 {
			return nil, fmt.Errorf("expected format: [Direction] [Steps], got: %v", line)
		}
		count, err := strconv.Atoi(splitted[1])
		if err != nil {
			return nil, err
		}
		direction, err := DirectionOf(splitted[0])
		if err != nil {
			return nil, err
		}
		moves = append(moves,
			KnotMove{
				Direction: direction,
				Count:     count,
			},
		)
	}
	return moves, nil
}

type Point struct {
	X int
	Y int
}

type KnotsState struct {
	Knots     []Point
	Recorders map[int]PosRecorder
}

type PosRecorder struct {
	Positions []Point
}

func (s *KnotsState) calcMove(move KnotMove) {
	for i := 0; i < move.Count; i++ {
		head := &s.Knots[0]
		s.move1Step(head, move.Direction)
		s.recordPoistion(0)
		for j := 1; j < len(s.Knots); j++ {
			hIdx := j - 1
			tIdx := j
			curHead := &s.Knots[hIdx]
			tail := &s.Knots[tIdx]
			tm := s.determineTailAction(*curHead, *tail)
			// if current knot haven't moved, other mustn't too
			if len(tm) == 0 {
				s.recordPoistion(tIdx)
				break
			}
			// there're can be 2 moves if it's diagonal move, we don't need to record intermidiate tail poistion
			// in this case, only the final
			for _, m := range tm {
				s.move1Step(tail, m.Direction)
			}
			s.recordPoistion(tIdx)
		}
	}
}

func (KnotsState) move1Step(head *Point, direction Direction) {
	switch direction {
	case UP:
		head.Y += 1
	case DOWN:
		head.Y -= 1
	case LEFT:
		head.X -= 1
	case RIGHT:
		head.X += 1
	}
}

func (s *KnotsState) recordPoistion(knot int) {
	_, ok := s.Recorders[knot]
	if !ok {
		s.Recorders[knot] = PosRecorder{Positions: []Point{}}
	}
	rcrd := s.Recorders[knot]
	rcrd.Positions = append(rcrd.Positions, s.Knots[knot])
	s.Recorders[knot] = rcrd
}

func (KnotsState) determineTailAction(h Point, t Point) []KnotMove {
	// T and H intersects
	if h.X == t.X && h.Y == t.Y {
		return []KnotMove{}
	}
	// T and H on the same horizontal line
	if h.X == t.X {
		// T and H adjacement by vertical
		if h.Y == t.Y+1 || h.Y == t.Y-1 {
			return []KnotMove{}
		} else {
			diff := h.Y - t.Y
			if diff > 0 {
				return []KnotMove{{Direction: UP, Count: 1}}
			} else {
				return []KnotMove{{Direction: DOWN, Count: 1}}
			}
		}
	}
	// T and H on the same vertical line
	if h.Y == t.Y {
		// T and H adjacement by horizontal
		if h.X == t.X+1 || h.X == t.X-1 {
			return []KnotMove{}
		} else {
			diff := h.X - t.X
			if diff > 0 {
				return []KnotMove{{Direction: RIGHT, Count: 1}}
			} else {
				return []KnotMove{{Direction: LEFT, Count: 1}}
			}
		}
	}
	// T and H diagonal
	if h.X != t.X && h.Y != t.Y {
		// T and H adjacement by diagonal
		if h.X == t.X-1 && h.Y == t.Y-1 ||
			h.X == t.X-1 && h.Y == t.Y+1 ||
			h.X == t.X+1 && h.Y == t.Y-1 ||
			h.X == t.X+1 && h.Y == t.Y+1 {
			return []KnotMove{}
		} else {
			diffX := h.X - t.X
			diffY := h.Y - t.Y
			if diffX > 0 && diffY > 0 {
				return []KnotMove{{Direction: UP, Count: 1}, {Direction: RIGHT, Count: 1}}
			} else if diffX > 0 && diffY < 0 {
				return []KnotMove{{Direction: DOWN, Count: 1}, {Direction: RIGHT, Count: 1}}
			} else if diffX < 0 && diffY > 0 {
				return []KnotMove{{Direction: UP, Count: 1}, {Direction: LEFT, Count: 1}}
			} else {
				return []KnotMove{{Direction: DOWN, Count: 1}, {Direction: LEFT, Count: 1}}
			}
		}
	}
	// Mustn't be here in any case
	panic(fmt.Errorf("unexpected knots: head:%v, tail:%v", h, t))
}

func Task9_1(ir InputReader, cnvrtInpt func(InputReader) ([]KnotMove, error), debug bool) (string, error) {
	moves, err := cnvrtInpt(ir)
	if err != nil {
		return "", err
	}

	knotsCount := 2
	knots := make([]Point, knotsCount)
	recorders := map[int]PosRecorder{}

	state := KnotsState{
		Knots:     knots,
		Recorders: recorders,
	}

	for _, move := range moves {
		state.calcMove(move)
	}

	uniqueTailPos := map[Point]bool{}
	for _, p := range state.Recorders[1].Positions {
		uniqueTailPos[p] = true
	}

	if debug {
		f, err := os.Create("debug_d9.debug")
		if err != nil {
			fmt.Printf("can't crete debug file: %v", err)
		}
		debugOutput(state.Recorders[1].Positions, f)
	}

	return fmt.Sprintf("Result: %v", len(uniqueTailPos)), nil
}

func Task9_2(ir InputReader, cnvrtInpt func(InputReader) ([]KnotMove, error), debug bool) (string, error) {
	moves, err := cnvrtInpt(ir)
	if err != nil {
		return "", err
	}

	knotsCount := 10
	knots := make([]Point, knotsCount)
	recorders := map[int]PosRecorder{}

	state := KnotsState{
		Knots:     knots,
		Recorders: recorders,
	}

	for _, move := range moves {
		state.calcMove(move)
	}

	uniqueTailPos := map[Point]bool{}
	for _, p := range state.Recorders[9].Positions {
		uniqueTailPos[p] = true
	}

	if debug {
		f, err := os.Create("debug_d9.debug")
		if err != nil {
			fmt.Printf("can't crete debug file: %v", err)
		} else {
			debugOutput(state.Recorders[9].Positions, f)
		}
	}

	return fmt.Sprintf("Result: %v", len(uniqueTailPos)), nil
}

func debugOutput(positions []Point, out io.Writer) {
	grid := [1000][1000]int{}
	x, y := 500, 500
	step := 0
	maxX, maxY := 0, 0
	minX, minY := int(math.Pow(2, 16)), int(math.Pow(2, 16))
	for _, p := range positions {
		yp := y - p.Y
		xp := x + p.X
		grid[yp][xp] += 1
		step++
		if yp < minY {
			minY = yp
		}
		if xp < minX {
			minX = xp
		}
		if yp > maxY {
			maxY = yp
		}
		if xp > maxX {
			maxX = xp
		}
	}
	s := strings.Builder{}
	for i := minY; i <= maxY; i++ {
		for j := minX; j <= maxX; j++ {
			if grid[i][j] > 0 {
				s.WriteString("#")
			} else {
				s.WriteString(".")
			}
		}
		s.WriteString("\n")
	}
	out.Write([]byte(s.String()))
}
