package benchmark

import (
	"sync"
	"time"

	"github.com/HdrHistogram/hdrhistogram-go"
)

// a recipe is an internal representation of a user-initiated benchmark
type Recipe struct {
	Duration time.Duration // the duration of the benchmark
	Amount   int           // the number of transactions to post. If set duration must be ignored
	Endpoint string        // the URL address of the blockchain node to benchmark
	Runs     []Run         // multiple runs are run concurrently
}

// a single synchronous execution of signed transaction file/s
type Run struct {
	TxPaths []string // a run executes one or more transaction files synchronously
}

type Monitoring struct {
	mu      sync.RWMutex
	txs     map[string]*TxPending
	Report  *Report
	Stream  chan *TxPending
	TxFired int
	Done    bool
}

type TxPending struct {
	Hash  string
	Start *time.Time
}

type Report struct {
	Latencies         *hdrhistogram.Histogram
	BlockTimes        *hdrhistogram.Histogram
	TPS               *hdrhistogram.Histogram
	BenchmarkDuration time.Duration
}
