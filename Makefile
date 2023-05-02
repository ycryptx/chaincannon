.PHONY: # ignore


help:
	@perl -nle'print $& if m{^[a-zA-Z_-]+:.*?## .*$$}' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

run: # run the package
	go run ./cmd/chaincannon/main.go

build: # build the package
	go build ./cmd/chaincannon

run-example: # runs an example cosmos chain and benchmarks that chain using dummy transactions
	docker-compose up example-cosmos

test: # run all tests
	gotestsum --format testname -- ./...