package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	bytes, err := ioutil.ReadFile(path + "/1/input")

	if err != nil {
		fmt.Println(err)
	}

	inputs := strings.Split(string(bytes), "\r\n")
	totalFuel := 0

	for _, moduleMass := range inputs {
		mass, err := strconv.Atoi(moduleMass)

		if err != nil {
			fmt.Println(err)
		}

		moduleFuel := calculateFuelReq(mass)
		totalFuel += moduleFuel

		for moduleFuel > 0 {
			moduleFuel = calculateFuelReq(moduleFuel)

			if moduleFuel > 0 {
				totalFuel += moduleFuel
			}
		}
	}

	fmt.Println(totalFuel)
}

func calculateFuelReq(mass int) int {
	fuelReq := math.Floor(float64(mass/3)) - 2
	return int(fuelReq)
}
