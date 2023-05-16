package main

import (
	"atom/phase"
	"sync"
)

func main() {
	//count := (phase.TEnd - phase.TStart) / phase.DeltaT
	wg := &sync.WaitGroup{}
	wg.Add(1)
	phase.Integration(100, phase.TV, ".temp", "output", wg)
	wg.Wait()
}
