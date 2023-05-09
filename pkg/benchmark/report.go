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

func (report *Report) RecordBenchmarkDuration(ctx context.Context, start time.Time, end time.Time) {
	report.BenchmarkDuration = end.Sub(start)
}

func (report *Report) RecordTPS(ctx context.Context, blockTimeStart time.Time, blockTimeEnd time.Time, txsInBlock int) {
	tps := float64(txsInBlock) / (float64((blockTimeEnd.Sub(blockTimeStart)).Seconds()))
	report.TPS.RecordValue(int64(tps))
}

func (report *Report) PrintReport(ctx context.Context) {
	recipe, _ := ctx.Value("recipe").(*Recipe)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowSeparator("-")
	table.SetRowLine(true)
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
	report.printTPS(table)
	table.Render()
	fmt.Println("")
	fmt.Printf("Benchmark ran %d concurrent processes\n", len(recipe.Runs))
	fmt.Printf("Executed %d txs and took %f seconds", report.Latencies.TotalCount(), report.BenchmarkDuration.Seconds())
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

func (report *Report) printTPS(table *tablewriter.Table) {
	table.Append([]string{
		chalk.Bold.TextStyle("TPS"),
		fmt.Sprintf("%v", report.TPS.ValueAtPercentile(2.5)),
		fmt.Sprintf("%v", report.TPS.ValueAtPercentile(50)),
		fmt.Sprintf("%v", report.TPS.ValueAtPercentile(97.5)),
		fmt.Sprintf("%v", report.TPS.ValueAtPercentile(99)),
		fmt.Sprintf("%.2f", report.TPS.Mean()),
		fmt.Sprintf("%.2f", report.TPS.StdDev()),
		fmt.Sprintf("%v", report.TPS.Max()),
	})
}
