package phase

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

const TStart = 295
const TEnd = 330
const DeltaT = 5

const inputPrefix string = "/Users/xuyi/Desktop/phase515_loop/290/1/"
const inputSuffix string = ".txt"
const outputPrefix = inputPrefix + "data/"
const outputSuffix string = ".txt"

var titleLine map[string]int = make(map[string]int, 0)

type DataLine struct {
	timeStep int64
	rouV     float64
	tempNow  float64
	temp     float64
	volume   float64
	density  float64
	enthalpy float64
}

func init() {
	titleLine[".temp"] = 14
}

// Integration 分段数据整合处理
func Integration(average int, formatOutput func(*os.File, []DataLine, int), inputType string, outputStr string, wg *sync.WaitGroup) {
	outputPath := outputPrefix + outputStr + outputSuffix
	writer := openOutput(average, outputPath, formatOutput, wg)
	cnt := dataProcess(TStart, TEnd, DeltaT, inputType, writer)
	fmt.Printf("Totally scan %d lines of data\n", cnt)
}

// Single 分段数据分段处理
func Single(target int, average int, formatOutput func(*os.File, []DataLine, int), inputType string, outputStr string, wg *sync.WaitGroup) {
	targetStr := strconv.FormatInt(int64(target), 10)
	outputStr += targetStr
	outputPath := outputPrefix + outputStr + outputSuffix
	writer := openOutput(average, outputPath, formatOutput, wg)
	cnt := dataProcess(target, target, 0, inputType, writer)
	fmt.Printf("Totally scan %d lines of data\n", cnt)
}

func openOutput(average int, outputPath string, formatOutput func(*os.File, []DataLine, int), wg *sync.WaitGroup) (writer chan []DataLine) {
	output, err := os.OpenFile(outputPath, os.O_RDWR, 0666)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		output, _ = os.Create(outputPath)
	}
	_ = output.Truncate(0)
	_, _ = output.Seek(0, 0)
	writer = make(chan []DataLine)
	go func() {
		for true {
			select {
			case data := <-writer:
				{
					// output
					formatOutput(output, data, average)
					wg.Done()
				}
			}
		}
	}()
	time.Sleep(1 * time.Second)
	return
}

func dataProcess(start, end, delta int, inputType string, writer chan []DataLine) int {
	var cnt = 0
	data := make([]DataLine, 0)
	for t := start; t <= end; t += delta {
		tStr := strconv.FormatInt(int64(t), 10)
		input := inputPrefix + tStr + inputType + inputSuffix
		dataFile, err := os.OpenFile(input, os.O_RDWR, 0666)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Reading data file: %s\n", input)
		var title string
		for i := 0; i < titleLine[inputType]; i++ {
			_, _ = fmt.Fscan(dataFile, &title)
		}
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
	}
	writer <- data
	return cnt
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

// HT
// enthalpy - temp - rouV
// 计算比热
func HT(f *os.File, data []DataLine, average int) {
	sort.Slice(data, func(i, j int) bool {
		return data[i].temp < data[j].tempNow
	})
	curr, n := 0, len(data)
	sumT, sumH, sumD := 0.0, 0.0, 0.0
	fmt.Fprintln(f, "Temp-Enthalpy-Density")
	for i := 0; i < n; i++ {
		sumT += data[i].temp
		sumH += data[i].enthalpy
		sumD += data[i].rouV
		curr += 1
		if curr == average {
			fmt.Fprintf(f, "%.3f %.2f %.1f\n", sumT/float64(average), sumH/float64(average), sumD/float64(average))
			curr = 0
			sumT, sumH, sumD = 0.0, 0.0, 0.0
		}
	}
}
