package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
)

var (
	InstructionLenght = map[int]int{
		1: 4, 2: 4, 3: 2, 4: 2, 5: 3, 6: 3, 7: 4, 8: 4, 99: 1,
	}
)

func main() {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	program := Program{Position: 0, Completed: false}

	program.loadCodeFromFile(path + "/5/code")
	program.execute()
}

type Instruction struct {
	OpCode int
	Length int
	Params []InstructionParam
}

func (i *Instruction) initialize(intCode []int, pIndex int) {
	instValue := intCode[pIndex]

	i.OpCode = instValue

	evalParamModes := false
	if instValue >= 100 {
		i.OpCode = instValue % 100
		evalParamModes = true
	}

	i.Length = InstructionLenght[i.OpCode]
	paramCount := i.Length - 1
	i.Params = make([]InstructionParam, paramCount, paramCount)

	for j := 0; j < paramCount; j++ {
		i.Params[j] = InstructionParam{0, intCode[pIndex+j+1]}

		if evalParamModes {
			i.Params[j].Mode = (instValue / int(math.Pow(float64(10), float64(j+2)))) % 10
		}
	}
}

func (i *Instruction) getValuesCount() int {
	switch i.OpCode {
	case 1, 2, 5, 6, 7, 8:
		return 2
	case 4:
		return 1
	default:
		return 0
	}
}

type InstructionParam struct {
	Mode  int
	Value int
}

type Program struct {
	IntCode   []int
	Position  int
	Completed bool
}

func (p *Program) loadCodeFromFile(file string) {
	bytes, err := ioutil.ReadFile(file)

	if err != nil {
		fmt.Println(err)
	}

	inputs := strings.Split(string(bytes), ",")
	intInputs, err := convertStringArray(inputs)

	if err != nil {
		fmt.Println(err)
	}

	p.IntCode = intInputs
}

func (p *Program) execute() {
	for !p.Completed {
		var instruction Instruction
		instruction.initialize(p.IntCode, p.Position)

		p.loadParameterValues(&instruction)

		switch instruction.OpCode {
		case 1:
			p.doAdd(&instruction)
		case 2:
			p.doMultiply(&instruction)
		case 3:
			p.doReadInput(&instruction)
		case 4:
			p.doWriteOutput(&instruction)
		case 5:
			p.doJumpIfTrue(&instruction)
		case 6:
			p.doJumpIfFalse(&instruction)
		case 7:
			p.doComparisonLessThan(&instruction)
		case 8:
			p.doComparisonEquals(&instruction)
		case 99:
			fmt.Println("Program finished")
			p.Completed = true
		default:
			fmt.Println("Encountered invalid OpCode: ", instruction.OpCode)
			p.Completed = true
		}
	}
}

func (p *Program) loadParameterValues(i *Instruction) {
	for j := 0; j < i.getValuesCount(); j++ {
		if i.Params[j].Mode == 0 {
			i.Params[j].Value = p.IntCode[i.Params[j].Value]
		}
	}
}

func (p *Program) doAdd(i *Instruction) {
	p.IntCode[i.Params[2].Value] = i.Params[0].Value + i.Params[1].Value
	p.Position += i.Length
}

func (p *Program) doMultiply(i *Instruction) {
	p.IntCode[i.Params[2].Value] = i.Params[0].Value * i.Params[1].Value
	p.Position += i.Length
}

func (p *Program) doReadInput(i *Instruction) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter value: ")
	value, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println(err)
	}

	intValue, err := strconv.Atoi(strings.TrimSuffix(value, "\n"))

	if err != nil {
		fmt.Println(err)
	}

	p.IntCode[i.Params[0].Value] = intValue
	p.Position += i.Length
}

func (p *Program) doWriteOutput(i *Instruction) {
	fmt.Println("Program outputs: ", i.Params[0].Value)
	p.Position += i.Length
}

func (p *Program) doJumpIfTrue(i *Instruction) {
	if i.Params[0].Value != 0 {
		p.Position = i.Params[1].Value
	} else {
		p.Position += i.Length
	}
}

func (p *Program) doJumpIfFalse(i *Instruction) {
	if i.Params[0].Value == 0 {
		p.Position = i.Params[1].Value
	} else {
		p.Position += i.Length
	}
}

func (p *Program) doComparisonLessThan(i *Instruction) {
	if i.Params[0].Value < i.Params[1].Value {
		p.IntCode[i.Params[2].Value] = 1
	} else {
		p.IntCode[i.Params[2].Value] = 0
	}
	p.Position += i.Length
}

func (p *Program) doComparisonEquals(i *Instruction) {
	if i.Params[0].Value == i.Params[1].Value {
		p.IntCode[i.Params[2].Value] = 1
	} else {
		p.IntCode[i.Params[2].Value] = 0
	}
	p.Position += i.Length
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
