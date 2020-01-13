package main

import (
    "fmt"
    "io/ioutil"
    "strings"
)

const CenterOfMass = "COM"

func main() {
    orbitMap := OrbitMap{CenterOfMass: CenterOfMass}

    orbitMapData := loadOrbitMapData("C:/Users/tomas.bobek/AdventOfCode19/6/orbitMap")
    orbitMap.loadSpatialObjectList(orbitMapData)
    orbitMap.constructOrbitMap(orbitMapData)
    orbitMap.calculateOrbitDepthsFromCenter()

    fmt.Println("Number of direct and indirect orbits: ", orbitMap.getTotalNumberOfOrbits())
}

type OrbitMap struct {
    CenterOfMass string
    ObjectList   map[string]bool
    Map          map[string]*SpatialObject
}

type SpatialObject struct {
    Name       string
    OrbitDepth int
    Satellites []*SpatialObject
}

func loadOrbitMapData(file string) []string {
    bytes, err := ioutil.ReadFile(file)

    if err != nil {
        fmt.Println(err)
    }

    inputs := strings.Split(string(bytes), "\r\n")

    return inputs
}

func (o *OrbitMap) loadSpatialObjectList(data []string) {
    objectList := make(map[string]bool)

    for _, record := range data {
        pair := strings.Split(record, ")")

        if _, ok := objectList[pair[0]]; !ok {
            objectList[pair[0]] = true
        }

        if _, ok := objectList[pair[1]]; !ok {
            objectList[pair[1]] = true
        }
    }

    o.ObjectList = objectList
}

func (o *OrbitMap) constructOrbitMap(data []string) {
    orbitMap := make(map[string]*SpatialObject)

    for _, record := range data {
        pair := strings.Split(record, ")")

        var so *SpatialObject
        if val, ok := orbitMap[pair[1]]; ok {
            so = val
        } else {
            so = &SpatialObject{pair[1], 0, []*SpatialObject{}}
            orbitMap[pair[1]] = so
        }

        if val, ok := orbitMap[pair[0]]; ok {
            val.Satellites = append(val.Satellites, so)
        } else {
            orbitMap[pair[0]] = &SpatialObject{pair[0], 0, []*SpatialObject{so}}
        }
    }

    o.Map = orbitMap
}

func (o *OrbitMap) calculateOrbitDepthsFromCenter() {
    fillOrbitDepth(o.Map[o.CenterOfMass], 0)
}

func (o *OrbitMap) getTotalNumberOfOrbits() int {
    totalOrbits := 0
    for objectName := range o.ObjectList {
        totalOrbits += o.Map[objectName].OrbitDepth
    }

    return totalOrbits
}

func fillOrbitDepth(node *SpatialObject, depth int) {
    node.OrbitDepth = depth

    for _, satellite := range node.Satellites {
        fillOrbitDepth(satellite, depth + 1)
    }
}
