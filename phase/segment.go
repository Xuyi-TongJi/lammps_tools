package phase

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
)

const tStart = 295
const tEnd = 330
const deltaT = 5

const inputPrefix string = "/Users/xuyi/Desktop/phase515_loop/290/1/"
const inputSuffix string = ".txt"
const outputPath string = "/Users/xuyi/Desktop/phase515_loop/290/1/data/output.txt"

const titleLine = 14

type DataLine struct {
	timeStep int64
	rouV     float64
	tempNow  float64
	temp     float64
	volume   float64
	density  float64
	enthalpy float64
}

func Segment(average int, formatOutput func(*os.File, []DataLine, int), inputType string) {
	var output *os.File
	output, err := os.OpenFile(outputPath, os.O_RDWR, 0666)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		output, _ = os.Create(outputPath)
	}
	_ = output.Truncate(0)
	output.Seek(0, 0)
	var cnt int
	for t := tStart; t <= tEnd; t += deltaT {
		tStr := strconv.FormatInt(int64(t), 10)
		input := inputPrefix + tStr + inputType + inputSuffix
		dataFile, err := os.OpenFile(input, os.O_RDWR, 0666)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Reading data file: %s\n", input)
		var title string
		for i := 0; i < titleLine; i++ {
			_, _ = fmt.Fscan(dataFile, &title)
		}
		data := make([]DataLine, 0)
		for true {
			step := DataLine{}
			n, err := fmt.Fscanf(dataFile, "%d %f %f %f %f %f %f\n",
				&step.timeStep, &step.rouV, &step.tempNow, &step.temp, &step.volume, &step.density, &step.enthalpy,
			)
			if n == 0 || err != nil {
				break
			}
			data = append(data, step)
			cnt += 1
		}
		formatOutput(output, data, average)
	}
	fmt.Printf("Totally scan %d lines of data\n", cnt)
}

// TV
// temp - V
func TV(f *os.File, data []DataLine, average int) {
	sort.Slice(data, func(i, j int) bool {
		return data[i].temp < data[j].tempNow
	})
	curr, n := 0, len(data)
	sumT, sumV := 0.0, 0.0
	for i := 0; i < n; i++ {
		sumT += data[i].temp
		sumV += data[i].volume
		curr += 1
		if curr == average {
			fmt.Fprintf(f, "%.3f %.1f\n", sumT/float64(average), sumV/float64(average))
			curr = 0
			sumT, sumV = 0.0, 0.0
		}
	}
}
