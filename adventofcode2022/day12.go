package adventofcode2022

import (
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

type Step struct {
	Destination Point
	Direction   Direction
}

type ElevationMap struct {
	Start  Point
	Finish Point
	Map    [][]int
}

func ToElevationMap(ir InputReader) (ElevationMap, error) {
	lines, err := ir.GetInput()
	if err != nil {
		return ElevationMap{}, err
	}
	res := ElevationMap{
		Map: [][]int{},
	}
	i, j := 0, 0
	for _, line := range lines {
		row := []int{}
		j = 0
		for _, symb := range line {
			if symb == 'S' {
				res.Start = Point{X: j, Y: i}
				row = append(row, 'a')
			} else if symb == 'E' {
				res.Finish = Point{X: j, Y: i}
				row = append(row, 'z')
			} else {
				row = append(row, int(symb))
			}
			j++
		}
		i++
		res.Map = append(res.Map, row)
	}
	return res, nil
}

func bfs(graphMatrix [][]int, start Point, finish Point) [][]int {
	visited := [][]int{}
	for i := 0; i < len(graphMatrix); i++ {
		visited = append(visited, make([]int, len(graphMatrix[0])))
	}

	queue := []Step{}

	queue = append(queue, Step{Destination: Point{X: start.X, Y: start.Y}, Direction: UP})
	visited[start.Y][start.X] = 0

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		i, j := curr.Destination.Y, curr.Destination.X
		if finish.X == j && finish.Y == i {
			finish = Point{X: j, Y: i}
			break
		}
		if i < len(graphMatrix)-1 && curr.Direction != UP {
			nextI, nextJ := i+1, j
			if (graphMatrix[i][j]-graphMatrix[nextI][nextJ] >= 0 || graphMatrix[i][j]-graphMatrix[nextI][nextJ] == -1) &&
				visited[nextI][nextJ] == 0 {
				queue = append(queue, Step{Destination: Point{X: nextJ, Y: nextI}, Direction: DOWN})
				visited[nextI][nextJ] = visited[i][j] + 1
			}
		}
		if j < len(graphMatrix[0])-1 && curr.Direction != LEFT {
			nextI, nextJ := i, j+1
			if (graphMatrix[i][j]-graphMatrix[nextI][nextJ] >= 0 || graphMatrix[i][j]-graphMatrix[nextI][nextJ] == -1) &&
				visited[nextI][nextJ] == 0 {
				queue = append(queue, Step{Destination: Point{X: nextJ, Y: nextI}, Direction: RIGHT})
				visited[nextI][nextJ] = visited[i][j] + 1
			}
		}
		if i > 0 && curr.Direction != DOWN {
			nextI, nextJ := i-1, j
			if (graphMatrix[i][j]-graphMatrix[nextI][nextJ] >= 0 || graphMatrix[i][j]-graphMatrix[nextI][nextJ] == -1) &&
				visited[nextI][nextJ] == 0 {
				queue = append(queue, Step{Destination: Point{X: nextJ, Y: nextI}, Direction: UP})
				visited[nextI][nextJ] = visited[i][j] + 1
			}
		}
		if j > 0 && curr.Direction != RIGHT {
			nextI, nextJ := i, j-1
			if (graphMatrix[i][j]-graphMatrix[nextI][nextJ] >= 0 || graphMatrix[i][j]-graphMatrix[nextI][nextJ] == -1) &&
				visited[nextI][nextJ] == 0 {
				queue = append(queue, Step{Destination: Point{X: nextJ, Y: nextI}, Direction: LEFT})
				visited[nextI][nextJ] = visited[i][j] + 1
			}
		}
	}
	return visited
}

func Task12_1(ir InputReader, cnvrtInpt func(ir InputReader) (ElevationMap, error), debug bool) (string, error) {
	data, err := cnvrtInpt(ir)
	if err != nil {
		return "", err
	}

	visited := bfs(data.Map, data.Start, data.Finish)

	if debug {
		f, err := os.Create("debug_12d_letters.debug")
		if err != nil {
			fmt.Printf("can't print debug: %v\n", err)
		} else {
			debugShowLetters(data.Map, visited, f)
		}
		f, err = os.Create("debug_12d_path.debug")
		if err != nil {
			fmt.Printf("can't print debug: %v\n", err)
		} else {
			debugShowPathL(visited, f)
		}
	}

	return fmt.Sprintf("Result: %v", visited[data.Finish.Y][data.Finish.X]), nil
}

func Task12_2(ir InputReader, cnvrtInpt func(ir InputReader) (ElevationMap, error), debug bool) (string, error) {
	data, err := cnvrtInpt(ir)
	if err != nil {
		return "", err
	}

	var minVisited [][]int

	minPath := math.MaxInt32
	for i := 0; i < len(data.Map); i++ {
		for j := 0; j < len(data.Map[i]); j++ {
			if ((i == 0 || i == len(data.Map)-1) && data.Map[i][j] == 'a') ||
				((j == 0 || j == len(data.Map[i])-1) && data.Map[i][j] == 'a') {
				currVisited := bfs(data.Map, Point{X: j, Y: i}, data.Finish)
				currPath := currVisited[data.Finish.Y][data.Finish.X]
				if currPath != 0 {
					if currPath < minPath {
						minPath = currPath
						minVisited = currVisited
					}
				}
			}
		}
	}

	if debug {
		f, err := os.Create("debug_12d_letters.debug")
		if err != nil {
			fmt.Printf("can't print debug: %v\n", err)
		} else {
			debugShowLetters(data.Map, minVisited, f)
		}
		f, err = os.Create("debug_12d_path.debug")
		if err != nil {
			fmt.Printf("can't print debug: %v\n", err)
		} else {
			debugShowPathL(minVisited, f)
		}
	}

	return fmt.Sprintf("Result: %v", minPath), nil
}

func debugShowLetters(letters [][]int, visited [][]int, writer io.Writer) {
	builder := strings.Builder{}
	for i := 0; i < len(visited); i++ {
		for j := 0; j < len(visited[i]); j++ {
			if visited[i][j] == 0 {
				builder.WriteString(".")
			} else {
				builder.WriteString(fmt.Sprintf("%v", string(rune(letters[i][j]))))
			}
		}
		builder.WriteString("\n")
	}
	writer.Write([]byte(builder.String()))
}

func debugShowPathL(visited [][]int, writer io.Writer) {
	builder := strings.Builder{}
	for i := 0; i < len(visited); i++ {
		for j := 0; j < len(visited[i]); j++ {
			if visited[i][j] == 0 {
				builder.WriteString(" .. ")
			} else {
				builder.WriteString(fmt.Sprintf("%3d ", visited[i][j]))
			}
		}
		builder.WriteString("\n")
	}
	writer.Write([]byte(builder.String()))
}
