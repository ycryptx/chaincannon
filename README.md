# Chaincannon

Chaincannon is an all-in-one blockchain benchmarking tool written in go, currently only available for [Cosmos-SDK](https://github.com/cosmos/cosmos-sdk) based chains. As the name suggests, this tool is inspired by Matteo Collina's [Autocannon](https://github.com/mcollina/autocannon)

The goal of Chaincannon is to be a memory-efficient, lightweight tool for stress testing and benchmarking a blockchain. A single benchmark can fire tens of thousands of concurrent transactions in a controlled, repeatable way. That makes it possible to integrate chaincannon as part of a the blockchain development`s ci/cd pipeline.

Because the generation, encoding, and signing of transactions is different between Cosmos-SDK chains, Chaincannon has been designed to be chain-agnostic and thus the transaction generation and signing is the responsibility of the user. That means that a chaincannon benchmark run must be provided with file/s containing the transactions to execute in the benchmark. 

## Installation

```bash
$ make build

```

## Usage
```
$ chaincannon 
Chaincannon is a blockchain benchmarking tool. Currently supported chains are: cosmos
Usage: chaincannon [opts]
  -amount int
        The number of requests to make before exiting the benchmark. If set, duration is ignored.
  -chain string
        The blockchain type (e.g. cosmos).
  -duration int
        The number of seconds to run the benchmark. (default 30)
  -endpoint string
        The node's RPC endpoint to call.
  -threads int
        The number of concurrent threads to use to make requests. (default: max)
  -tx-file value
        Path to a file containing signed transactions. This flag can be used more than once.
```

## Sample Run

```bash
$ chaincannon -chain cosmos -endpoint 0.0.0.0:9090 -duration 60 -tx-file ./example/cosmos/data/run1.json  -tx-file ./example/cosmos/data/run3.json -tx-file ./example/cosmos/data/run4.json
 100% |█████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████| (3/100, 2 it/s)        

█▀▀ █░█ ▄▀█ █ █▄░█ █▀▀ ▄▀█ █▄░█ █▄░█ █▀█ █▄░█
█▄▄ █▀█ █▀█ █ █░▀█ █▄▄ █▀█ █░▀█ █░▀█ █▄█ █░▀█

+------------+------+---------+---------+---------+------------+----------+---------+-------+
|    STAT    | 2.5% |   50%   |  97.5%  |   99%   |    AVG     |  STDEV   |   MAX   | COUNT |
+------------+------+---------+---------+---------+------------+----------+---------+-------+
| Tx         | 0 ms | 613 ms  | 648 ms  | 648 ms  | 612.78 ms  | 22.94 ms | 648 ms  |    18 |
| Latency    |      |         |         |         |            |          |         |       |
+------------+------+---------+---------+---------+------------+----------+---------+-------+
| Block      | 0 ms | 1013 ms | 1016 ms | 1016 ms | 1014.50 ms | 1.50 ms  | 1016 ms |     2 |
| Time       |      |         |         |         |            |          |         |       |
+------------+------+---------+---------+---------+------------+----------+---------+-------+
| TPS        |    0 |      17 |      17 |      17 |      17.00 |     0.00 |      17 |
+------------+------+---------+---------+---------+------------+----------+---------+-------+

Benchmark ran 3 concurrent processes
Executed 18 txs and took 1.648335 seconds                                                                                                               
```

## Benchmark Metrics

1. Transaction Latency
      - Represented and stored as a histogram
      - Measures time interval between when transaction is sent to when the transaction is included in a confirmed block. The time used is the timestamp in the benchmark machine before transaction is sent, and the timestamp of the benchmark machine when it hears about the block in which the transcion was included
2. Block time
      - Represented and stored as a histogram
      - Measures time interval between blocks. The time used is the timestamp of the block itself
3. TPS (Transactions per second)
      - Represented and stored as a histogram
      - Measures the # of transactions per block per blockTime per second

Currently unimplemented:
- Transaction success rate (txs dropped)
- Uncle Block rate (if applicable to the chain)

## About Transaction Files

Transaction files should contain protobuf-encoded signed transactions. Each transaction in a file should be new-line-delimited. Each transaction file can be run concurrently or synchronously depending on the `--threads` user flags and on how many cores are available on the machine running the benchmark. An example of how transactions can be generated can [be found here](./example/cosmos/data/tx_gen.sh). Note that the demo transactions are specific to [the demo blockchain](./example/cosmos/chain/Dockerfile). The demo transactions use [Cosmos' x/bank module](https://docs.cosmos.network/v0.46/modules/bank/) to pass tokens between pre-funded users. You will probably want to generate transactions specific to the modules in your blockchain.

## Demo

Use this demo to see how chaincannon works. It spins up a local cosmos-sdk chain using Docker (you must have Docker [installed](https://docs.docker.com/get-docker/)), and hits the chain with a series of concurrent transactions from a few transacions files. You may look at the [`Makefile`](./Makefile) to see how the chaincannon cli command is structured.

```bash
$ make setup-example
```
And in another terminal window run:
```bash
$ make run-example
```

## Licensing

[`Apache License 2.0`](./LICENSE)

### Authors 

[`Yonatan Medina`](github.com/ycryptx)
