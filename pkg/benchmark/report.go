package benchmark

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/ttacon/chalk"
)

func (report *Report) RecordLatency(ctx context.Context, start time.Time, end time.Time) {
	report.Latencies.RecordValue(end.UnixMilli() - start.UnixMilli())
}

func (report *Report) RecordBlockTime(ctx context.Context, start time.Time, end time.Time) {
	report.BlockTimes.RecordValue(end.UnixMilli() - start.UnixMilli())
}

func (report *Report) PrintReport(ctx context.Context) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowSeparator("-")
	table.SetHeader([]string{
		"Stat",
		"2.5%",
		"50%",
		"97.5%",
		"99%",
		"Avg",
		"Stdev",
		"Max",
		"Count",
	})
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor})
	fmt.Println(`
█▀▀ █░█ ▄▀█ █ █▄░█ █▀▀ ▄▀█ █▄░█ █▄░█ █▀█ █▄░█
█▄▄ █▀█ █▀█ █ █░▀█ █▄▄ █▀█ █░▀█ █░▀█ █▄█ █░▀█`)
	fmt.Println("")
	report.printTxLatencies(table)
	report.printBlockTimes(table)
	table.Render()
	fmt.Println("")
	fmt.Println("")
}

func (report *Report) printTxLatencies(table *tablewriter.Table) {
	table.Append([]string{
		chalk.Bold.TextStyle("Tx Latency"),
		fmt.Sprintf("%v ms", report.Latencies.ValueAtPercentile(2.5)),
		fmt.Sprintf("%v ms", report.Latencies.ValueAtPercentile(50)),
		fmt.Sprintf("%v ms", report.Latencies.ValueAtPercentile(97.5)),
		fmt.Sprintf("%v ms", report.Latencies.ValueAtPercentile(99)),
		fmt.Sprintf("%.2f ms", report.Latencies.Mean()),
		fmt.Sprintf("%.2f ms", report.Latencies.StdDev()),
		fmt.Sprintf("%v ms", report.Latencies.Max()),
		fmt.Sprintf("%v", report.Latencies.TotalCount()),
	})
}

func (report *Report) printBlockTimes(table *tablewriter.Table) {
	table.Append([]string{
		chalk.Bold.TextStyle("Block Time"),
		fmt.Sprintf("%v ms", report.BlockTimes.ValueAtPercentile(2.5)),
		fmt.Sprintf("%v ms", report.BlockTimes.ValueAtPercentile(50)),
		fmt.Sprintf("%v ms", report.BlockTimes.ValueAtPercentile(97.5)),
		fmt.Sprintf("%v ms", report.BlockTimes.ValueAtPercentile(99)),
		fmt.Sprintf("%.2f ms", report.BlockTimes.Mean()),
		fmt.Sprintf("%.2f ms", report.BlockTimes.StdDev()),
		fmt.Sprintf("%v ms", report.BlockTimes.Max()),
		fmt.Sprintf("%v", report.BlockTimes.TotalCount()),
	})
}
