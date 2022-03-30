package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

const DataFile = "loremipsum.txt"

// Return the word frequencies of the text argument.
//
// Split load optimally across processor cores.
func WordCount(text string) map[string]int {
	words := strings.Fields(strings.ToLower(text))
	wordlists := Split(words, 30000)
	ch := make(chan map[string]int)
	out := make(chan map[string]int)
	wgc := new(sync.WaitGroup)
	wgr := new(sync.WaitGroup)
	wgc.Add(len(wordlists))
	wgr.Add(len(wordlists))
	for _, ws := range wordlists {
		go Counter(ws, wgc, ch)
	}
	go Reducer(wgr, ch, out)
	wgc.Wait()
	wgr.Wait()
	close(ch)
	return <-out
}

func Split(words []string, size int) [][]string {
	out := make([][]string, len(words)/size+1)
	exit := false
	bot := 0
	top := size
	i := 0
	for !exit {
		if top > len(words) {
			top = len(words)
			exit = true
		}

		out[i] = words[bot:top]

		bot += size
		top += size
		i += 1
	}
	return out
}

func Counter(words []string, wgc *sync.WaitGroup, ch chan<- map[string]int) {
	freqs := make(map[string]int)
	for _, w := range words {
		freqs[strings.Trim(w, ",.")] += 1
	}
	ch <- freqs
	wgc.Done()
}

func Reducer(wgr *sync.WaitGroup, ch <-chan map[string]int, out chan<- map[string]int) {
	freqs := make(map[string]int)
	for fs := range ch {
		for key, value := range fs {
			freqs[key] += value
		}
		wgr.Done()
	}
	out <- freqs
}

// Benchmark how long it takes to count word frequencies in text numRuns times.
//
// Return the total time elapsed.
func benchmark(text string, numRuns int) int64 {
	start := time.Now()
	for i := 0; i < numRuns; i++ {
		WordCount(text)
	}
	runtimeMillis := time.Since(start).Nanoseconds() / 1e6

	return runtimeMillis
}

// Print the results of a benchmark
func printResults(runtimeMillis int64, numRuns int) {
	fmt.Printf("amount of runs: %d\n", numRuns)
	fmt.Printf("total time: %d ms\n", runtimeMillis)
	average := float64(runtimeMillis) / float64(numRuns)
	fmt.Printf("average time/run: %.2f ms\n", average)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	data, err := os.ReadFile(DataFile)
	check(err)

	numRuns := 100
	runtimeMillis := benchmark(string(data), numRuns)
	printResults(runtimeMillis, numRuns)
}
