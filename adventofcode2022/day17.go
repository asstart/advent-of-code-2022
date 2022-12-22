package adventofcode2022

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

var (
	Hlineb    = []int{15 << 1}
	Vlineb    = []int{1 << 4, 1 << 4, 1 << 4, 1 << 4}
	Xb        = []int{2 << 2, 7 << 2, 2 << 2}
	Reverselb = []int{1 << 2, 1 << 2, 7 << 2}
	Squareb   = []int{3 << 3, 3 << 3}

	forderb = [][]int{Hlineb, Xb, Reverselb, Vlineb, Squareb}
)

func ToDirections(ir InputReader) ([]Direction, error) {
	lines, err := ir.GetInput()
	if err != nil {
		return nil, err
	}
	if len(lines) != 1 {
		return nil, fmt.Errorf("expected 1 line, got: %v", len(lines))
	}

	dirs := []Direction{}
	for _, s := range lines[0] {
		switch s {
		case '<':
			dirs = append(dirs, LEFT)
		case '>':
			dirs = append(dirs, RIGHT)
		default:
			return nil, fmt.Errorf("unexpected symbol: %v", s)
		}
	}
	return dirs, nil
}

// This approach can't handle number of iteration needed for 2nd part, possibly need to find repeatable pattern
func Task17_1(ir InputReader, cnvrtInpt func(InputReader) ([]Direction, error), debug bool) (string, error) {
	dirs, err := cnvrtInpt(ir)
	if err != nil {
		return "", err
	}

	figCounter := 0
	dirsCounter := 0
	height := 100 // half size of matrix which will be kept in memory
	field := initFieldB(height)

	// current shape of figure
	var shape []int = make([]int, len(forderb[figCounter]))
	copy(shape, forderb[figCounter])

	// Y coordinate of figure
	fbotidx := height - 4
	prevTop := height
	// Height of block's tower
	floors := 0

	figCount := 2022 //100000000
	var stop bool

	for i := 0; i < figCount; {
		d := dirs[dirsCounter]
		fbotidx, stop = emulateMoveB(field, shape, fbotidx, d)

		if len(dirs)-1 == dirsCounter {
			dirsCounter = 0
		} else {
			dirsCounter++
		}
		if stop {
			if len(forderb)-1 == figCounter {
				figCounter = 0
			} else {
				figCounter++
			}

			i++
			nTop := Min(prevTop, fbotidx-len(shape)+1)
			floors += prevTop - nTop
			prevTop = nTop
			fbotidx = nTop - 4

			field, _ = extendFieldIfNeed(field, fbotidx, height, &prevTop)

			shape = make([]int, len(forderb[figCounter]))
			copy(shape, forderb[figCounter])
		}
	}

	return fmt.Sprintf("%v", floors), nil
}

func initFieldB(h int) []int {
	arr := make([]int, h+1)
	arr[h] = 127
	return arr
}

func extendFieldIfNeed(field []int, fbotIdx int, h int, prevTop *int) ([]int, int) {
	nf := field
	if fbotIdx < 20 { // 51 is random threshold here, can be configured
		nf = make([]int, h)
		nf = append(nf, field[:h]...)
		fbotIdx = fbotIdx + h
		*prevTop = *prevTop + h
	}
	return nf, fbotIdx
}

func emulateMoveB(field []int, shape []int, fbotidx int, direction Direction) (int, bool) {
	switch direction {
	case LEFT:
		return moveLeftIfCan(field, shape, fbotidx)
	case RIGHT:
		return moveRightIfCan(field, shape, fbotidx)
	}
	panic(fmt.Errorf("unexpected situation, supported only [left, right] directions, got: %v", direction))
}

func moveLeftIfCan(field []int, shape []int, fbotidx int) (int, bool) {
	move := true
	var l int
	for idx, p := range shape {
		l = bitLen(p)
		if l == 7 {
			move = false
			break
		}
		check := 1 << l
		if field[fbotidx-len(shape)+1+idx]&check != 0 {
			move = false
			break
		}
	}
	if move {
		for i := len(shape) - 1; i >= 0; i-- {
			shape[i] = shape[i] << 1
		}
	}
	return moveDownIfCan(field, shape, fbotidx)
}

func moveRightIfCan(field []int, shape []int, fbotidx int) (int, bool) {
	move := true
	for idx, p := range shape {

		// check if we already in the most right position: [0000001]
		if p%2 == 1 {
			move = false
			break
		}
		var cr int
		var lastZeroIdx int = 0
		// find position of last zero before 1 in binary representation of number
		for i := 0; i < 7; i++ {
			cr = p & (1 << i)
			if cr == 0 {
				lastZeroIdx = i
			} else {
				break
			}
		}
		// check if move current number 1 got intersection or not
		if (1<<lastZeroIdx)&field[fbotidx-len(shape)+1+idx] != 0 {
			move = false
			break
		}
	}
	if move {
		for i := len(shape) - 1; i >= 0; i-- {
			shape[i] = shape[i] >> 1
		}
	}
	return moveDownIfCan(field, shape, fbotidx)
}

func moveDownIfCan(field []int, shape []int, fbotidx int) (int, bool) {
	for idx, p := range shape {
		check := p & field[fbotidx-len(shape)+idx+1+1]
		if check != 0 {
			fillField(field, shape, fbotidx)
			return fbotidx, true
		}
	}
	return fbotidx + 1, false
}

func fillField(field []int, shape []int, fbotidx int) {
	for idx, p := range shape {
		field[fbotidx-len(shape)+idx+1] = field[fbotidx-len(shape)+idx+1] | p
	}
}

// simple comparison works much faster, then calculating log2(N + 1)
func bitLen(n int) int {
	if n <= 1 {
		return 1
	} else if n <= 3 {
		return 2
	} else if n <= 7 {
		return 3
	} else if n <= 15 {
		return 4
	} else if n <= 31 {
		return 5
	} else if n <= 63 {
		return 6
	} else {
		return 7
	}
	// return int(math.Ceil(math.Log2(float64(n + 1))))
}

func debugP2B(field []int, i int, writer io.Writer, top int) {
	bld := strings.Builder{}
	bld.WriteString(fmt.Sprintf("%v---\n", i))
	skip := true
	var upto = top + 51
	if upto > len(field) {
		upto = len(field)
	}
	for _, i := range field[top:upto] {
		if i != 0 {
			skip = false
		}
		if !skip {
			bld.WriteString(fmt.Sprintf("%07v\n", strconv.FormatInt(int64(i), 2)))
		}
	}
	writer.Write([]byte(bld.String()))
}
