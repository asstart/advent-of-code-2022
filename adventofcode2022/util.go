package adventofcode2022

import (
	"bufio"
	"fmt"
	"os"
)

type InputReader interface {
	GetInput() (interface{}, error)
}

type FileToStringsInputReader struct {
	Path string
}

func (fts *FileToStringsInputReader) GetInput() (interface{}, error) {
	f, err := os.Open(fts.Path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	lines := []string{}

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	if err := sc.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func ToLines(ir InputReader) []string {
	content, err := ir.GetInput()
	if err != nil {
		panic(err)
	}

	lines, ok := content.([]string)

	if !ok {
		panic(fmt.Errorf("can't convert to []string"))
	}

	return lines
}
