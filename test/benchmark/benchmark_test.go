package benchmark_test

import (
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/ycryptx/chaincannon/pkg/benchmark"
)

var (
	endpoint           = "0.0.0.0:123"
	tendermintEndpoint = ""
	maxCores           = 4
	maxDuration        = time.Duration(10000) * time.Second
)

type addTest struct {
	testName           string
	endpoint           string
	tendermintEndpoint string
	txPaths            []string
	duration           int
	amount             int
	threads            int
	expected           benchmark.Recipe
}

var addTests = []addTest{
	{"when amount is set duration should == maxDuration", endpoint, tendermintEndpoint, []string{"/path1"}, 123, 456, 0, benchmark.Recipe{
		Endpoint:           endpoint,
		TendermintEndpoint: tendermintEndpoint,
		Duration:           maxDuration,
		Amount:             456,
		Runs:               []benchmark.Run{{TxPaths: []string{"/path1"}}},
	}},
	{"when amount == 0 duration should be used", endpoint, tendermintEndpoint, []string{"/path1"}, 123, 0, 0, benchmark.Recipe{
		Endpoint:           endpoint,
		TendermintEndpoint: tendermintEndpoint,
		Duration:           time.Duration(123) * time.Second,
		Amount:             0,
		Runs:               []benchmark.Run{{TxPaths: []string{"/path1"}}},
	}},
	{"should split runs equally between specified threads", endpoint, tendermintEndpoint, []string{"/path1", "/path2"}, 123, 0, 2, benchmark.Recipe{
		Endpoint:           endpoint,
		TendermintEndpoint: tendermintEndpoint,
		Duration:           time.Duration(123) * time.Second,
		Amount:             0,
		Runs:               []benchmark.Run{{TxPaths: []string{"/path1"}}, {TxPaths: []string{"/path2"}}},
	}},
	{"should more than one tx file in same thread if tx_files > threads", endpoint, tendermintEndpoint, []string{"/path1", "/path2", "/path3"}, 123, 0, 2, benchmark.Recipe{
		Endpoint:           endpoint,
		TendermintEndpoint: tendermintEndpoint,
		Duration:           time.Duration(123) * time.Second,
		Amount:             0,
		Runs:               []benchmark.Run{{TxPaths: []string{"/path1", "/path3"}}, {TxPaths: []string{"/path2"}}},
	}},
	{"should cap to max threads if specified threads is greater than available cores", endpoint, tendermintEndpoint, []string{"/path1", "/path2", "/path3", "/path4", "/path5"}, 123, 0, 5, benchmark.Recipe{
		Endpoint:           endpoint,
		TendermintEndpoint: tendermintEndpoint,
		Duration:           time.Duration(123) * time.Second,
		Amount:             0,
		Runs:               []benchmark.Run{{TxPaths: []string{"/path1", "/path5"}}, {TxPaths: []string{"/path2"}}, {TxPaths: []string{"/path3"}}, {TxPaths: []string{"/path4"}}},
	}},
	{"should run correctly when many more files than threads", endpoint, tendermintEndpoint, []string{"/path1", "/path2", "/path3", "/path4", "/path5"}, 123, 0, 2, benchmark.Recipe{
		Endpoint:           endpoint,
		TendermintEndpoint: tendermintEndpoint,
		Duration:           time.Duration(123) * time.Second,
		Amount:             0,
		Runs:               []benchmark.Run{{TxPaths: []string{"/path1", "/path3", "/path5"}}, {TxPaths: []string{"/path2", "/path4"}}},
	}},
}

func TestProcessFlags(t *testing.T) {
	for _, test := range addTests {
		if diff := deep.Equal(test.expected, *benchmark.ProcessFlags(test.endpoint, test.tendermintEndpoint, test.txPaths, test.duration, test.amount, test.threads, maxCores, maxDuration)); diff != nil {
			t.Errorf("%s: %q", test.testName, diff)
		}
	}
}
