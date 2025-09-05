

build: 
	go run ./cmd/geninit 
	go build ./...
	(cd cmd/cardinal; go build .; mv cardinal ../..)

generate:
	go generate ./...

lint:
	go mod tidy
	find . -name '*.go' | xargs gofmt -w -s
	golangci-lint run .

# Run all tests
test: build
	go test ./...


prof:
	go test -cpuprofile cpu.prof -memprofile mem.prof -test.run=SRE -bench=SRE2b/MatchAny3,binding ./core/...

# Coverage report
# Go default is crap
# cov0 is red (not covered)
# cov8 is the green (covered)
cover:
	rm -f coverage*
	go test -coverpkg=./... -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html-tmp
	cat coverage.html-tmp | sed 's/background: black/background: whitesmoke/g' | sed 's/80, 80, 80/0,0,0/g' | sed 's/Menlo/ui-monospace/g' | sed 's/bold/normal/g' | sed 's/rgb(192, 0, 0)/rgb(255,0,0);font-weight:bold;/g' > coverage.html

clean:
	rm -f cpu.prof mem.prof
	rm -f repl cmd/repl/repl
	find . -name '*.bak*' | xargs rm -f
	rm -f engine/attribute_string.go
	rm -f core/symbol/symbols.go

setup:
	go get golang.org/x/tools/cmd/stringer
	go get -u github.com/google/pprof
