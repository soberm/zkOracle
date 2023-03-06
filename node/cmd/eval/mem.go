package main

import (
	"runtime"
	"time"
)

func MemUsage(stop chan struct{}) uint64 {
	var m runtime.MemStats
	var memory uint64
loop:
	for {
		select {
		case <-stop: // triggered when the stop channel is closed
			break loop // exit
		default:
			runtime.ReadMemStats(&m)
			current := bToMb(m.Sys)
			if memory < current {
				memory = current
			}
			time.Sleep(time.Millisecond)
		}
	}
	return memory
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
