package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func main() {
main:
	for n := 0; n < 100; n++ {
		for v := 0; v < 100; v++ {
			sequence := loadSequence()
			returnCode := executeProgram(sequence, n, v)

			if returnCode == 19690720 {
				fmt.Println(fmt.Sprintf("Noun: %d, Verb: %d, Code: %d", n, v, 100*n+v))
				break main
			}
		}
	}
}

func loadSequence() []int {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	bytes, err := ioutil.ReadFile(path + "/2/code")

	if err != nil {
		fmt.Println(err)
	}

	inputs := strings.Split(string(bytes), ",")
	intInputs, err := convertStringArray(inputs)

	if err != nil {
		fmt.Println(err)
	}

	return intInputs
}

func executeProgram(sequence []int, noun, verb int) int {
	run := true
	seqPos := 0

	sequence[1] = noun
	sequence[2] = verb

	for run {
		switch sequence[seqPos] {
		case 1:
			add(sequence, seqPos)
		case 2:
			multiply(sequence, seqPos)
		case 99:
			run = false
		default:
			fmt.Println("Encountered invalid OpCode: ", sequence[seqPos])
			run = false
		}

		seqPos += 4
	}

	return sequence[0]
}

func add(a []int, pos int) {
	a[a[pos+3]] = a[a[pos+1]] + a[a[pos+2]]
}

func multiply(a []int, pos int) {
	a[a[pos+3]] = a[a[pos+1]] * a[a[pos+2]]
}

func convertStringArray(strArr []string) ([]int, error) {
	iArr := make([]int, 0, len(strArr))
	for _, str := range strArr {
		i, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		iArr = append(iArr, i)
	}
	return iArr, nil
}
