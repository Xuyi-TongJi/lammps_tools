package main

import (
	"atom/postProcess"
	. "atom/util"
	"sync"
)

const (
	path = "/Users/xuyi/Desktop/phase518_tme/280/0/"
	s    = 1
	e    = 3
	it   = "heat"
	task = 1
	dt   = postProcess.HEAT
)

func main() {
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
		OutputSuffix: postProcess.XLSX,
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
