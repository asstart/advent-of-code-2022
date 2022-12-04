package adventofcode2022

import (
	"bufio"
	"os"
	"strings"
)

type InputReader interface {
	GetInput() ([]string, error)
}

type FileToStringsInputReader struct {
	Path string
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
		line := strings.TrimSpace(sc.Text())
		lines = append(lines, line)

	}

	if err := sc.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
