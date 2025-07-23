
build:
	go run cmd/wrapgen/main.go cmd/wrapgen/reflect_helper.go
	go run cmd/wrapgen/main.go cmd/wrapgen/reflect_helper.go -setup builtin_setup.go
	go build .
	(cd cmd/repl; go build .; mv repl ../..)

lint:
	go mod tidy
	gofmt -w -s *.go
	golangci-lint run .

test:  build
	go test .

clean:
	rm -f repl cmd/repl/repl
	rm -f main wrapgen
	rm -f *.bak*
	rm -f *_wrappers.go
	rm -f builtin_setup.go
