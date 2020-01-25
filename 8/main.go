package main

import (
    "fmt"
    "image"
    "image/color"
    "image/png"
    "io/ioutil"
    "math"
    "os"
    "strconv"
)

func main() {
    path, err := os.Getwd()
    if err != nil {
        fmt.Println(err)
    }

    // -----------------------------------------------------------------------------------------------------------------
    // Here we solve problem for Part One
    processor := ImageProcessor{ImageWidth: 25, ImageHeight: 6}
    processor.loadDataFromFile(path + "/8/imageData")
    processor.constructLayers()

    minZeroCount := math.MaxInt32
    var leastZeroLayer *ImageLayer
    for _, layer := range processor.Layers {
        zeroCount := layer.getNumberOf(0)
        if zeroCount < minZeroCount {
            minZeroCount = zeroCount
            leastZeroLayer = layer
        }
    }

    if leastZeroLayer != nil {
        fmt.Println("Checksum of layer with least amount of zeros: ", leastZeroLayer.calculateChecksum())
    } else {
        fmt.Println("Error: Could not locate layer with least amount of zeros!")
    }

    // -----------------------------------------------------------------------------------------------------------------
    // Here we solve problem for Part Two
    processor.processImage()
    processor.renderImage(path + "/8/elvenImage.png")
}

type ImageProcessor struct {
    ImageWidth  int
    ImageHeight int

    RawData     []int
    ImageData   []int

    Layers      []*ImageLayer
}

func (ip *ImageProcessor) loadDataFromFile(file string) {
    bytes, err := ioutil.ReadFile(file)

    if err != nil {
        fmt.Println(err)
    }

    for _, val := range bytes {
        intVal, err := strconv.Atoi(string(val))

        if err != nil {
            fmt.Println(err)
        }

        ip.RawData = append(ip.RawData, intVal)
    }
}

func (ip *ImageProcessor) constructLayers() {
    layerLength := ip.ImageWidth * ip.ImageHeight
    totalLayers := len(ip.RawData) / layerLength

    for i := 0; i < totalLayers; i++ {
        var layer ImageLayer
        layer.RawData = ip.RawData[(i * layerLength):((i + 1) * layerLength)]
        ip.Layers = append(ip.Layers, &layer)
    }
}

func (ip *ImageProcessor) processImage() {
    ip.resetImageData()

    for _, layer := range ip.Layers {
        for i := 0; i < len(layer.RawData); i++ {
            if ip.ImageData[i] == -1 || ip.ImageData[i] == 2 {
                ip.ImageData[i] = layer.RawData[i]
            }
        }
    }
}

func (ip *ImageProcessor) renderImage(output string) {
    startPoint := image.Point{0, 0}
    endPoint := image.Point{ip.ImageWidth, ip.ImageHeight}

    img := image.NewRGBA(image.Rectangle{startPoint, endPoint})

    for x := 0; x < ip.ImageWidth; x++ {
        for y := 0; y < ip.ImageHeight; y++ {
            img.Set(x, y, ip.getPixelColor(x, y))
        }
    }

    f, err := os.Create(output)
    if err != nil {
        fmt.Println(err)
    }

    err = png.Encode(f, img)
    if err != nil {
        fmt.Println(err)
    }
}

func (ip *ImageProcessor) getPixelColor(x, y int) color.RGBA {
    colorCode := ip.ImageData[y * ip.ImageWidth + x]

    switch colorCode {
    case 0:
        return color.RGBA{0, 0, 0, 0xff}
    case 1:
        return color.RGBA{255, 255, 255, 0xff}
    default:
        return color.RGBA{0, 0, 0, 0x00}
    }
}

func (ip *ImageProcessor) resetImageData() {
    dataLength := ip.ImageWidth * ip.ImageHeight
    imageData := make([]int, dataLength, dataLength)

    for i := 0; i < dataLength; i++ {
        imageData[i] = -1
    }

    ip.ImageData = imageData
}

type ImageLayer struct {
    RawData     []int
}

func (il *ImageLayer) getNumberOf(n int) int {
    counter := 0
    for _, num := range il.RawData {
        if num == n {
            counter++
        }
    }

    return counter
}

func (il *ImageLayer) calculateChecksum() int {
    return il.getNumberOf(1) * il.getNumberOf(2)
}
