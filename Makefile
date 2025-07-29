

build: generate
	go run cmd/wrapgen/main.go cmd/wrapgen/reflect_helper.go cmd/wrapgen/symbols.go
	go build ./...
	(cd cmd/repl; go build .; mv repl ../..)

generate:
	go get golang.org/x/tools/cmd/stringer
	go generate ./...

lint:
	go mod tidy
	gofmt -w -s *.go core/*.go stdlib/*.go cmd/wrapgen/*.go cmd/repl/*.go tests/integration/*.go
	golangci-lint run .

# Run all tests
test: build test-unit test-integration

# Unit tests (fast) - pure Go function tests
test-unit:
	go test ./stdlib/... ./core/...

# Integration tests (slower) - end-to-end evaluation tests  
test-integration:
	go test ./tests/integration/...

# Core infrastructure tests (parser, evaluator, etc.)
test-core:
	go test -run="TestParser|TestLexer|TestEvaluator|TestFunctionRegistry|TestError" .

# Quick smoke test
test-quick:
	go test -short ./...

# Coverage report
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean:
	rm -f repl cmd/repl/repl
	rm -f main wrapgen
	find . -name '*.bak*' | xargs rm -f
	rm -rf wrapped
	rm -f builtin_setup.go
	rm -f attribute_string.go

setup:
	go get golang.org/x/tools/cmd/stringer
