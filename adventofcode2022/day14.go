package adventofcode2022

import (
	"fmt"
	"image"
	"image/draw"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var mostLeft Point = Point{X: math.MaxInt32, Y: math.MaxInt32}
var mostRight Point = Point{X: -1, Y: math.MaxInt32}
var lowest Point = Point{X: math.MaxInt32}

var zeroPoint = Point{X: 500, Y: 0}

type SandGrainState int

const (
	Falling SandGrainState = iota
	Fell
	InfinityFalling
)

type Cave struct {
	Rocks map[Point]byte
	MinX  int
	MaxX  int
	// lowest point
	MaxY int
	//highest point
	MinY int
}

func ToRockMap(ir InputReader) (Cave, error) {
	lines, err := ir.GetInput()
	if err != nil {
		return Cave{}, err
	}
	minX := math.MaxInt32
	maxX := -1
	maxY := -1
	field := map[Point]byte{}
	for _, line := range lines {
		splitted := strings.Split(line, "->")
		if len(splitted) == 0 {
			return Cave{}, fmt.Errorf("expected at least one pair of [x,y -> x,y], got: %v", line)
		}
		var from string
		var to string
		for i := 0; i < len(splitted)-1; i++ {
			from = strings.TrimSpace(splitted[i])
			to = strings.TrimSpace(splitted[i+1])

			fSplitted := strings.Split(from, ",")
			if len(fSplitted) != 2 {
				return Cave{}, fmt.Errorf("expected point in format [x,y], got: %v", from)
			}
			fx, err := strconv.Atoi(fSplitted[0])
			if err != nil {
				return Cave{}, fmt.Errorf("expected X to be number, go: %v", fSplitted[0])
			}
			fy, err := strconv.Atoi(fSplitted[1])
			if err != nil {
				return Cave{}, fmt.Errorf("expected Y to be number, go: %v", fSplitted[1])
			}
			tSplitted := strings.Split(to, ",")
			if len(tSplitted) != 2 {
				return Cave{}, fmt.Errorf("expected point in format [x,y], got: %v", to)
			}
			tx, err := strconv.Atoi(tSplitted[0])
			if err != nil {
				return Cave{}, fmt.Errorf("expected X to be number, go: %v", tSplitted[0])
			}
			ty, err := strconv.Atoi(tSplitted[1])
			if err != nil {
				return Cave{}, fmt.Errorf("expected Y to be number, go: %v", tSplitted[1])
			}

			if fx == tx {
				if fy > ty {
					Swap(&fy, &ty)
				}
				for i := fy; i <= ty; i++ {
					maxY = Max(maxY, i)
					p := Point{X: fx, Y: i}
					field[p] = 1
				}
			}

			if fy == ty {
				if fx > tx {
					Swap(&fx, &tx)
				}
				for i := fx; i <= tx; i++ {
					minX = Min(minX, i)
					maxX = Max(maxX, i)
					p := Point{X: i, Y: fy}
					field[p] = 1
				}
			}
		}
	}
	c := Cave{
		Rocks: field,
		MinX:  minX,
		MaxX:  maxX,
		MinY:  0,
		MaxY:  maxY,
	}
	return c, nil
}

func Task14_1(ir InputReader, cnvrtInpt func(InputReader) (Cave, error), debug bool) (string, error) {
	cave, err := cnvrtInpt(ir)
	if err != nil {
		return "", err
	}

	count := emulateSandfallWithInfinityFloor(&cave)

	if debug {
		f, err := os.OpenFile("day14p1_debug.debug", os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return "", err
		}
		debugD14p1(cave, f)
	}
	return fmt.Sprintf("%v", count), nil
}

func Task14_2(ir InputReader, cnvrtInpt func(InputReader) (Cave, error), debug bool) (string, error) {
	cave, err := cnvrtInpt(ir)
	if err != nil {
		return "", err
	}

	cave.MaxY += 2

	count := emulateSanfallWithFloor(&cave)

	if debug {
		f, err := os.OpenFile("day14p2_debug.debug", os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return "", err
		}
		debugD14p2(cave, f)
	}
	return fmt.Sprintf("%v", count), nil
}

func Task14_2V(ir InputReader, cnvrtInpt func(InputReader) (Cave, error)) {
	cave, err := cnvrtInpt(ir)
	if err != nil {
		panic(err)
	}

	cave.MaxY += 2

	mltpl := 3
	shitf := 170
	width := ((cave.MaxX + shitf) - (cave.MinX - shitf)) * mltpl
	height := (cave.MaxY) * mltpl
	shiftX := cave.MinX - shitf

	field := make([][]int, height)
	for i := 0; i < height; i++ {
		field[i] = make([]int, width)
		for j := 0; j < width; j++ {
			_, ok := cave.Rocks[Point{X: j + shiftX, Y: i}]
			if ok {
				field[i][j] = 1
			}
		}
	}

	cfg := pixelgl.WindowConfig{
		Title:  "AoC2022 - Day 14, Part 2",
		Bounds: pixel.R(0, 0, float64(width), float64(height)),
		VSync: true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// sig := make(chan Point)

	ticker := time.Tick(time.Nanosecond * 10000)
	go emulateSanfallWithFloorSync(&cave, field, shiftX, ticker)
	for !win.Closed() {
		win.Update()
		p := drawCaveField(field, width, height, mltpl)
		s := pixel.NewSprite(p, p.Bounds())
		m := pixel.IM.Moved(win.Bounds().Center())
		s.Draw(win, m)
	}

}

var (
	rocks = colornames.Gray
	sand  = colornames.Sandybrown
	back  = colornames.Wheat
)

func drawCaveField(field [][]int, w int, h int, pxlMult int) *pixel.PictureData {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(img, img.Bounds(), &image.Uniform{back}, image.Point{}, draw.Src)
	for i := 0; i < len(field); i++ {
		for j := 0; j < len(field[i]); j++ {
			c := back
			if field[i][j] == 1 {
				c = rocks
			} else if field[i][j] == 2 {
				c = sand
			}
			rect := image.Rect(j*pxlMult, i*pxlMult, (j+1)*pxlMult, (i+1)*pxlMult)
			draw.Draw(img, rect, &image.Uniform{c}, image.Point{}, draw.Over)
		}
	}
	return pixel.PictureDataFromImage(img)
}

func emulateSandfallWithInfinityFloor(cave *Cave) int {
	var currPoint Point
	nextPoint := zeroPoint
	prevGrainState := Falling
	counter := 0
	for {
		if prevGrainState == InfinityFalling {
			break
		}
		currPoint = nextPoint

		nextPoint, prevGrainState = emulateSandGrainFallIfNoFloor(currPoint, cave)
		if prevGrainState == Fell {
			counter++
			nextPoint = zeroPoint
		}

	}
	return counter
}

func emulateSanfallWithFloorSync(cave *Cave, field [][]int, shiftX int, delay <-chan time.Time) int {
	var currPoint Point
	nextPoint := zeroPoint
	prevGrainState := Falling
	counter := 0
	for {
		<-delay
		currPoint = nextPoint
		nextPoint, prevGrainState = emulateSandGrainFall(currPoint, cave)
		if prevGrainState == Fell {
			counter++
			// if grain fell and fell position was the same as start
			if currPoint == zeroPoint {
				break
			}
			nextPoint = zeroPoint
		}
		field[currPoint.Y][currPoint.X-shiftX] = 2
	}
	return counter
}

func emulateSanfallWithFloor(cave *Cave) int {
	var currPoint Point
	nextPoint := zeroPoint
	prevGrainState := Falling
	counter := 0
	for {
		currPoint = nextPoint
		nextPoint, prevGrainState = emulateSandGrainFall(currPoint, cave)
		if prevGrainState == Fell {
			counter++
			// if grain fell and fell position was the same as start
			if currPoint == zeroPoint {
				break
			}
			nextPoint = zeroPoint
		}

	}
	return counter
}

func isInfinityFalling(point Point, field map[Point]byte) bool {
	for k, _ := range field {
		if k.X == point.X && k.Y >= point.Y {
			return false
		}
	}
	return true
}

func isDiagonalMoveAvlb(point Point, cave *Cave) bool {
	//check if destination point is free
	_, ok := cave.Rocks[point]
	if ok {
		return false
	}

	return point.Y < cave.MaxY
}

func isDownMoveAvlb(point Point, cave *Cave) bool {
	//check if destination point is free
	_, ok := cave.Rocks[point]
	if ok {
		return false
	}

	return point.Y < cave.MaxY
}

func emulateSandGrainFallIfNoFloor(startPoint Point, cave *Cave) (Point, SandGrainState) {
	var next Point
	//check down
	next = Point{X: startPoint.X, Y: startPoint.Y + 1}
	_, ok := cave.Rocks[next]
	if !ok {
		return next, Falling
	}

	//check left
	next = Point{X: startPoint.X - 1, Y: startPoint.Y + 1}
	if isInfinityFalling(next, cave.Rocks) {
		return next, InfinityFalling
	}
	if isDiagonalMoveAvlb(next, cave) {
		return next, Falling
	}
	//check right
	next = Point{X: startPoint.X + 1, Y: startPoint.Y + 1}
	if isInfinityFalling(next, cave.Rocks) {
		return next, InfinityFalling
	}
	if isDiagonalMoveAvlb(next, cave) {
		return next, Falling
	}
	cave.Rocks[startPoint] = 2
	return Point{}, Fell
}

func emulateSandGrainFall(startPoint Point, cave *Cave) (Point, SandGrainState) {
	var next Point
	//check down
	next = Point{X: startPoint.X, Y: startPoint.Y + 1}
	if isDownMoveAvlb(next, cave) {
		return next, Falling
	}
	//check left
	next = Point{X: startPoint.X - 1, Y: startPoint.Y + 1}
	if isDiagonalMoveAvlb(next, cave) {
		cave.MinX = Min(cave.MinX, next.X)
		return next, Falling
	}
	//check right
	next = Point{X: startPoint.X + 1, Y: startPoint.Y + 1}
	if isDiagonalMoveAvlb(next, cave) {
		cave.MaxX = Max(cave.MaxX, next.X)
		return next, Falling
	}
	cave.Rocks[startPoint] = 2
	return Point{}, Fell
}

func debugD14p1(cave Cave, writer io.Writer) {
	minX := cave.MinX - 2
	maxX := cave.MaxX + 2
	minY := cave.MinY
	maxY := cave.MaxY

	builder := strings.Builder{}
	for i := minY; i <= maxY; i++ {
		for j := minX; j <= maxX; j++ {
			if j == 500 && i == 0 {
				builder.WriteString("!\n")
				continue
			}
			v, ok := cave.Rocks[Point{X: j, Y: i}]
			if !ok {
				builder.WriteString(".")
				continue
			}
			if v == 1 {
				builder.WriteString("#")
				continue
			} else {
				builder.WriteString("o")
				continue
			}
		}
		builder.WriteString("\n")
	}
	writer.Write([]byte(builder.String()))
}

func debugD14p2(cave Cave, writer io.Writer) {
	minX := cave.MinX - 2
	maxX := cave.MaxX + 2
	minY := cave.MinY
	maxY := cave.MaxY

	builder := strings.Builder{}
	for i := minY; i <= maxY; i++ {
		for j := minX; j <= maxX; j++ {
			if j == 500 && i == 0 {
				builder.WriteString("!")
				continue
			} else if i == maxY {
				builder.WriteString("=")
				continue
			}

			v, ok := cave.Rocks[Point{X: j, Y: i}]
			if !ok {
				builder.WriteString(".")
				continue
			}
			if v == 1 {
				builder.WriteString("#")
				continue
			} else {
				builder.WriteString("o")
				continue
			}
		}
		builder.WriteString("\n")
	}
	writer.Write([]byte(builder.String()))
}
