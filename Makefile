

build: generate
	go run cmd/wrapgen/main.go cmd/wrapgen/reflect_helper.go cmd/wrapgen/symbols.go
	go build ./...
	(cd cmd/repl; go build .; mv repl ../..)
	go test ./...

generate:
	go get golang.org/x/tools/cmd/stringer
	go generate ./...

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
	rm -rf wrapped
	rm -f builtin_setup.go
	rm -f attribute_string.go

setup:
	go get golang.org/x/tools/cmd/stringer
