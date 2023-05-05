package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/ycryptx/chaincannon/pkg/benchmark"
	"github.com/ycryptx/chaincannon/pkg/cosmos"
	"github.com/ycryptx/chaincannon/pkg/flags"
)

var SUPPORTED_CHAINS = []string{"cosmos"}

var txPaths flags.StringArray

func main() {
	chain := flag.String("chain", "", fmt.Sprintf("The blockchain type (e.g. %s).", strings.Join(SUPPORTED_CHAINS, ", ")))
	endpoint := flag.String("endpoint", "", "The node's RPC endpoint to call.")
	flag.Var(&txPaths, "tx-file", "Path to a file containing signed transactions. This flag can be used more than once.")
	duration := flag.Int("duration", 10, "The number of seconds to run the benchmark.")
	amount := flag.Int("amount", 0, "The number of requests to make before exiting the benchmark. If set, duration is ignored.")
	threads := flag.Int("threads", 0, "The number of concurrent threads to use to make requests. (default: max)")
	flag.Parse()

	recipe := benchmark.ProcessFlags(*endpoint, txPaths, *duration, *amount, *threads, runtime.NumCPU())

	switch *chain {
	case "cosmos":
		err := cosmos.Handle(recipe)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	default:
		fmt.Printf("Chaincannon is a blockchain benchmarking tool. Currently supported chains are: %s\n", strings.Join(SUPPORTED_CHAINS, ", "))
		fmt.Println("Usage: chaincannon [opts]")
		flag.PrintDefaults()
		return
	}

}
