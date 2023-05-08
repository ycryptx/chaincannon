package benchmark

import (
	"context"
	"time"

	"github.com/HdrHistogram/hdrhistogram-go"
)

func InitMonitoring(ctx context.Context, cancel context.CancelFunc) *Monitoring {
	recipe, _ := ctx.Value("recipe").(*Recipe)
	highestDiscernibleValue := recipe.Duration
	report := &Report{
		Latencies:  hdrhistogram.New(1, highestDiscernibleValue.Milliseconds(), 5),
		BlockTimes: hdrhistogram.New(1, highestDiscernibleValue.Milliseconds(), 5),
	}
	m := &Monitoring{Report: report, Stream: make(chan *TxPending), txs: map[string]*TxPending{}}
	go m.monitor(ctx, cancel)
	return m
}

func (monitoring *Monitoring) IsOurTx(ctx context.Context, hash string) bool {
	monitoring.mu.RLock()
	_, ok := monitoring.txs[hash]
	monitoring.mu.RUnlock()
	return ok
}

func (monitoring *Monitoring) NoMorePendingTxs() bool {
	var keyLen int
	monitoring.mu.Lock()
	keyLen = len(monitoring.txs)
	monitoring.mu.Unlock()
	return keyLen == 0
}

func (monitoring *Monitoring) RecordTxs(ctx context.Context, hashes []string, endTime time.Time) {
	for _, hash := range hashes {
		tx, ok := monitoring.txs[hash]
		if !ok {
			panic("cannot find tx")
		}
		monitoring.Report.RecordLatency(ctx, *tx.Start, endTime)
	}
	monitoring.mu.Lock()
	for _, hash := range hashes {
		delete(monitoring.txs, hash)
	}
	monitoring.mu.Unlock()
}

func (monitoring *Monitoring) monitor(ctx context.Context, cancel context.CancelFunc) {
	for {
		tx := <-monitoring.Stream
		monitoring.addTx(ctx, tx, cancel)
	}
}

func (monitoring *Monitoring) addTx(ctx context.Context, tx *TxPending, cancel context.CancelFunc) {
	recipe, _ := ctx.Value("recipe").(*Recipe)
	if recipe.Amount > 0 && monitoring.TxFired >= recipe.Amount {
		// don't track txs after they exceed the max tx amount
		return
	}
	monitoring.mu.Lock()
	if _, ok := monitoring.txs[tx.Hash]; !ok {
		monitoring.txs[tx.Hash] = tx
		monitoring.TxFired += 1
	}
	monitoring.mu.Unlock()
}
