package main

import (
	"context"
	"flag"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/ycryptx/chaincannon/pkg/benchmark"
	"github.com/ycryptx/chaincannon/pkg/cosmos"
	"github.com/ycryptx/chaincannon/pkg/flags"
	"github.com/ycryptx/chaincannon/pkg/logger"
)

// Default Benchmark duration is 5 mintues
const maxDuration = time.Minute * 5

var SUPPORTED_CHAINS = []string{"cosmos"}

var txPaths flags.StringArray

func main() {
	log := logger.InitLogger()
	chain := flag.String("chain", "", fmt.Sprintf("The blockchain type (e.g. %s).", strings.Join(SUPPORTED_CHAINS, ", ")))
	endpoint := flag.String("endpoint", "", "The node's RPC endpoint to call.")
	flag.Var(&txPaths, "tx-file", "Path to a file containing signed transactions. This flag can be used more than once.")
	duration := flag.Int("duration", 10, "The number of seconds to run the benchmark.")
	amount := flag.Int("amount", 0, "The number of requests to make before exiting the benchmark. If set, duration is ignored.")
	threads := flag.Int("threads", 0, "The number of concurrent threads to use to make requests. (default: max)")
	flag.Parse()

	recipe := benchmark.ProcessFlags(*endpoint, txPaths, *duration, *amount, *threads, runtime.NumCPU(), maxDuration)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "log", log)
	ctx = context.WithValue(ctx, "recipe", recipe)

	bar := progressbar.Default(100)
	go func() {
		for i := 0; i < 100; i++ {
			bar.Add(1)
			time.Sleep(recipe.Duration / 100)
		}
	}()
	ctx = context.WithValue(ctx, "bar", bar)

	switch *chain {
	case "cosmos":
		cosmos.Handle(ctx)
	default:
		fmt.Printf("Chaincannon is a blockchain benchmarking tool. Currently supported chains are: %s\n", strings.Join(SUPPORTED_CHAINS, ", "))
		fmt.Println("Usage: chaincannon [opts]")
		flag.PrintDefaults()
		return
	}

}
