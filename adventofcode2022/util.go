package adventofcode2022

import (
	"bufio"
	"os"
	"strings"
)

type InputReader interface {
	GetInput() ([]string, error)
}

type Options struct {
	TrimLine bool
}

type FileToStringsInputReader struct {
	Path string
	Opts Options
}

func (fts *FileToStringsInputReader) GetInput() ([]string, error) {
	f, err := os.Open(fts.Path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	lines := []string{}

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if fts.Opts.TrimLine {
			line = strings.TrimSpace(line)
		}
		lines = append(lines, line)
	}

	if err := sc.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func Max(arg1 int, arg2 int) int {
	if arg1 > arg2 {
		return arg1
	}
	return arg2
}

func Min(arg1 int, arg2 int) int {
	if arg1 < arg2 {
		return arg1
	}
	return arg2
}

func Abs(arg int) int {
	if arg > 0 {
		return arg
	}
	return arg * -1
}

func Swap(arg1 *int, arg2 *int) {
	tmp := *arg2
	*arg2 = *arg1
	*arg1 = tmp
}
