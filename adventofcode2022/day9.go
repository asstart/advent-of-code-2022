package adventofcode2022

import (
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

type Direction int

const (
	UP Direction = iota
	DOWN
	LEFT
	RIGHT
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

type M struct {
	Direction Direction
	Count     int
}

func ToMoves(ir InputReader) ([]M, error) {
	lines, err := ir.GetInput()
	if err != nil {
		return nil, err
	}

	moves := []M{}
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
			M{
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

type THState struct {
	Knots     []Point
	Recorders map[int]PosRecorder
}

type PosRecorder struct {
	Positions []Point
}

func (s *THState) calcMove(move M) {
	for i := 0; i < move.Count; i++ {
		head := &s.Knots[0]
		s.moveHead1Step(head, move.Direction)
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
				s.moveTail1Step(tail, m.Direction)
			}
			s.recordPoistion(tIdx)
		}
	}
}

func (THState) moveHead1Step(head *Point, direction Direction) {
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

func (THState) moveTail1Step(tail *Point, direction Direction) {
	switch direction {
	case UP:
		tail.Y += 1
	case DOWN:
		tail.Y -= 1
	case LEFT:
		tail.X -= 1
	case RIGHT:
		tail.X += 1
	}
}

func (s *THState) recordPoistion(knot int) {
	_, ok := s.Recorders[knot]
	if !ok {
		s.Recorders[knot] = PosRecorder{Positions: []Point{}}
	}
	rcrd := s.Recorders[knot]
	rcrd.Positions = append(rcrd.Positions, s.Knots[knot])
	s.Recorders[knot] = rcrd
}

func (THState) determineTailAction(h Point, t Point) []M {
	// T and H intersects
	if h.X == t.X && h.Y == t.Y {
		return []M{}
	}
	// T and H on the same horizontal line
	if h.X == t.X {
		// T and H adjacement by vertical
		if h.Y == t.Y+1 || h.Y == t.Y-1 {
			return []M{}
		} else {
			diff := h.Y - t.Y
			if diff > 0 {
				return []M{M{Direction: UP, Count: 1}}
			} else {
				return []M{M{Direction: DOWN, Count: 1}}
			}
		}
	}
	// T and H on the same vertical line
	if h.Y == t.Y {
		// T and H adjacement by horizontal
		if h.X == t.X+1 || h.X == t.X-1 {
			return []M{}
		} else {
			diff := h.X - t.X
			if diff > 0 {
				return []M{M{Direction: RIGHT, Count: 1}}
			} else {
				return []M{M{Direction: LEFT, Count: 1}}
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
			return []M{}
		} else {
			diffX := h.X - t.X
			diffY := h.Y - t.Y
			if diffX > 0 && diffY > 0 {
				return []M{M{Direction: UP, Count: 1}, M{Direction: RIGHT, Count: 1}}
			} else if diffX > 0 && diffY < 0 {
				return []M{M{Direction: DOWN, Count: 1}, M{Direction: RIGHT, Count: 1}}
			} else if diffX < 0 && diffY > 0 {
				return []M{M{Direction: UP, Count: 1}, M{Direction: LEFT, Count: 1}}
			} else {
				return []M{M{Direction: DOWN, Count: 1}, M{Direction: LEFT, Count: 1}}
			}
		}
	}
	// Shouldn't be here
	panic("unexpected")
}

func Task9_1(ir InputReader, cnvrtInpt func(InputReader) ([]M, error), debug bool) (string, error) {
	moves, err := cnvrtInpt(ir)
	if err != nil {
		return "", err
	}

	knotsCount := 2
	knots := make([]Point, knotsCount)
	recorders := map[int]PosRecorder{}
	for i := 0; i < 2; i++ {
		knots[i] = Point{X: 0, Y: 0}
		recorders[i] = PosRecorder{Positions: []Point{}}
	}

	state := THState{
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

func Task9_2(ir InputReader, cnvrtInpt func(InputReader) ([]M, error), debug bool) (string, error) {
	moves, err := cnvrtInpt(ir)
	if err != nil {
		return "", err
	}

	knotsCount := 10
	knots := make([]Point, knotsCount)
	recorders := map[int]PosRecorder{}
	for i := 0; i < 2; i++ {
		knots[i] = Point{X: 0, Y: 0}
		recorders[i] = PosRecorder{Positions: []Point{}}
	}

	state := THState{
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
		}
		debugOutput(state.Recorders[9].Positions, f)
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
