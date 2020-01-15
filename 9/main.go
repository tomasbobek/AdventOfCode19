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

type InstructionOperation int

const (
    Add             InstructionOperation = 1
    Multiply        InstructionOperation = 2
    Read            InstructionOperation = 3
    Write           InstructionOperation = 4
    JumpIfTrue      InstructionOperation = 5
    JumpIfFalse     InstructionOperation = 6
    LessThan        InstructionOperation = 7
    Equals          InstructionOperation = 8
    SetRelativeBase InstructionOperation = 9
    Terminate       InstructionOperation = 99
)

var (
    InstructionLength = map[InstructionOperation]int{
        Add:4, Multiply:4, Read:2, Write:2, JumpIfTrue:3, JumpIfFalse:3, LessThan:4, Equals:4, SetRelativeBase:2, Terminate:1,
    }
)

func main() {
    // -----------------------------------------------------------------------------------------------------------------
    // Here we solve problem for Part One
    program := Program{}
    program.loadCodeFromFile("C:/Users/tomas.bobek/AdventOfCode19/9/code")
    program.execute()

    fmt.Println("Program generated following BOOST code: ", program.DataStack[len(program.DataStack) - 1])
}

type Instruction struct {
    Operation InstructionOperation
    Length    int
    Params    []InstructionParam
}

func (i *Instruction) initialize(intCode []int64, pIndex int) {
    instValue := int(intCode[pIndex])

    i.Operation = InstructionOperation(instValue)

    // Standard Operation Codes are between 1 and 99, larger number means that Parameter Modes are included there
    evalParamModes := false
    if instValue >= 100 {
        i.Operation = InstructionOperation(instValue % 100)
        evalParamModes = true
    }

    i.Length = InstructionLength[i.Operation]
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
    switch i.Operation {
    case Add, Multiply, JumpIfTrue, JumpIfFalse, LessThan, Equals:
        return 2
    case Write, SetRelativeBase:
        return 1
    default:
        return 0
    }
}

func (i *Instruction) doesStoreOutputInMemory() bool {
    return i.Operation == Read || i.Operation == Add || i.Operation == Multiply || i.Operation == LessThan || i.Operation == Equals
}

type InstructionParam struct {
    Mode  int
    Value int64
}

type Program struct {
    Memory       []int64
    MemorySize   int
    Position     int
    RelativeBase int
    Completed    bool
    Halt         bool

    DataStack    []int64
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

    p.MemorySize = len(intInputs) * 10
    p.Memory = make([]int64, p.MemorySize, p.MemorySize)
    for i := 0; i < len(intInputs); i++ {
        p.Memory[i] = intInputs[i]
    }
}

func (p *Program) resetState() {
    p.Position = 0
    p.Completed = false
    p.Halt = false
}

func (p *Program) resetMemory() {
    p.DataStack = make([]int64, p.MemorySize, p.MemorySize)
}

func (p *Program) execute() {
    for !p.Completed && !p.Halt {
        var instruction Instruction
        instruction.initialize(p.Memory, p.Position)

        p.loadParameterValues(&instruction)

        switch instruction.Operation {
        case Add:
            p.doAdd(&instruction)
        case Multiply:
            p.doMultiply(&instruction)
        case Read:
            p.doReadInput(&instruction)
        case Write:
            p.doWriteOutput(&instruction)
        case JumpIfTrue:
            p.doJumpIfTrue(&instruction)
        case JumpIfFalse:
            p.doJumpIfFalse(&instruction)
        case LessThan:
            p.doComparisonLessThan(&instruction)
        case Equals:
            p.doComparisonEquals(&instruction)
        case SetRelativeBase:
            p.doUpdateRelativeBase(&instruction)
        case Terminate:
            fmt.Println("Program finished")
            p.Completed = true
        default:
            fmt.Println("Encountered invalid OpCode: ", instruction.Operation)
            p.Completed = true
        }
    }
}

// Parameters can be handled "by value" or "by reference" and this function supplies the end value in each case
func (p *Program) loadParameterValues(i *Instruction) {
    for j := 0; j < i.getValuesCount(); j++ {
        switch i.Params[j].Mode {
        case 0:
            i.Params[j].Value = p.Memory[i.Params[j].Value]
        case 2:
            i.Params[j].Value = p.Memory[p.RelativeBase + int(i.Params[j].Value)]
        }
    }

    if i.doesStoreOutputInMemory() {
        if i.Params[i.getValuesCount()].Mode == 2 {
            i.Params[i.getValuesCount()].Value = int64(p.RelativeBase) + i.Params[i.getValuesCount()].Value
        }
    }
}

func (p *Program) doAdd(i *Instruction) {
    p.Memory[i.Params[2].Value] = i.Params[0].Value + i.Params[1].Value
    p.Position += i.Length
}

func (p *Program) doMultiply(i *Instruction) {
    p.Memory[i.Params[2].Value] = i.Params[0].Value * i.Params[1].Value
    p.Position += i.Length
}

// Inputs are primarily read from DataStack of the Program, if it is empty, input is prompted from Standard Input
func (p *Program) doReadInput(i *Instruction) {
    var input int64

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

        inputInt, err := strconv.Atoi(strings.TrimSuffix(value, "\n"))

        if err != nil {
            fmt.Println(err)
        }

        input = int64(inputInt)
    }

    p.Memory[i.Params[0].Value] = input
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
        p.Position = int(i.Params[1].Value)
    } else {
        p.Position += i.Length
    }
}

func (p *Program) doJumpIfFalse(i *Instruction) {
    if i.Params[0].Value == 0 {
        p.Position = int(i.Params[1].Value)
    } else {
        p.Position += i.Length
    }
}

func (p *Program) doComparisonLessThan(i *Instruction) {
    if i.Params[0].Value < i.Params[1].Value {
        p.Memory[i.Params[2].Value] = 1
    } else {
        p.Memory[i.Params[2].Value] = 0
    }
    p.Position += i.Length
}

func (p *Program) doComparisonEquals(i *Instruction) {
    if i.Params[0].Value == i.Params[1].Value {
        p.Memory[i.Params[2].Value] = 1
    } else {
        p.Memory[i.Params[2].Value] = 0
    }
    p.Position += i.Length
}

func (p *Program) doUpdateRelativeBase(i *Instruction) {
    p.RelativeBase += int(i.Params[0].Value)
    p.Position += i.Length
}

func convertStringArray(strArr []string) ([]int64, error) {
    iArr := make([]int64, 0, len(strArr))
    for _, str := range strArr {
        i, err := strconv.Atoi(str)
        if err != nil {
            return nil, err
        }
        iArr = append(iArr, int64(i))
    }
    return iArr, nil
}
