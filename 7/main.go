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
    InstructionLength = map[int]int{
        1: 4, 2: 4, 3: 2, 4: 2, 5: 3, 6: 3, 7: 4, 8: 4, 99: 1,
    }
)

func main() {
    // -----------------------------------------------------------------------------------------------------------------
    // Here we solve problem for Part One
    program := Program{Position: 0, Completed: false}
    program.loadCodeFromFile("C:/Users/tomas.bobek/AdventOfCode19/7/code")

    bestSignal := 0
    for _, thrusterConfig := range getArrayPermutations([]int{0,1,2,3,4}) {
        signal := program.launchThrusterSequence(thrusterConfig)

        if signal > bestSignal {
            bestSignal = signal
        }
    }

    fmt.Println("Sequence that generates max power: ", bestSignal)

    // -----------------------------------------------------------------------------------------------------------------
    // Here we solve problem for Part Two (feedback loop)
    amplifierA := &Program{Position: 0, Completed: false, HaltOnOutput: true}
    amplifierA.loadCodeFromFile("C:/Users/tomas.bobek/AdventOfCode19/7/code")
    amplifierB := &Program{Position: 0, Completed: false, HaltOnOutput: true}
    amplifierB.loadCodeFromFile("C:/Users/tomas.bobek/AdventOfCode19/7/code")
    amplifierC := &Program{Position: 0, Completed: false, HaltOnOutput: true}
    amplifierC.loadCodeFromFile("C:/Users/tomas.bobek/AdventOfCode19/7/code")
    amplifierD := &Program{Position: 0, Completed: false, HaltOnOutput: true}
    amplifierD.loadCodeFromFile("C:/Users/tomas.bobek/AdventOfCode19/7/code")
    amplifierE := &Program{Position: 0, Completed: false, HaltOnOutput: true}
    amplifierE.loadCodeFromFile("C:/Users/tomas.bobek/AdventOfCode19/7/code")
    amplifiers := []*Program{amplifierA, amplifierB, amplifierC, amplifierD, amplifierE}

    bestSignal = 0
    for _, thrusterConfig := range getArrayPermutations([]int{5,6,7,8,9}) {

        // Reset memory and state of each Amplifier for new set of input data
        for _, amplifier := range amplifiers {
            amplifier.resetMemory()
            amplifier.resetState()
        }

        amplifierPosition := 0
        firstRound := true
        for !amplifiers[len(amplifiers) - 1].Completed {

            // Amplifiers consume input sequence during first loop
            if firstRound {
                if amplifierPosition == 0 {
                    amplifiers[0].DataStack = append(amplifiers[0].DataStack, 0)
                }
                amplifiers[amplifierPosition].DataStack = append(amplifiers[amplifierPosition].DataStack, thrusterConfig[amplifierPosition])
            }

            // Return back to first Amplifier once reached the last one
            nextPosition := amplifierPosition + 1
            if amplifierPosition == len(amplifiers) - 1 {
                nextPosition = 0
                firstRound = false
            }

            // Amplifiers' program is configured to halt on output, then we can unblock it and pass its last output as input to the next Amplifier
            amplifiers[amplifierPosition].execute()
            amplifiers[amplifierPosition].Halt = false
            amplifiers[nextPosition].DataStack = append(amplifiers[nextPosition].DataStack, amplifiers[amplifierPosition].DataStack[0])
            amplifiers[amplifierPosition].DataStack = amplifiers[amplifierPosition].DataStack[:len(amplifiers[amplifierPosition].DataStack)-1]

            amplifierPosition = nextPosition
        }

        // When last Amplifier completes, it passes its last output to the next Amplifier, which is the first one
        signal := amplifiers[0].DataStack[0]

        if signal > bestSignal {
            bestSignal = signal
        }
    }

    fmt.Println("Feedback loop sequence that generates max power: ", bestSignal)
}

type Instruction struct {
    OpCode int
    Length int
    Params []InstructionParam
}

func (i *Instruction) initialize(intCode []int, pIndex int) {
    instValue := intCode[pIndex]

    i.OpCode = instValue

    // Standard Operation Codes are between 1 and 99, larger number means that Parameter Modes are included there
    evalParamModes := false
    if instValue >= 100 {
        i.OpCode = instValue % 100
        evalParamModes = true
    }

    i.Length = InstructionLength[i.OpCode]
    paramCount := i.Length - 1
    i.Params = make([]InstructionParam, paramCount, paramCount)

    for j := 0; j < paramCount; j++ {
        i.Params[j] = InstructionParam{0, intCode[pIndex+j+1]}

        // Parameter Mode is either 0 (by reference) or 1 (by value) and this mode is specified
        // in the Instruction code itself (as given number at respective position)
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
    IntCode      []int
    Position     int
    Completed    bool
    Halt         bool

    DataStack    []int
    HaltOnOutput bool
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

// Runs the program as many time as there is number of inputs in sequence
func (p *Program) launchThrusterSequence(sequence []int) int {
    p.resetMemory()
    p.DataStack = append(p.DataStack, 0)

    for _, input := range sequence {
        p.resetState()
        p.DataStack = append(p.DataStack, input)
        p.execute()
    }

    return p.DataStack[0]
}

func (p *Program) resetState() {
    p.Position = 0
    p.Completed = false
    p.Halt = false
}

func (p *Program) resetMemory() {
    p.DataStack = []int{}
}

func (p *Program) execute() {
    for !p.Completed && !p.Halt {
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

// Parameters can be handled "by value" or "by reference" and this function supplies the end value in each case
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

// Inputs are primarily read from DataStack of the Program, if it is empty, input is prompted from Standard Input
func (p *Program) doReadInput(i *Instruction) {
    var input int

     if len(p.DataStack) > 0 {
        input = p.DataStack[len(p.DataStack)-1]
        p.DataStack = p.DataStack[:len(p.DataStack)-1]
    } else {
        reader := bufio.NewReader(os.Stdin)
        fmt.Print("Enter value: ")
        value, err := reader.ReadString('\n')

        if err != nil {
            fmt.Println(err)
        }

        input, err = strconv.Atoi(strings.TrimSuffix(value, "\n"))

        if err != nil {
            fmt.Println(err)
        }
    }

    p.IntCode[i.Params[0].Value] = input
    p.Position += i.Length
}

// Program outputs are logged to Standard Output and stored in internal Data Stack
func (p *Program) doWriteOutput(i *Instruction) {
    fmt.Println("Program outputs: ", i.Params[0].Value)
    p.DataStack = append(p.DataStack, i.Params[0].Value)
    p.Position += i.Length

    if p.HaltOnOutput {
        p.Halt = true
    }
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

func getArrayPermutations(inputs []int) [][]int {
    var permutations [][]int
    var getPerm func([]int, int)
    getPerm = func(a []int, k int) {
        if k == len(a) {
            permutations = append(permutations, append([]int{}, a...))
        } else {
            for i := k; i < len(inputs); i++ {
                a[k], a[i] = a[i], a[k]
                getPerm(a, k+1)
                a[k], a[i] = a[i], a[k]
            }
        }
    }
    getPerm(inputs, 0)

    return permutations
}
