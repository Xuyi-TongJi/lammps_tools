package main

import (
	"atom/postProcess"
	"atom/preProcess"
	. "atom/util"
	"strconv"
	"sync"
)

var suf = "/Users/xuyi/Desktop/phase520/"
var path = "/Users/xuyi/Desktop/phase520/"

const (
	s    = 1
	e    = 3
	it   = postProcess.HEAT
	task = 1
	dt   = postProcess.ENERGY
)

func main() {
	for i := 260; i <= 290; i += 30 {
		for j := 0; j <= 3; j++ {
			path = suf + strconv.FormatInt(int64(i), 10) + "/" + strconv.FormatInt(int64(j), 10) + "/"
			post()
		}
	}
}

func post() {
	wg := &sync.WaitGroup{}
	wg.Add(task)
	// output1
	var outputs []OutputFormatSample
	writers := make([]chan []Data, 0)
	writer := make(chan []Data)
	writers = append(writers, writer)
	outputs = append(outputs,
		OutputFormatSample{OutputFormat: postProcess.TVDiagram, Start: -1000, End: 1000})
	output := Output{
		Average:      100,
		OutputPrefix: path + "data/",
		OutputSuffix: postProcess.TXT,
		OutputFormat: outputs,
	}
	postProcess.OpenOutput(output, wg, writer)
	// HTTP
	for i := 0; i < task; i++ {
		input := Input{
			Start:       s,
			End:         e,
			InputType:   it,
			DataType:    dt,
			InputPrefix: path,
			InputSuffix: postProcess.TXT,
		}
		postProcess.OpenInput(input, writers)
	}
	wg.Wait()
}

func pre() {
	path := "/Users/xuyi/Desktop/npg_data/npg_atom.txt"
	input := PreInput{
		InputType: preProcess.Tme,
		Path:      path,
	}
	preProcess.AddMol(input)
}
