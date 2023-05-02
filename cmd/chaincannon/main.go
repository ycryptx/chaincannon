package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/ycryptx/chaincannon/pkg/cosmos"
)

type Color string

const (
	ColorBlack  Color = "\u001b[30m"
	ColorRed          = "\u001b[31m"
	ColorGreen        = "\u001b[32m"
	ColorYellow       = "\u001b[33m"
	ColorBlue         = "\u001b[34m"
	ColorReset        = "\u001b[0m"
)

var SUPPORTED_CHAINS = []string{"cosmos"}

func colorize(color Color, message string) {
	fmt.Println(string(color), message, string(ColorReset))
}

func main() {
	chain := flag.String("chain", "", fmt.Sprintf("The blockchain type (e.g. %s).", strings.Join(SUPPORTED_CHAINS, ", ")))
	endpoint := flag.String("endpoint", "", "The node's RPC endpoint to call.")
	connections := flag.Int("connections", 1, "The number of concurrent connections to make.")
	duration := flag.Int("duration", 10, "The number of seconds to run the benchmark.")
	amount := flag.Int("amount", 0, "The number of requests to make before exiting the benchmark. If set, duration is ignored.")
	threads := flag.Int("threads", 0, "The number of worker threads to use to make requests. (default: max)")
	flag.Parse()

	switch *chain {
	case "cosmos":
		cosmos.Handle(*endpoint, *connections, *duration, *amount, *threads)
	default:
		fmt.Printf("Chaincannon is a blockchain benchmarking tool. Currently supported chains are: %s\n", strings.Join(SUPPORTED_CHAINS, ", "))
		fmt.Println("Usage: chaincannon [opts]")
		flag.PrintDefaults()
		return
	}

}
