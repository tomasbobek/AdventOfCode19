package main

import (
	"fmt"
	"strconv"
)

const (
	RangeMin = 171309
	RangeMax = 643603
)

func main() {
	totalMatches := 0

	for i := RangeMin; i <= RangeMax; i++ {
		literal := strconv.Itoa(i)

		if checkCriteria(literal) {
			totalMatches++
		}
	}

	fmt.Println("Total eligible passwords: ", totalMatches)
}

func checkCriteria(literal string) bool {
	containsPair := false
	areDigitsIncreasing := true

	var lastChar int32
	for index, char := range literal {
		if char == lastChar && !containsPair {
			containsPair = true

			if index > 1 && char == int32(literal[index-2]) || index < 5 && char == int32(literal[index+1]) {
				containsPair = false
			}
		}

		if char < lastChar {
			areDigitsIncreasing = false
		}

		lastChar = char
	}

	return containsPair && areDigitsIncreasing
}
