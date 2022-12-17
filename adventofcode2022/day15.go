package adventofcode2022

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type PointType int

const (
	NotCovered PointType = iota
	Covered
	Sensor
	Beacon
)

type Line struct {
	From int
	To   int
}

func (l Line) Contains(x int) bool {
	return x >= l.From && x <= l.To
}

type SensorsBeaconsField struct {
	Points []SensorBeacon
	MaxX   int
	MinX   int
	MaxY   int
	MinY   int
}
type SensorBeacon struct {
	Sensor Point
	Beacon Point
}

func ToSensorsBeacons(ir InputReader) (SensorsBeaconsField, error) {
	lines, err := ir.GetInput()
	if err != nil {
		return SensorsBeaconsField{}, err
	}
	res := []SensorBeacon{}
	re := regexp.MustCompile(`.+x=(-?\d+), y=(-?\d+)`)
	minX := math.MaxInt32
	maxX := math.MinInt32
	minY := math.MaxInt32
	maxY := math.MinInt32
	for _, line := range lines {
		splitted := strings.Split(line, ":")
		if len(splitted) != 2 {
			return SensorsBeaconsField{}, fmt.Errorf("expected format: [...x=1, y=2: ... x=2, y=4], got: %v", line)
		}
		foundSensorPoints := re.FindStringSubmatch(splitted[0])
		if len(foundSensorPoints) != 3 {
			return SensorsBeaconsField{}, fmt.Errorf("expexted format [...x=1, y=2], got: %v", splitted[0])
		}
		foundBeaconPoints := re.FindStringSubmatch(splitted[1])
		if len(foundBeaconPoints) != 3 {
			return SensorsBeaconsField{}, fmt.Errorf("expexted format [...x=1, y=2], got: %v", splitted[1])
		}
		sx, err := strconv.Atoi(foundSensorPoints[1])
		if err != nil {
			return SensorsBeaconsField{}, err
		}
		sy, err := strconv.Atoi(foundSensorPoints[2])
		if err != nil {
			return SensorsBeaconsField{}, err
		}
		bx, err := strconv.Atoi(foundBeaconPoints[1])
		if err != nil {
			return SensorsBeaconsField{}, err
		}
		by, err := strconv.Atoi(foundBeaconPoints[2])
		if err != nil {
			return SensorsBeaconsField{}, err
		}

		maxY = Max(maxY, Max(sy, by))
		minY = Min(minY, Min(sy, by))
		maxX = Max(maxX, Max(sx, bx))
		minX = Min(minX, Min(sx, bx))

		res = append(res, SensorBeacon{
			Sensor: Point{X: sx, Y: sy},
			Beacon: Point{X: bx, Y: by},
		})
	}
	return SensorsBeaconsField{
		Points: res,
		MaxX:   maxX,
		MinX:   minX,
		MaxY:   maxY,
		MinY:   minY,
	}, nil
}

func Task15_1(ir InputReader, cnvrtInpt func(InputReader) (SensorsBeaconsField, error), debug bool) (string, error) {
	sb, err := cnvrtInpt(ir)
	if err != nil {
		return "", err
	}

	field := map[Point]PointType{}
	for _, p := range sb.Points {
		field[p.Sensor] = Sensor
		field[p.Beacon] = Beacon
	}

	// k - Y (line num)
	// []Line - list of coverage ranges, in case if several range on a line
	coverage := map[int][]Line{}
	for _, p := range sb.Points {
		markRadius(p.Sensor, m1Distance(p.Sensor, p.Beacon), coverage)
	}

	line := 2000000
	res, err := countCovered(line, coverage, field, sb.MinX, sb.MaxX, sb.MinY)
	if err != nil {
		return "", err
	}

	if debug {
		debugD15P1(coverage)
	}

	return fmt.Sprintf("%v", res), nil
}

func Task15_2(ir InputReader, cnvrtInpt func(InputReader) (SensorsBeaconsField, error), debug bool) (string, error) {
	sb, err := cnvrtInpt(ir)
	if err != nil {
		return "", err
	}

	field := map[Point]PointType{}
	for _, p := range sb.Points {
		field[p.Sensor] = Sensor
		field[p.Beacon] = Beacon
	}

	// k - Y (line num)
	// []Line - list of coverage ranges, in case if several range on a line
	coverage := map[int][]Line{}
	for _, p := range sb.Points {
		markRadius(p.Sensor, m1Distance(p.Sensor, p.Beacon), coverage)
	}

	tx, ty := 0, 0
	for i := 0; i <= 4000000; i++ {
		cont := true
		for j := 0; j <= 4000000; j++ {
			sgs, ok := coverage[i]
			if !ok {
				continue
			}
			cont = true
			for _, sg := range sgs {
				if sg.Contains(j) {
					j = sg.To
					break
				} else {
					cont = cont && false
				}
			}
			if !cont {
				tx = j
				ty = i
				break
			}
		}
		if !cont {
			break
		}
	}

	res := int64(tx)*4000000 + int64(ty)
	return fmt.Sprintf("%v", res), nil
}

func countCovered(line int, field map[int][]Line, initField map[Point]PointType, minX int, maxX int, minY int) (int, error) {
	// count := 0

	tls, ok := field[line]
	if !ok {
		return 0, fmt.Errorf("line: %v not found", line)
	}
	// filtered := []Line{}
	length := 0
	for _, l := range tls {
		length += l.To - l.From + 1
	}

	for k, _ := range initField {
		if k.Y == line {
			length--
		}
	}

	return length, nil
}

func m1Distance(p1 Point, p2 Point) int {
	return Abs(p1.X-p2.X) + Abs(p1.Y-p2.Y)
}

func markRadius(p Point, m1 int, field map[int][]Line) {
	top := p.Y - m1
	bot := p.Y + m1
	x0 := p.X
	shiftX := 0
	shiftY := m1
	for i := top; i <= bot; i++ {
		shiftX = m1 - Abs(shiftY)
		shiftY--
		lines, ok := field[i]
		nl := Line{From: x0 - shiftX, To: x0 + shiftX}
		if !ok {
			field[i] = []Line{nl}
		} else {
			lines = mergeSegment(nl, lines)
			field[i] = lines
		}
	}
}

func mergeSegment(line Line, lines []Line) []Line {
	sort.Slice(lines, func(i int, j int) bool {
		return lines[i].From < lines[j].From
	})

	// segment starting before frist segment in arr
	if line.To < lines[0].From {
		tmp := make([]Line, 1)
		tmp[0] = line
		lines = append(tmp, lines...)
		return lines
	} else if line.To == lines[0].From {
		lines[0] = Line{line.From, lines[0].To}
	}
	// segment starting further then last segment in arr
	if line.From > lines[len(lines)-1].To {
		return append(lines, line)
	}

	overlap := 0
	overlapFrom := 0
	overlapIdx := 0
	for i := 0; i < len(lines); i++ {
		// Either line before first segment or between any two segments, without overlaping
		if overlap == 0 && line.To < lines[i].From { // full before
			if i == 0 {
				tmp := make([]Line, 1)
				tmp[0] = line
				lines = append(tmp, lines...)
			} else {
				tmp := make([]Line, len(lines[i:]))
				copy(tmp, lines[i:])
				lines = append(append(lines[:i], line), tmp...)
			}
			return lines
		} else if overlap == 0 && line.From < lines[i].From && line.To >= lines[i].From && line.To <= lines[i].To {
			newLine := Line{From: line.From, To: lines[i].To}
			lines[i] = newLine
			return lines
		} else if overlap == 0 && line.From < lines[i].From {
			overlapIdx = i
			overlap++
			overlapFrom = line.From
		} else if overlap == 0 && line.From >= lines[i].From && line.To <= lines[i].To { // full inside
			return lines
		} else if overlap == 0 && line.From >= lines[i].From && line.From <= lines[i].To { // only from inside
			overlapIdx = i
			overlapFrom = lines[i].From
			overlap++
		} else if overlap > 0 && line.To < lines[i].From {
			newLine := Line{From: overlapFrom, To: line.To}
			lines = append(append(lines[:overlapIdx], newLine), lines[i:]...)
			return lines
		} else if overlap > 0 && line.To >= lines[i].From && line.To <= lines[i].To {
			newLine := Line{From: overlapFrom, To: lines[i].To}
			lines[overlapIdx] = newLine
			lines = append(lines[:overlapIdx+1], lines[overlapIdx+overlap+1:]...)
			return lines
		} else if overlap > 0 && line.To > lines[i].To {
			overlap++
		}
	}
	if overlap > 0 {
		if line.To >= lines[overlapIdx].To {
			lines[overlapIdx] = Line{From: overlapFrom, To: line.To}
			lines = lines[:overlapIdx+1]

		}
	}

	return lines
}

func debugD15P1(field map[int][]Line) {
	minX := math.MaxInt32
	maxX := math.MinInt32
	minY := math.MaxInt32
	maxY := math.MinInt32
	for k, v := range field {
		for _, s := range v {
			minX = Min(s.From, minX)
			maxX = Max(s.To, maxX)
		}
		minY = Min(k, minY)
		maxY = Max(k, maxY)
	}

	fmt.Printf("x0: %v, y0: %v, x1: %v, y1: %v\n", minX, minY, maxX, maxY)

	for i := minY; i <= maxY; i++ {
		sgmts, ok := field[i]
		if !ok {
			for j := minX; j <= maxX; j++ {
				fmt.Printf(".")
			}
		} else {
			for j := minX; j <= maxX; j++ {
				f := false
				for _, s := range sgmts {
					if s.From <= j && j <= s.To {
						f = true
						break
					}
				}
				if f {
					fmt.Printf("#")
				} else {
					fmt.Printf(".")
				}
				f = false
			}
		}
		fmt.Println()
	}
}
