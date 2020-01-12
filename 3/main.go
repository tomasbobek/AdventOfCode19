package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"sort"
	"strconv"
	"strings"
)

type Point struct {
	X int
	Y int
}

func main() {
	wireA, wireB := loadSteps()

	wirePathA := constructWirePath(wireA)
	wirePathB := constructWirePath(wireB)

	intersections := getPathIntersections(wirePathA, wirePathB)
	fmt.Println("Nearest intersection distance: ", nearestDistance(intersections))

	crossDistances := make([]int, 0, 0)
	for _, intersection := range intersections {
		crossDistances = append(crossDistances, getPathLengthToPoint(wirePathA, intersection)+getPathLengthToPoint(wirePathB, intersection))
	}
	sort.Ints(crossDistances)
	fmt.Println("Nearest intersection for wire length: ", crossDistances[0])
}

func loadSteps() ([]string, []string) {
	bytes, err := ioutil.ReadFile("C:/Users/tomas.bobek/AdventOfCode19/3/input")

	if err != nil {
		fmt.Println(err)
	}

	wires := strings.Split(string(bytes), "\r\n")

	return strings.Split(wires[0], ","), strings.Split(wires[1], ",")
}

func constructWirePath(steps []string) []Point {
	currentPos := Point{0, 0}
	path := make([]Point, 0, 0)

	for _, step := range steps {
		direction := string(step[0])
		length, err := strconv.Atoi(step[1:])

		if err != nil {
			fmt.Println(err)
		}

		for i := 1; i <= length; i++ {
			var stepPoint Point

			switch direction {
			case "U":
				stepPoint = Point{currentPos.X, currentPos.Y + 1}
			case "D":
				stepPoint = Point{currentPos.X, currentPos.Y - 1}
			case "L":
				stepPoint = Point{currentPos.X - 1, currentPos.Y}
			case "R":
				stepPoint = Point{currentPos.X + 1, currentPos.Y}
			}

			path = append(path, stepPoint)
			currentPos = stepPoint
		}
	}

	return path
}

func getPathLengthToPoint(steps []Point, final Point) int {
	lenght := 0

	for _, step := range steps {
		lenght++

		if step.X == final.X && step.Y == final.Y {
			break
		}
	}

	return lenght
}

func getPathIntersections(pathA, pathB []Point) []Point {
	intersections := make([]Point, 0, 0)

	for _, p1 := range pathA {
		for _, p2 := range pathB {
			if p1.X == p2.X && p1.Y == p2.Y {
				fmt.Println(fmt.Sprintf("Found intersection at: [%d,%d]", p1.X, p1.Y))
				intersections = append(intersections, p1)
			}
		}
	}

	return intersections
}

func nearestDistance(points []Point) int {
	distance := math.MaxInt32

	for _, point := range points {
		currentDist := int(math.Abs(float64(point.X)) + math.Abs(float64(point.Y)))
		if currentDist < distance {
			distance = currentDist
		}
	}

	return distance
}
