package benchmark

import "time"

type Recipe struct {
	Duration time.Duration // the duration of the benchmark
	Amount   int           // the number of transactions to post. If set duration must be ignored
	Endpoint string        // the URL address of the blockchain node to benchmark
	Runs     []Run         // multiple runs are run concurrently
}

type Run struct {
	TxPaths []string // a run executes one or more transaction files synchronously
}
