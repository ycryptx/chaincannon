.PHONY: # ignore

current_dir := $(notdir $(patsubst %/,%,$(dir $(mkfile_path))))

help:
	@perl -nle'print $& if m{^[a-zA-Z_-]+:.*?## .*$$}' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

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
	./chaincannon -chain cosmos -endpoint 0.0.0.0:9090 -duration 60 -tx-file ./example/cosmos/data/run1.json -tx-file ./example/cosmos/data/run2.json -tx-file ./example/cosmos/data/run3.json -tx-file ./example/cosmos/data/run4.json

run-tests: # run all tests
	gotestsum --format testname -- ./...