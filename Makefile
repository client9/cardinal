
generate:
	go generate .
	go build .

lint:
	go mod tidy
	gofmt -w -s *.go
	golangci-lint run .

test:
	go generate ./...
	go test

clean:
	rm -f repl cmd/repl/repl
	rm -f wrapgen
	rm -f *.bak*
