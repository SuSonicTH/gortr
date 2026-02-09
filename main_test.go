package main

import (
	"log"
	"os"
	"runtime/pprof"
	"testing"
)

// run benchmark/profiling with `goz test -bench=.`
// show profile data with `pprof -http=: profile`
func BenchmarkRefresh(b *testing.B) {
	f, err := os.Create("profile")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)

	defer pprof.StopCPUProfile()

	for b.Loop() {
		refresh()
	}

}
