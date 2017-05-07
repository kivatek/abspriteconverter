package main

import (
	"fmt"
	"flag"
	"os"
	"image"
	_ "image/png"
	_ "image/color"
	"math"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

func main() {
	if parseArguments() {
		buf, err := ioutil.ReadFile(configName)
		if err != nil {
			panic(err)
		}
		var config Config
		err = yaml.Unmarshal(buf, &config)
		if err != nil {
			panic(err)
		}
		process(config)
	} else {
		printUsage()
	}
}

var configName string
var rThreshold uint32 = 127
var gThreshold uint32 = 127
var bThreshold uint32 = 127

func parseArguments() bool {
	flag.Parse()
	if flag.NArg() == 1 {
		configName = flag.Arg(0)
		return true
	}
	return false
}

func printUsage() {
	fmt.Println("usage: abspriteconverter config.yml")
}

func process(config Config) {
	sb := make([]byte, 0)
	sb = append(sb, fmt.Sprintf("PROGMEM const byte %s[] = {\n", config.SpriteName)...)
	for index, source := range config.Images {
		s, err := perFileProcess(index, source.FileName)
		if err != nil {
			panic(err)
		}
		sb = append(sb, s...)
	}
	sb = append(sb, fmt.Sprintf("};\n")...)
	fmt.Printf("%s", string(sb))
}

func perFileProcess(index int, filename string) (string, error) {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return "", err
	}
	source, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	sb := make([]byte, 0)

	xSize := source.Bounds().Size().X
	ySize := source.Bounds().Size().Y
	if index == 0 {
		sb = append(sb, fmt.Sprintf("%d, %d,\n", xSize, ySize)...)
	}
	sb = append(sb, fmt.Sprintf("// %d: %s\n", index, filename)...)

	yLimit := (ySize + (8 - 1)) / 8
	for yCount := 0; yCount < yLimit; yCount++ {
		yRange := int(math.Min(float64(ySize - (yCount * 8)), 8))
		var count = 0
		for x := 0; x < xSize; x++ {
			hex := 0
			bit := 0x01
			for y := 0; y < yRange; y++ {
				pixel := source.At(x, y + (yCount * 8))
				r, g, b, a := pixel.RGBA()
				if a == 0 {
					r = 0
					g = 0
					b = 0
				}
				if r > rThreshold || g > gThreshold || b > bThreshold {
					hex |= bit
				}
				bit <<= 1
			}
			bit <<= uint32(8 - yRange)
			sb = append(sb, fmt.Sprintf("0x%02x,", hex)...)
			count++
			if (count % 8) == 0 {
				sb = append(sb, "\n"...)
			}
		}
		if (count % 8) != 0 {
			sb = append(sb, "\n"...)
		}
	}
	return string(sb), nil
}
