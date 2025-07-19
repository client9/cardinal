

lint:
	go mod tidy
	gofmt -w -s *.go
	golangci-lint run .
test:
	go test

clean:
	rm -f *.bak*
