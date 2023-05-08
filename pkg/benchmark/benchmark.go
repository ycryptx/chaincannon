package benchmark

import (
	"time"
)

// Processes user CLI flags and determines the level of concurrency by which to run the benchmark.
func ProcessFlags(endpoint string, txPaths []string, duration int, amount int, threads int, maxCores int, maxDuration time.Duration) *Recipe {
	// cap threads to max available cores, or if threads is unset default to max available cores
	if maxCores < threads || threads == 0 {
		threads = maxCores
	}

	asDuration := time.Duration(duration) * time.Second
	if amount > 0 || asDuration > maxDuration {
		asDuration = maxDuration
	}

	runs := []Run{}
	for i, path := range txPaths {
		if len(runs) < i+1 && i < threads {
			runs = append(runs, Run{})
		}
		putInThread := i % (threads)
		runs[putInThread].TxPaths = append(runs[putInThread].TxPaths, path)
	}

	return &Recipe{
		Endpoint: endpoint,
		Duration: asDuration,
		Amount:   amount,
		Runs:     runs,
	}
}
