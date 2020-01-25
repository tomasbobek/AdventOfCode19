package main

import (
    "fmt"
    "io/ioutil"
    "math"
    "os"
    "sort"
    "strings"
)

func main() {
    path, err := os.Getwd()
    if err != nil {
        fmt.Println(err)
    }

    // -----------------------------------------------------------------------------------------------------------------
    // Here we solve problem for Part One
    asteroidMap := &AsteroidMap{}
    asteroidMap.loadMapDataFromFile(path + "/10/map")

    asteroid, visibilityCount := asteroidMap.getAsteroidWithBestVisibility()

    fmt.Println(fmt.Sprintf("Asteroid is on [%d, %d]", asteroid.X, asteroid.Y))
    fmt.Println(fmt.Sprintf("The most suitable asteroid has clear visibility on %d other asteroids", visibilityCount))

    // -----------------------------------------------------------------------------------------------------------------
    // Here we solve problem for Part One
    laserGun := &LaserGun{AsteroidMap: asteroidMap, BaseStation: asteroid, HitAsteroidsMap: make(map[*Point]bool)}
    laserGun.calculateAsteroidAngles()
    laserGun.divideAsteroidsToQuadrants()

    quadrant := 1

    // We'll shoot until there is no asteroid left (apart the base station :)).
    for len(laserGun.HitAsteroids) < len(asteroidMap.Asteroids) - 1 {
        laserGun.shootQuadrant(quadrant)

        if quadrant == 4 {
            quadrant = 1
        } else {
            quadrant++
        }
    }

    // 200th asteroid is located at index 199
    goal := laserGun.HitAsteroids[199]
    fmt.Println(fmt.Sprintf("200th asteroid vaporized is located at [%d,%d], code is %d", goal.X, goal.Y, goal.X * 100 + goal.Y))
}

type Asteroid struct {
    Coordinates   *Point
    AngleFromBase float64
}

type AsteroidList []Asteroid

func (al AsteroidList) Len() int { return len(al) }
func (al AsteroidList) Less(i, j int) bool { return al[i].AngleFromBase < al[j].AngleFromBase }
func (al AsteroidList) Swap(i, j int){ al[i], al[j] = al[j], al[i] }

type LaserGun struct {
    ScannedQuadrantData map[int]AsteroidList
    BaseStation         *Point
    HitAsteroids        []*Point
    HitAsteroidsMap     map[*Point]bool
    AsteroidAngles      map[*Point]float64
    AsteroidMap         *AsteroidMap
}

func (g *LaserGun) calculateAsteroidAngles() {
    g.AsteroidAngles = make(map[*Point]float64)
    baseVector := Vector{X: 0, Y: g.BaseStation.Y}

    for _, asteroid := range g.AsteroidMap.Asteroids {
        asteroidVector := Vector{X: asteroid.X - g.BaseStation.X, Y: asteroid.Y - g.BaseStation.Y}
        g.AsteroidAngles[asteroid] = baseVector.angleWith(asteroidVector)
    }
}

func (g *LaserGun) divideAsteroidsToQuadrants() {
    g.ScannedQuadrantData = make(map[int]AsteroidList)

    for _, asteroid := range g.AsteroidMap.Asteroids {
        quadrant := 0
        if asteroid.X >= g.BaseStation.X && asteroid.Y < g.BaseStation.Y {
            quadrant = 1
        }
        if asteroid.X > g.BaseStation.X && asteroid.Y >= g.BaseStation.Y {
            quadrant = 2
        }
        if asteroid.X <= g.BaseStation.X && asteroid.Y > g.BaseStation.Y {
            quadrant = 3
        }
        if asteroid.X < g.BaseStation.X && asteroid.Y <= g.BaseStation.Y {
            quadrant = 4
        }

        if quadrant > 0 {
            g.ScannedQuadrantData[quadrant] = append(g.ScannedQuadrantData[quadrant], Asteroid{Coordinates: asteroid, AngleFromBase: g.AsteroidAngles[asteroid]})
        }
    }

    // Sorting asteroids aby angle from baseline (horizontal line through Base Station).
    // For quadrants 2 (lower-right) and 4 (upper-left) we start with angles perpendicular to baseline.
    sort.Sort(g.ScannedQuadrantData[1])
    sort.Sort(g.ScannedQuadrantData[3])
    sort.Sort(sort.Reverse(g.ScannedQuadrantData[2]))
    sort.Sort(sort.Reverse(g.ScannedQuadrantData[4]))
}

func (g *LaserGun) shootQuadrant(quadrant int) {
    var vaporized []*Point
    for _, asteroid := range g.ScannedQuadrantData[quadrant] {
        if _, ok := g.HitAsteroidsMap[asteroid.Coordinates]; ok {
            continue
        }

        line := getLineFromTwoPoints(g.BaseStation, asteroid.Coordinates)

        blocked := false
        for _, blockingAsteroid := range g.AsteroidMap.Asteroids {
            if _, ok := g.HitAsteroidsMap[blockingAsteroid]; ok {
                continue
            }

            if line.containsPointInBetween(blockingAsteroid) {
                blocked = true
            }
        }

        if !blocked {
            vaporized = append(vaporized, asteroid.Coordinates)
        }
    }

    // We have to add vaporized asteroids to the final list only after the quadrant shooting is completed.
    // Otherwise we'd get false negatives for blocking asteroids.
    for _, asteroid := range vaporized {
        g.HitAsteroids = append(g.HitAsteroids, asteroid)
        g.HitAsteroidsMap[asteroid] = true
    }
}

type AsteroidMap struct {
    Asteroids []*Point
    Width     int
    Height    int
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
        m.Height = len(rows)
        m.Width = len(row)

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

// --------------------------------------------------------------------------------------------------------------------
// Generic structures

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

type Vector struct {
    X int
    Y int
}

func (v Vector) angleWith(ov Vector) float64 {
    return math.Acos(math.Abs(float64(v.X * ov.X + v.Y * ov.Y)) / (math.Sqrt(math.Pow(float64(v.X), 2) + math.Pow(float64(v.Y), 2)) * math.Sqrt(math.Pow(float64(ov.X), 2) + math.Pow(float64(ov.Y), 2))))
}
