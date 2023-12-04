.PHONY: # ignore

DATE := $(shell date '+%Y-%m-%dT%H:%M:%S')
HEAD = $(shell git rev-parse HEAD)
LD_FLAGS = -X github.com/ycryptx/chaincannon/version.Head='$(HEAD)' \
	-X github.com/ycryptx/chaincannon/version.Date='$(DATE)'
BUILD_FLAGS = -mod=readonly -ldflags='$(LD_FLAGS)'

help:
	@perl -nle'print $& if m{^[a-zA-Z_-]+:.*?## .*$$}' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

install: # install the binary
	@echo Installing Chaincannon...
	@go install $(BUILD_FLAGS) ./...
	@echo Chaincannon installed!

run: # run the package
	go run ./cmd/chaincannon/main.go

build: # build the package
	go build ./cmd/chaincannon

gen-example-txs: # NOTE: this has been pre-generated so no need to run it. Used to generate signed transactions files to be used in the example.
	chmod +x ./example/cosmos/data/tx_gen.sh
	./example/cosmos/data/tx_gen.sh

build-example-docker: # build the example docker image
	docker build -t example-cosmos example/cosmos/chain

setup-example: # set an example cosmos chain
	docker-compose up example-cosmos --force-recreate

run-example: # runs the example cosmos benchmark
	go build ./cmd/chaincannon
	./chaincannon -chain cosmos -endpoint 0.0.0.0:9090 -tendermintEndpoint 0.0.0.0:26657 -duration 60 -tx-file ./example/cosmos/data/run1.json  -tx-file ./example/cosmos/data/run3.json -tx-file ./example/cosmos/data/run4.json

run-tests: # run all tests
	gotestsum --format testname -- ./...