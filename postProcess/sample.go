package postProcess

import (
	. "atom/util"
	"bufio"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
	"sort"
	"strconv"
	"sync"
)

const (
	XLSX         = ".xlsx"
	TXT          = ".txt"
	TVDiagram    = "TV"
	DefaultSheet = "Sheet1"
	HEAT         = "heat"
	ENERGY       = "energy"
)

var titleLine = make(map[string]int, 0)
var formatMap = make(map[string]func(any, []Data, int, int, int, string))

func init() {
	titleLine[HEAT] = 13
	titleLine[ENERGY] = 13
	formatMap[TVDiagram] = TV
}

func OpenInput(input Input, writers []chan []Data) {
	data := make([]Data, 0)
	for i := input.Start; i <= input.End; i += 1 {
		tStr := input.DataType + strconv.FormatInt(int64(i), 10) + "." + input.InputType
		inputPath := input.InputPrefix + tStr + input.InputSuffix
		f, err := os.OpenFile(inputPath, os.O_RDWR, 0666)
		if err != nil {
			panic(err)
		}
		f.Seek(0, 0)
		data = dataProcess(f, input, data)
	}
	for _, writer := range writers {
		writer <- data
	}
	log.Printf("Totally scan %d lines of data\n", len(data))
}

func OpenOutput(output Output, wg *sync.WaitGroup, writer chan []Data) {
	go func() {
		current := 0
		for true {
			select {
			case data := <-writer:
				{
					current += 1
					for _, format := range output.OutputFormat {
						function, ext := formatMap[format.OutputFormat]
						if !ext {
							panic("Invalid output format function ...")
						}
						outputStr := format.OutputFormat + "_" + strconv.FormatInt(int64(format.Start), 10) + "_" + strconv.FormatInt(int64(format.End), 10) + "_" + strconv.FormatInt(int64(current), 10)
						outputPath := output.OutputPrefix + outputStr + output.OutputSuffix
						switch output.OutputSuffix {
						case XLSX:
							{
								f := excelize.NewFile()
								index, err := f.NewSheet(DefaultSheet)
								if err != nil {
									panic(err)
								}
								function(f, data, output.Average, format.Start, format.End, output.OutputSuffix)
								f.SetActiveSheet(index)
								if err := f.SaveAs(outputPath); err != nil {
									panic(err)
								}
							}
						case TXT:
							{
								f, err := os.OpenFile(outputPath, os.O_RDWR, 0666)
								if err != nil && errors.Is(err, os.ErrNotExist) {
									f, _ = os.Create(outputPath)
								} else if err != nil {
									panic(err)
								}
								f.Truncate(0)
								f.Seek(0, 0)
								function(f, data, output.Average, format.Start, format.End, output.OutputSuffix)
							}
						}
					}
					wg.Done()
				}
			}
		}
	}()
}

func dataProcess(file *os.File, input Input, data []Data) []Data {
	log.Printf("Reading data file: %s\n", file.Name())
	var title string
	for i := 0; i < titleLine[input.InputType]; i++ {
		_, _ = fmt.Fscan(file, &title)
	}
	for true {
		switch input.DataType {
		case HEAT:
			{
				step := HeatData{}
				n, err := fmt.Fscanf(file, "%d %f %f %f %f %f\n",
					&step.TimeStep, &step.RouV, &step.Temp, &step.Volume, &step.Density, &step.Enthalpy,
				)
				if n == 0 || err != nil {
					return data
				}
				data = append(data, step)
			}
		case ENERGY:
			{
				step := EnergyData{}
				n, err := fmt.Fscanf(file, "%d %f %f %f %f %f %f\n",
					&step.TimeStep, &step.Ke, &step.Pe, &step.ETotal, &step.EMol, &step.EPair, &step.ELong,
				)
				if n == 0 || err != nil {
					return data
				}
				data = append(data, step)
			}
		}
	}
	return data
}

// TV
// temp - V
func TV(file any, data []Data, average int, startSample int, endSample int, outputFile string) {
	log.Println("Output TV data ...")
	sort.Slice(data, func(i, j int) bool {
		d1, ok1 := data[i].(HeatData)
		d2, ok2 := data[j].(HeatData)
		if !ok1 || !ok2 {
			panic("Invalid input type...")
		}
		return d1.Temp < d2.Temp
	})
	curr, n := 0, len(data)
	sumT, sumV := 0.0, 0.0
	row := 1
	for i := 0; i < n; i++ {
		d, _ := data[i].(HeatData)
		if d.Temp >= float64(startSample) && d.Temp <= float64(endSample) {
			sumT += d.Temp
			sumV += d.Volume
			curr += 1
			if curr == average {
				avgT, avgV := sumT/float64(average), sumV/float64(average)
				switch outputFile {
				case XLSX:
					{
						f, valid := file.(*excelize.File)
						if !valid {
							panic("Invalid output file type...")
						}
						cellT := "A" + strconv.FormatInt(int64(row), 10)
						cellV := "B" + strconv.FormatInt(int64(row), 10)
						f.SetCellValue(DefaultSheet, cellT, strconv.FormatFloat(avgT, 'f', 3, 64))
						f.SetCellValue(DefaultSheet, cellV, strconv.FormatFloat(avgV, 'f', 1, 64))
						row += 1
					}
				case TXT:
					{
						f, valid := file.(*os.File)
						if !valid {
							panic("Invalid output file type...")
						}
						buffer := bufio.NewWriter(f)
						fmt.Fprintf(buffer, "%.3f %.1f\n", avgT, avgV)
					}
				}
				curr = 0
				sumT, sumV = 0.0, 0.0
			}
		} else if d.Temp > float64(endSample) {
			break
		}
	}
}
