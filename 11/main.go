package main

import (
    "bufio"
    "fmt"
    "io/ioutil"
    "math"
    "os"
    "strconv"
    "strings"
    "sync"
    "time"
)

type instructionOperation int

const (
    Add             instructionOperation = 1
    Multiply        instructionOperation = 2
    Read            instructionOperation = 3
    Write           instructionOperation = 4
    JumpIfTrue      instructionOperation = 5
    JumpIfFalse     instructionOperation = 6
    LessThan        instructionOperation = 7
    Equals          instructionOperation = 8
    SetRelativeBase instructionOperation = 9
    Terminate       instructionOperation = 99
)

type direction int

const (
    up direction = iota
    right
    down
    left
)

var (
    InstructionLength = map[instructionOperation]int{
        Add:4, Multiply:4, Read:2, Write:2, JumpIfTrue:3, JumpIfFalse:3, LessThan:4, Equals:4, SetRelativeBase:2, Terminate:1,
    }
)

func main() {
    path, err := os.Getwd()
    if err != nil {
        fmt.Println(err)
    }

    // -----------------------------------------------------------------------------------------------------------------
    // Here we solve problem for Part One
    robot := newPaintingRobotWithProgram(path + "/11/code")
    robot.run()

    fmt.Println("robot painted ", len(robot.paintedPoints), " tiles on the ship hull")
}

func newPaintingRobotWithProgram(programPath string) *paintingRobot {
    inChannel := make(chan int64)
    outChannel := make(chan int64)
    doneChannel := make(chan interface{})

    program := &program{
        inChannel: inChannel,
        outChannel: outChannel,
        done: doneChannel,
    }
    program.loadCodeFromFile(programPath)

    return &paintingRobot{
        brain: program,
        direction: up,
        position:  point{
            x:     0,
            y:     0,
        },
    }
}

type paintingRobot struct {
    brain         *program
    position      point
    direction     direction
    paintedPoints []point
}

func (r *paintingRobot) run() {
    wg := sync.WaitGroup{}
    wg.Add(1)
    go r.brain.execute()
    go func() {
        r.brain.inChannel <- 0
        readingColor := true

        robotLoop: for {
            var scannedColor int
            select {
            case reading := <-r.brain.outChannel:
                // Program outputs have 2 possible meanings that switch periodically:
                //  * color (0 - black, 1 - white)
                //  * rotation (0 - CCW, 1 - CW)
                if readingColor {
                    r.paint(int(reading))
                } else {
                    r.changeDirection(int(reading))
                    r.move()
                    scannedColor = r.scanColor()

                    // After orientation change the program expects the code of detected color on that position as input.
                    select {
                    case r.brain.inChannel <- int64(scannedColor):
                        fmt.Println("robot detected color ", scannedColor)
                    case <-r.brain.done:
                    }
                }

                readingColor = !readingColor

            case <-r.brain.done:
                wg.Done()
                break robotLoop
            }
        }
    }()
    wg.Wait()
}

// Gives the tile a color based on input (0 - black, 1 - white).
// In order to keep track of unique painted tiles we keep record in slice and just repaint existing items.
func (r *paintingRobot) paint(color int) {
    fmt.Println(fmt.Sprintf("robot paints [%d,%d] to color %d", r.position.x, r.position.y, color))

    for _, p := range r.paintedPoints {
        if p.x == r.position.x && p.y == r.position.y {
            p.color = color
            fmt.Println("just repainted, # of painted tiles: ", len(r.paintedPoints))
            return
        }
    }

    r.position.color = color
    r.paintedPoints = append(r.paintedPoints, r.position)
    fmt.Println("NEW painting, # of painted tiles: ", len(r.paintedPoints))
}

// Rotates the direction robot is facing - 0 for CW rotation and 1 for CCW rotation.
func (r *paintingRobot) changeDirection(input int) {
    if input == 1 {
        if r.direction == up {
            r.direction = left
        } else {
            r.direction -= 1
        }
    } else {
        if r.direction == left {
            r.direction = up
        } else {
            r.direction += 1
        }
    }
}

// Moves the robot by 1 distance point in the direction it is currently facing.
func (r *paintingRobot) move() {
    posX, posY := r.position.x, r.position.y
    switch r.direction {
    case up:
        posY -= 1
    case right:
        posX += 1
    case down:
        posY += 1
    case left:
        posX -= 1
    }

    r.position = point{
        x:     posX,
        y:     posY,
    }

    fmt.Println(fmt.Sprintf("robot moved to [%d,%d]", r.position.x, r.position.y))
}

// Gets the color of underlying tile (based on robot's position). Default color is black (0).
func (r *paintingRobot) scanColor() int {
    for _, p := range r.paintedPoints {
        if p.x == r.position.x && p.y == r.position.y {
            return p.color
        }
    }

    return 0
}

type point struct {
    x     int
    y     int
    color int
}

type instruction struct {
    operation instructionOperation
    length int
    params []instructionParam
}

func (i *instruction) initialize(intCode []int64, pIndex int) {
    instValue := int(intCode[pIndex])

    i.operation = instructionOperation(instValue)

    // Standard operation Codes are between 1 and 99, larger number means that Parameter Modes are included there
    evalParamModes := false
    if instValue >= 100 {
        i.operation = instructionOperation(instValue % 100)
        evalParamModes = true
    }

    i.length = InstructionLength[i.operation]
    paramCount := i.length - 1
    i.params = make([]instructionParam, paramCount, paramCount)

    for j := 0; j < paramCount; j++ {
        i.params[j] = instructionParam{0, intCode[pIndex+j+1]}

        // Parameter mode is either 0 (by reference) or 1 (by value) and this mode is specified
        // in the instruction code itself (as given number at respective position)
        if evalParamModes {
            i.params[j].mode = (instValue / int(math.Pow(float64(10), float64(j+2)))) % 10
        }
    }
}

func (i *instruction) getValuesCount() int {
    switch i.operation {
    case Add, Multiply, JumpIfTrue, JumpIfFalse, LessThan, Equals:
        return 2
    case Write, SetRelativeBase:
        return 1
    default:
        return 0
    }
}

func (i *instruction) doesStoreOutputInMemory() bool {
    return i.operation == Read || i.operation == Add || i.operation == Multiply || i.operation == LessThan || i.operation == Equals
}

type instructionParam struct {
    mode  int
    value int64
}

type program struct {
    memory       []int64
    memorySize   int
    position     int
    relativeBase int
    completed    bool
    halt         bool

    inChannel    chan int64
    outChannel   chan int64
    done         chan interface{}

    dataStack    []int64
    haltOnOutput bool
}

func (p *program) loadCodeFromFile(file string) {
    bytes, err := ioutil.ReadFile(file)

    if err != nil {
        fmt.Println(err)
    }

    inputs := strings.Split(string(bytes), ",")
    intInputs, err := convertStringArray(inputs)

    if err != nil {
        fmt.Println(err)
    }

    p.memorySize = len(intInputs) * 10
    p.memory = make([]int64, p.memorySize, p.memorySize)
    for i := 0; i < len(intInputs); i++ {
        p.memory[i] = intInputs[i]
    }
}

func (p *program) resetState() {
    p.position = 0
    p.completed = false
    p.halt = false
}

func (p *program) resetMemory() {
    p.dataStack = make([]int64, p.memorySize, p.memorySize)
}

func (p *program) execute() {
    for !p.completed && !p.halt {
        var instruction instruction
        instruction.initialize(p.memory, p.position)

        p.loadParameterValues(&instruction)

        switch instruction.operation {
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
            fmt.Println("program finished")
            p.completed = true
            close(p.done)
            close(p.outChannel)
        default:
            fmt.Println("Encountered invalid OpCode: ", instruction.operation)
            p.completed = true
            close(p.done)
            close(p.outChannel)
        }
    }
}

// Parameters can be handled "by value" or "by reference" and this function supplies the end value in each case
func (p *program) loadParameterValues(i *instruction) {
    for j := 0; j < i.getValuesCount(); j++ {
        switch i.params[j].mode {
        case 0:
            i.params[j].value = p.memory[i.params[j].value]
        case 2:
            i.params[j].value = p.memory[p.relativeBase+ int(i.params[j].value)]
        }
    }

    if i.doesStoreOutputInMemory() {
        if i.params[i.getValuesCount()].mode == 2 {
            i.params[i.getValuesCount()].value = int64(p.relativeBase) + i.params[i.getValuesCount()].value
        }
    }
}

func (p *program) doAdd(i *instruction) {
    p.memory[i.params[2].value] = i.params[0].value + i.params[1].value
    p.position += i.length
}

func (p *program) doMultiply(i *instruction) {
    p.memory[i.params[2].value] = i.params[0].value * i.params[1].value
    p.position += i.length
}

// Inputs are primarily read from dataStack of the program, if it is empty, input is prompted from Standard Input
func (p *program) doReadInput(i *instruction) {
    var input int64
    channelReadOk := false

    if p.inChannel != nil {
        select {
        case <-time.After(10 * time.Second):
            fmt.Println("waiting for input timed-out, trying to read from dataStack")
        case input = <-p.inChannel:
            channelReadOk = true
        }
    }

    if !channelReadOk {
        if len(p.dataStack) > 0 {
            input = p.dataStack[len(p.dataStack)-1]
            p.dataStack = p.dataStack[:len(p.dataStack)-1]
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
    }

    p.memory[i.params[0].value] = input
    p.position += i.length
}

// program outputs are logged to Standard Output and stored in internal Data Stack
func (p *program) doWriteOutput(i *instruction) {
    if p.outChannel != nil {
        p.outChannel <- i.params[0].value
    } else {
        p.dataStack = append(p.dataStack, i.params[0].value)
    }
    p.position += i.length

    if p.haltOnOutput {
        p.halt = true
    }
}

func (p *program) doJumpIfTrue(i *instruction) {
    if i.params[0].value != 0 {
        p.position = int(i.params[1].value)
    } else {
        p.position += i.length
    }
}

func (p *program) doJumpIfFalse(i *instruction) {
    if i.params[0].value == 0 {
        p.position = int(i.params[1].value)
    } else {
        p.position += i.length
    }
}

func (p *program) doComparisonLessThan(i *instruction) {
    if i.params[0].value < i.params[1].value {
        p.memory[i.params[2].value] = 1
    } else {
        p.memory[i.params[2].value] = 0
    }
    p.position += i.length
}

func (p *program) doComparisonEquals(i *instruction) {
    if i.params[0].value == i.params[1].value {
        p.memory[i.params[2].value] = 1
    } else {
        p.memory[i.params[2].value] = 0
    }
    p.position += i.length
}

func (p *program) doUpdateRelativeBase(i *instruction) {
    p.relativeBase += int(i.params[0].value)
    p.position += i.length
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
