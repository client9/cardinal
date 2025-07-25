
build:
	go run cmd/wrapgen/main.go cmd/wrapgen/reflect_helper.go
	go run cmd/wrapgen/main.go cmd/wrapgen/reflect_helper.go -setup builtin_setup.go
	go build ./...
	(cd cmd/repl; go build .; mv repl ../..)

lint:
	go mod tidy
	gofmt -w -s *.go core/*.go stdlib/*.go cmd/wrapgen/*.go cmd/repl/*.go
	golangci-lint run .

test:  build
	go test ./...

clean:
	rm -f repl cmd/repl/repl
	rm -f main wrapgen
	find . -name '*.bak*' | xargs rm -f
	rm -f *_wrappers.go
	rm -f builtin_setup.go
