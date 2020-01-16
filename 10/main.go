package main

import (
    "fmt"
    "io/ioutil"
    "math"
    "strings"
)

func main() {
    asteroidMap := &AsteroidMap{}
    asteroidMap.loadMapDataFromFile("C:/Users/tomas.bobek/AdventOfCode19/10/map")

    asteroid, visibilityCount := asteroidMap.getAsteroidWithBestVisibility()

    fmt.Println(fmt.Sprintf("Asteroid is on [%d, %d]", asteroid.X, asteroid.Y))
    fmt.Println(fmt.Sprintf("The most suitable asteroid has clear visibility on %d other asteroids", visibilityCount))
}

type Point struct {
    X int
    Y int
}

func (p *Point) distanceFrom(op *Point) float64 {
    return math.Sqrt(math.Pow(float64(op.X - p.X), 2) + math.Pow(float64(op.Y - p.Y), 2))
}

func (p *Point) equals(op *Point) bool {
    return (p.X == op.X) && (p.Y == op.Y)
}

type Line struct {
    A int
    B int
    C int
    StartPoint *Point
    EndPoint   *Point
}

func (l *Line) containsPointInBetween(p *Point) bool {
    // Return immediately if the point is not on line
    if (l.A * p.X + l.B * p.Y + l.C) != 0 {
        return false
    }

    mainDistance := l.StartPoint.distanceFrom(l.EndPoint)
    return p.distanceFrom(l.StartPoint) < mainDistance && p.distanceFrom(l.EndPoint) < mainDistance
}

func getLineFromTwoPoints(p1, p2 *Point) *Line {
    a := p1.Y - p2.Y
    b := p2.X - p1.X
    c := (p1.X - p2.X) * p1.Y + (p2.Y - p1.Y) * p1.X

    return &Line{a, b, c, p1, p2}
}

type AsteroidMap struct {
    Asteroids []*Point
}

func (m *AsteroidMap) loadMapDataFromFile(file string) {
    bytes, err := ioutil.ReadFile(file)

    if err != nil {
        fmt.Println(err)
    }

    rows := strings.Split(string(bytes), "\r\n")

    if err != nil {
        fmt.Println(err)
    }

    for i, row := range rows {
        for j, value := range row {
            if string(value) == "#" {
                m.Asteroids = append(m.Asteroids, &Point{j,i})
            }
        }
    }
}

func (m *AsteroidMap) getAsteroidWithBestVisibility() (*Point, int) {
    maxVisibleAsteroids := 0
    var bestAsteroid *Point

    for _, asteroid := range m.Asteroids {
        visibleAsteroids := len(m.getAsteroidsInSight(asteroid))
        if visibleAsteroids > maxVisibleAsteroids {
            bestAsteroid = asteroid
            maxVisibleAsteroids = visibleAsteroids
        }
    }

    return bestAsteroid, maxVisibleAsteroids
}

func (m *AsteroidMap) getAsteroidsInSight(origin *Point) []*Point {
    var asteroidsInSight []*Point

    for _, asteroid := range m.Asteroids {
        if origin.equals(asteroid) {
            continue
        }

        blocked := false
        line := getLineFromTwoPoints(origin, asteroid)

        for _, asteroidToCheck := range m.Asteroids {
            if line.containsPointInBetween(asteroidToCheck) {
                blocked = true
                break
            }
        }

        if !blocked {
            asteroidsInSight = append(asteroidsInSight, asteroid)
        }
    }

    return asteroidsInSight
}
