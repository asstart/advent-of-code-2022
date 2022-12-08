package adventofcode2022

import (
	"fmt"
	"strconv"
)

type VStatus int

const (
	NotVisible VStatus = iota
	NotVisibleSameHight
	Visible
)

type TreeInfo struct {
	Hight      int
	VLeft  VStatus
	DLeft  int
	VTop   VStatus
	DTop   int
	VRight VStatus
	DRight int
	VBot   VStatus
	DBot   int
}

func (v TreeInfo) visibilityArea() int {
	return v.DLeft * v.DBot * v.DTop * v.DRight
}

func To2DTreeInfoArray(ir InputReader) ([][]TreeInfo, error) {
	content, err := ir.GetInput()
	if err != nil {
		return nil, err
	}

	out := [][]TreeInfo{}
	for _, line := range content {
		row := []TreeInfo{}
		for i := 0; i < len(line); i++ {
			parsed, err := strconv.Atoi(string(line[i]))
			if err != nil {
				return nil, err
			}
			row = append(row, TreeInfo{
				Hight: parsed,
			})
		}
		out = append(out, row)
	}

	return out, nil
}

// set visibility from top of particular tree
func populateVisibilityTop(row int, col int, data [][]TreeInfo) {
	v := &data[row][col]
	top := data[row-1][col]
	if v.Hight > top.Hight {
		v.VTop = Visible
		propagateVisiblityChangeToTop(v, row-1, col, data)
	} else if v.Hight == top.Hight {
		v.VTop = NotVisibleSameHight
		v.DTop += 1
	} else {
		v.VTop = NotVisible
		v.DTop += 1
	}
}

// update visibility from particular tree to every tree to the top
func propagateVisiblityChangeToTop(initial *TreeInfo, fromRow int, col int, data [][]TreeInfo) {
	for i := fromRow; i >= 0; i-- {
		initial.DTop += 1
		if data[i][col].Hight < initial.Hight {
			data[i][col].VBot = NotVisible
		} else if data[i][col].Hight == initial.Hight {
			data[i][col].VBot = NotVisibleSameHight
			initial.VTop = NotVisible
			return
		} else {
			initial.VTop = NotVisible
			return
		}
	}
}

// set visibility from left of particular tree
func populateVisibilityLeft(row int, col int, data [][]TreeInfo) {
	v := &data[row][col]
	left := data[row][col-1]
	if v.Hight > left.Hight {
		v.VLeft = Visible
		propagateVisiblityChangeToLeft(v, row, col-1, data)
	} else if v.Hight == left.Hight {
		v.VLeft = NotVisibleSameHight
		v.DLeft += 1
	} else {
		v.VLeft = NotVisible
		v.DLeft += 1
	}
}

// update visibility from particular tree to every tree to the left
func propagateVisiblityChangeToLeft(initial *TreeInfo, row int, fromCol int, data [][]TreeInfo) {
	for i := fromCol; i >= 0; i-- {
		initial.DLeft += 1
		if data[row][i].Hight < initial.Hight {
			data[row][i].VRight = NotVisible
		} else if data[row][i].Hight == initial.Hight {
			data[row][i].VRight = NotVisibleSameHight
			initial.VLeft = NotVisible
			return
		} else {
			initial.VLeft = NotVisible
			return
		}
	}
}

// set visibility from bottom of particular tree
func populateVisibilityBot(row int, col int, data [][]TreeInfo) {
	v := &data[row][col]
	bot := data[row+1][col]
	if v.Hight > bot.Hight {
		v.VBot = Visible
		propagateVisiblityChangeToBot(v, row+1, col, data)
	} else if v.Hight == bot.Hight {
		v.VBot = NotVisibleSameHight
		v.DBot += 1
	} else {
		v.VBot = NotVisible
		v.DBot += 1
	}
}

// update visibility from particular tree to every tree to the bottom
func propagateVisiblityChangeToBot(initial *TreeInfo, fromRow int, col int, data [][]TreeInfo) {
	for i := fromRow; i < len(data); i++ {
		initial.DBot += 1
		if data[i][col].Hight < initial.Hight {
			data[i][col].VTop = NotVisible
		} else if data[i][col].Hight == initial.Hight {
			data[i][col].VTop = NotVisibleSameHight
			initial.VBot = NotVisible
			return
		} else {
			initial.VBot = NotVisible
			return
		}
	}
}

// set visibility from right of particular tree
func populateVisibilityRight(row int, col int, data [][]TreeInfo) {
	v := &data[row][col]
	right := data[row][col+1]
	if v.Hight > right.Hight {
		v.VRight = Visible
		propagateVisiblityChangeToRight(v, row, col+1, data)
	} else if v.Hight == right.Hight {
		v.DRight += 1
		v.VRight = NotVisibleSameHight
	} else {
		v.DRight += 1
		v.VRight = NotVisible
	}
}

// update visibility from particular tree to every tree to the right
func propagateVisiblityChangeToRight(initial *TreeInfo, row int, fromCol int, data [][]TreeInfo) {
	for i := fromCol; i < len(data[row]); i++ {
		initial.DRight += 1
		if data[row][i].Hight < initial.Hight {
			data[row][i].VLeft = NotVisible
		} else if data[row][i].Hight == initial.Hight {
			data[row][i].VLeft = NotVisibleSameHight
			initial.VRight = NotVisible
			return
		} else {
			initial.VRight = NotVisible
			return
		}
	}
}

// Populate information about visibility and visibility area for each tree
func populateVisibility(data [][]TreeInfo) {
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i]); j++ {
			if i == 0 {
				data[i][j].VTop = Visible
			} else if i > 0 {
				populateVisibilityTop(i, j, data)
			}

			if j == 0 {
				data[i][j].VLeft = Visible
			} else if j > 0 {
				populateVisibilityLeft(i, j, data)
			}

			if i == len(data)-1 {
				data[i][j].VBot = Visible
				propagateVisiblityChangeToTop(&data[i][j], i-1, j, data)
			} else if i < len(data)-1 {
				populateVisibilityBot(i, j, data)
			}

			if j == len(data[i])-1 {
				data[i][j].VRight = Visible
				propagateVisiblityChangeToLeft(&data[i][j], i, j-1, data)
			} else if j < len(data[i])-1 {
				populateVisibilityRight(i, j, data)
			}
		}
	}
}

func Task8_1(ir InputReader, cnvrtInpt func(ir InputReader) ([][]TreeInfo, error), debug bool) (string, error) {
	data, err := cnvrtInpt(ir)
	if err != nil {
		return "", err
	}

	populateVisibility(data)

	counter := 0
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i]); j++ {
			item := data[i][j]
			if item.VTop == Visible ||
				item.VLeft == Visible ||
				item.VBot == Visible ||
				item.VRight == Visible {
				counter++
			}
		}
	}
	if debug {
		shortDebugVisibilityOutput(data)
		fullDebugVisibilityOutput(data)
	}
	return fmt.Sprintf("Result: %v", counter), nil
}

func Task8_2(ir InputReader, cnvrtInpt func(ir InputReader) ([][]TreeInfo, error), debug bool) (string, error) {
	data, err := cnvrtInpt(ir)
	if err != nil {
		return "nil", err
	}

	populateVisibility(data)

	maxArea := 0
	for i := 1; i < len(data)-1; i++ {
		for j := 1; j < len(data[i])-1; j++ {
			item := data[i][j]
			curArea := item.visibilityArea()
			if curArea > maxArea {
				maxArea = curArea
			}
		}
	}
	if debug {
		shortDebugAreaOutput(data)
	}
	return fmt.Sprintf("Result: %v", maxArea), nil
}

func fullDebugVisibilityOutput(data [][]TreeInfo) {
	fmt.Println("Full debug")
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i]); j++ {
			item := data[i][j]
			stat := "-"
			if item.VTop == Visible ||
				item.VLeft == Visible ||
				item.VBot == Visible ||
				item.VRight == Visible {
				stat = "+"
			}
			var vt string = "-"
			var vl string = "-"
			var vb string = "-"
			var vr string = "-"
			if item.VTop == Visible {
				vt = "+"
			}
			if item.VLeft == Visible {
				vl = "+"
			}
			if item.VBot == Visible {
				vb = "+"
			}
			if item.VRight == Visible {
				vr = "+"
			}
			fmt.Printf("%v%v[t%vl%vb%vr%v]", item.Hight, stat, vt, vl, vb, vr)
		}
		fmt.Println()
	}
}

func shortDebugVisibilityOutput(data [][]TreeInfo) {
	fmt.Println("Short debug")
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i]); j++ {
			item := data[i][j]
			stat := "-"
			if item.VTop == Visible ||
				item.VLeft == Visible ||
				item.VBot == Visible ||
				item.VRight == Visible {
				stat = "+"
			}
			fmt.Printf("%v%v", item.Hight, stat)
		}
		fmt.Println()
	}
}

func shortDebugAreaOutput(data [][]TreeInfo) {
	fmt.Println("Short debug")
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i]); j++ {
			item := data[i][j]
			fmt.Printf("[%v:%2d]", item.Hight, item.visibilityArea())
		}
		fmt.Println()
	}
}
