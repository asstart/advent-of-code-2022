package adventofcode2022

import (
	"fmt"
	"strconv"
	"strings"
)

type Point3D struct {
	X int
	Y int
	Z int
}

type Side int

type PointStatus struct {
	LeftFree  bool
	RightFree bool
	UpFree    bool
	DownFree  bool
	FrontFree bool
	BackFree  bool
}

func ToArrPoint3D(ir InputReader) ([]Point3D, error) {
	lines, err := ir.GetInput()
	if err != nil {
		return nil, err
	}

	points := make([]Point3D, len(lines))
	for i, line := range lines {
		splitted := strings.Split(line, ",")
		if len(splitted) != 3 {
			return nil, fmt.Errorf("expected: [1,1,1] got: %v", line)
		}
		x, err := strconv.Atoi(splitted[0])
		if err != nil {
			return nil, err
		}
		y, err := strconv.Atoi(splitted[1])
		if err != nil {
			return nil, err
		}
		z, err := strconv.Atoi(splitted[2])
		if err != nil {
			return nil, err
		}
		points[i] = Point3D{
			X: x,
			Y: y,
			Z: z,
		}
	}

	return points, nil
}

func Task18_1(ir InputReader, cnvrtInp func(InputReader) ([]Point3D, error), debug bool) (string, error) {
	points, err := cnvrtInp(ir)
	if err != nil {
		return "", err
	}

	storage := map[Point3D]bool{}
	freeSides := 0

	for _, p := range points {
		adjps := generateAdjacementCoords(p)
		for _, adjp := range adjps {
			_, ok := storage[adjp]
			if !ok {
				freeSides += 1
			} else {
				freeSides -= 1
			}
		}
		storage[p] = true
	}

	return fmt.Sprintf("%v", freeSides), nil
}

func generateAdjacementCoords(p Point3D) []Point3D {
	return []Point3D{
		Point3D{X: p.X + 1, Y: p.Y, Z: p.Z},
		Point3D{X: p.X - 1, Y: p.Y, Z: p.Z},
		Point3D{X: p.X, Y: p.Y + 1, Z: p.Z},
		Point3D{X: p.X, Y: p.Y - 1, Z: p.Z},
		Point3D{X: p.X, Y: p.Y, Z: p.Z + 1},
		Point3D{X: p.X, Y: p.Y, Z: p.Z - 1},
	}
}
