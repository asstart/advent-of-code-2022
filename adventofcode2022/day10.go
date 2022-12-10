package adventofcode2022

import (
	"fmt"
	"strconv"
	"strings"
)

type CmdState struct {
	Cycle int
	Value int
	Saved []int
}

type StatefullCmd interface {
	Execute(s *CmdState, needSave func(cycle int, state int) bool, save func(cycle int, state int) int)
}

type Addx struct {
	Cycle int
	Arg   int
}

func NewAddxCmd(arg int) Addx {
	return Addx{Cycle: 2, Arg: arg}
}

func (c Addx) Execute(s *CmdState, needSave func(cycle int, state int) bool, save func(cycle int, state int) int) {
	for i := 0; i < 2; i++ {
		s.Cycle += 1
		if needSave(s.Cycle, s.Value) {
			s.Saved = append(s.Saved, save(s.Cycle, s.Value))
		}
	}
	s.Value += c.Arg
}

type Noop struct {
	Cycle int
}

func NewNoopCmd() Noop {
	return Noop{Cycle: 1}
}

func (c Noop) Execute(s *CmdState, needSave func(cycle int, state int) bool, save func(cycle int, state int) int) {
	s.Cycle += 1
	if needSave(s.Cycle, s.Value) {
		s.Saved = append(s.Saved, save(s.Cycle, s.Value))
	}
}

func ToStatefulCmds(ir InputReader) ([]StatefullCmd, error) {
	lines, err := ir.GetInput()
	if err != nil {
		return nil, err
	}

	cmds := []StatefullCmd{}
	for _, line := range lines {
		splitted := strings.Split(line, " ")
		if len(splitted) == 0 || len(splitted) > 2 {
			return nil, fmt.Errorf("expected format: <cmd> <optional: arg>, got: %v", line)
		}
		cmd := splitted[0]
		switch cmd {
		case "noop":
			cmds = append(cmds, NewNoopCmd())
		case "addx":
			if len(splitted) != 2 {
				return nil, fmt.Errorf("for [addx] cmd, expected: addx <arg>, got: %v", line)
			}
			arg, err := strconv.Atoi(splitted[1])
			if err != nil {
				return nil, err
			}
			cmds = append(cmds, NewAddxCmd(arg))
		default:
			return nil, fmt.Errorf("unexpected line: %v", line)
		}
	}
	return cmds, nil
}

func Task10_1(ir InputReader, cnvrtInpt func(InputReader) ([]StatefullCmd, error), debug bool) (string, error) {
	cmnds, err := cnvrtInpt(ir)
	if err != nil {
		return "", err
	}

	state := CmdState{Value: 1}

	needSave := func(cycle int, state int) bool {
		return cycle == 20 || (cycle-20)%40 == 0
	}
	save := func(cycle int, state int) int {
		return cycle * state
	}

	for _, cmd := range cmnds {
		cmd.Execute(&state, needSave, save)
	}

	strength := 0
	for _, s := range state.Saved {
		strength += s
	}

	return fmt.Sprintf("Result: %v", strength), nil
}

func Task10_2(ir InputReader, cnvrtInpt func(InputReader) ([]StatefullCmd, error), debug bool) (string, error) {
	cmnds, err := cnvrtInpt(ir)
	if err != nil {
		return "", err
	}

	state := CmdState{Value: 1}

	needSave := func(cycle int, state int) bool {
		currPixel := cycle - 1
		if currPixel >= 40 {
			currPixel = currPixel % 40
		}
		return currPixel == state-1 || currPixel == state || currPixel == state+1
	}
	save := func(cycle int, state int) int {
		return cycle - 1
	}
	for _, cmd := range cmnds {
		cmd.Execute(&state, needSave, save)
	}

	pic := [6][40]string{}

	for _, i := range state.Saved {
		row := i / 40
		col := i % 40
		pic[row][col] = "#"
	}

	res := strings.Builder{}
	for i := 0; i < len(pic); i++ {
		for j := 0; j < len(pic[i]); j++ {
			if pic[i][j] != "#" {
				pic[i][j] = " "
			}
			res.WriteString(pic[i][j])
		}
		res.WriteString("\n")
	}

	return fmt.Sprintf("Result:\n%v\n", res.String()), nil
}
