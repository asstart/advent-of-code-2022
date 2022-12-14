package adventofcode2022

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
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

func bfsIter(graphMatrix [][]int, visited [][]int, start Point, finish Point, queue *[]Step) bool {
	cq := *queue
	curr := cq[0]
	cq = cq[1:]

	i, j := curr.Destination.Y, curr.Destination.X
	if finish.X == j && finish.Y == i {
		finish = Point{X: j, Y: i}
		return true
	}
	if i < len(graphMatrix)-1 && curr.Direction != UP {
		nextI, nextJ := i+1, j
		if (graphMatrix[i][j]-graphMatrix[nextI][nextJ] >= 0 || graphMatrix[i][j]-graphMatrix[nextI][nextJ] == -1) &&
			visited[nextI][nextJ] == 0 {
			cq = append(cq, Step{Destination: Point{X: nextJ, Y: nextI}, Direction: DOWN})
			visited[nextI][nextJ] = visited[i][j] + 1
		}
	}
	if j < len(graphMatrix[0])-1 && curr.Direction != LEFT {
		nextI, nextJ := i, j+1
		if (graphMatrix[i][j]-graphMatrix[nextI][nextJ] >= 0 || graphMatrix[i][j]-graphMatrix[nextI][nextJ] == -1) &&
			visited[nextI][nextJ] == 0 {
			cq = append(cq, Step{Destination: Point{X: nextJ, Y: nextI}, Direction: RIGHT})
			visited[nextI][nextJ] = visited[i][j] + 1
		}
	}
	if i > 0 && curr.Direction != DOWN {
		nextI, nextJ := i-1, j
		if (graphMatrix[i][j]-graphMatrix[nextI][nextJ] >= 0 || graphMatrix[i][j]-graphMatrix[nextI][nextJ] == -1) &&
			visited[nextI][nextJ] == 0 {
			cq = append(cq, Step{Destination: Point{X: nextJ, Y: nextI}, Direction: UP})
			visited[nextI][nextJ] = visited[i][j] + 1
		}
	}
	if j > 0 && curr.Direction != RIGHT {
		nextI, nextJ := i, j-1
		if (graphMatrix[i][j]-graphMatrix[nextI][nextJ] >= 0 || graphMatrix[i][j]-graphMatrix[nextI][nextJ] == -1) &&
			visited[nextI][nextJ] == 0 {
			cq = append(cq, Step{Destination: Point{X: nextJ, Y: nextI}, Direction: LEFT})
			visited[nextI][nextJ] = visited[i][j] + 1
		}
	}
	*queue = cq
	return false
}

func bfs(graphMatrix [][]int, start Point, finish Point) [][]int {
	visited := [][]int{}
	for i := 0; i < len(graphMatrix); i++ {
		visited = append(visited, make([]int, len(graphMatrix[0])))
	}

	queue := []Step{}

	queue = append(queue, Step{Destination: Point{X: start.X, Y: start.Y}, Direction: UP})
	visited[start.Y][start.X] = 0
	var finished bool

	for len(queue) > 0 {
		finished = bfsIter(graphMatrix, visited, start, finish, &queue)
		if finished {
			break
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

var threshold = 98

var colorMap = map[int]color.RGBA{
	97: color.RGBA{R: 44, G: 54, B: 57},
	98: color.RGBA{R: 63, G: 78, B: 79},
	99: color.RGBA{R: 162, G: 123, B: 92},

	100: color.RGBA{R: 186, G: 144, B: 110},
	101: color.RGBA{R: 186, G: 144, B: 110},
	102: color.RGBA{R: 186, G: 144, B: 110},

	103: color.RGBA{R: 186, G: 144, B: 110},
	104: color.RGBA{R: 190, G: 144, B: 110},
	105: color.RGBA{R: 186, G: 150, B: 110},

	106: color.RGBA{R: 186, G: 144, B: 115},
	107: color.RGBA{R: 195, G: 144, B: 110},
	108: color.RGBA{R: 195, G: 140, B: 110},

	109: color.RGBA{R: 190, G: 140, B: 100},
	110: color.RGBA{R: 190, G: 140, B: 102},
	111: color.RGBA{R: 95, G: 130, B: 75},

	112: color.RGBA{R: 95, G: 128, B: 75},
	113: color.RGBA{R: 95, G: 130, B: 78},
	114: color.RGBA{R: 95, G: 135, B: 78},

	115: color.RGBA{R: 95, G: 141, B: 78},
	116: color.RGBA{R: 95, G: 141, B: 78},
	117: color.RGBA{R: 95, G: 141, B: 78},
	118: color.RGBA{R: 95, G: 141, B: 78},

	119: color.RGBA{R: 73, G: 113, B: 116},
	120: color.RGBA{R: 214, G: 228, B: 229},
	121: color.RGBA{R: 214, G: 228, B: 229},
	122: color.RGBA{R: 239, G: 245, B: 245},
}

var visitedColorMap = map[int]color.RGBA{
	10:   color.RGBA{R: 228, G: 88, B: 38},
	15:   color.RGBA{R: 215, G: 88, B: 38},
	50:   color.RGBA{R: 215, G: 115, B: 57},
	80:   color.RGBA{R: 200, G: 115, B: 57},
	120:  color.RGBA{R: 190, G: 120, B: 57},
	200:  color.RGBA{R: 185, G: 110, B: 57},
	300:  color.RGBA{R: 179, G: 90, B: 30},
	400:  color.RGBA{R: 133, G: 67, B: 29},
	1000: color.RGBA{R: 27, G: 25, B: 41},
}

var colorkeys = []int{10, 15, 50, 80, 120, 200, 300, 400, 1000}

var pxlMultiplicator = 15
var (
	startPointColor = colornames.Blueviolet
	endPointColor   = colornames.Red
	currPointColor  = colornames.Orangered
)

func Task12_1V(ir InputReader, cnvrtInpt func(InputReader) (ElevationMap, error)) {
	data, err := cnvrtInpt(ir)
	if err != nil {
		panic(err)
	}

	width := len(data.Map[0]) * pxlMultiplicator
	height := len(data.Map) * pxlMultiplicator

	cfg := pixelgl.WindowConfig{
		Title:  "AoC2022 - Day 12, Part 1",
		Bounds: pixel.R(0, 0, float64(width), float64(height)),
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	var lastV *Point = &Point{}
	zero := 0
	var lastStep *int = &zero

	var visitedOnStep [][]int
	for i := 0; i < len(data.Map); i++ {
		visitedOnStep = append(visitedOnStep, make([]int, len(data.Map[0])))
	}

	tick := make(chan struct{}, 1)
	d := 6000000000 * time.Nanosecond
	var tickInt *time.Duration = &d

	go bfsV(data.Map, data.Start, data.Finish, lastV, tick, lastStep, visitedOnStep)

	go func() {
		for {
			<-time.After(*tickInt)
			tick <- struct{}{}
		}
	}()

	for !win.Closed() {
		win.Update()
		p := drawField(width, height, data, lastV.X, lastV.Y, visitedOnStep, *lastStep)
		s := pixel.NewSprite(p, p.Bounds())
		s.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

		if win.JustReleased(pixelgl.KeyUp) {
			*tickInt = *tickInt / 10
		} else if win.JustReleased(pixelgl.KeyDown) {
			*tickInt = *tickInt * 10
		}

	}
}

func darken(c color.RGBA) color.RGBA {
	nc := color.RGBA{
		R: uint8(float32(c.R) * 0.1),
		G: uint8(float32(c.G) * 0.1),
		B: uint8(float32(c.B) * 0.1),
	}
	return nc
}

func visitedPointColor(step int, visitedOn int, height int) color.RGBA {
	stepDiff := step - visitedOn
	minDiff := math.MaxInt32
	closest := 10000
	for _, item := range colorkeys {
		diff := Abs(item - stepDiff)
		if diff < minDiff {
			closest = item
			minDiff = diff
		}
	}
	c := visitedColorMap[closest]
	if height < threshold && closest > 120 {
		c = darken(c)
	}
	return c
}

func drawField(w int, h int, data ElevationMap, lx int, ly int, visitedOnStep [][]int, step int) *pixel.PictureData {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(img, img.Bounds(), &image.Uniform{colornames.Black}, image.Point{}, draw.Src)
	rand.Seed(time.Now().Unix())
	for i := 0; i < len(data.Map); i++ {
		for j := 0; j < len(data.Map[0]); j++ {
			var c color.RGBA
			if i == ly && j == lx {
				c = currPointColor
			} else if visitedOnStep[i][j] != 0 {
				c = visitedPointColor(step, visitedOnStep[i][j], data.Map[i][j])
			} else if i == data.Start.Y && j == data.Start.X {
				c = startPointColor
			} else if i == data.Finish.Y && j == data.Finish.X {
				c = endPointColor
			} else {
				c = colorMap[data.Map[i][j]]
			}
			rect := image.Rect(j*pxlMultiplicator, i*pxlMultiplicator, (j+1)*pxlMultiplicator, (i+1)*pxlMultiplicator)
			draw.Draw(img, rect, &image.Uniform{c}, image.Point{}, draw.Over)
		}
	}
	return pixel.PictureDataFromImage(img)
}

func bfsV(graphMatrix [][]int, start Point, finish Point, lastV *Point, tick <-chan struct{}, step *int, visitedOnStep [][]int) {
	visited := [][]int{}
	for i := 0; i < len(graphMatrix); i++ {
		visited = append(visited, make([]int, len(graphMatrix[0])))
	}
	queue := []Step{}

	queue = append(queue, Step{Destination: Point{X: start.X, Y: start.Y}, Direction: UP})
	visited[start.Y][start.X] = 0

	var finished bool
	for len(queue) > 0 {
		select {
		case <-tick:
			curr := queue[0]
			*step++
			visitedOnStep[curr.Destination.Y][curr.Destination.X] = *step

			*lastV = curr.Destination

			finished = bfsIter(graphMatrix, visited, start, finish, &queue)
			if finished {
				return
			}
		}
	}
}
