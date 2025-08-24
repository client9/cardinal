
test:
	go test .

lint:
	gofmt -w -s *.go

testre:
	go test -v  -test.run=SRE

benchre:
	go test  -test.run=SRE -bench=SRE
benchcompile:
	go test  -test.run=Compile -bench=Compile

benchmemre:
	go test  -test.run=SRE -bench=SRE -benchmem

# go tool pprof cpu.prof
profilere:
	go test -cpuprofile cpu.prof -memprofile mem.prof -test.run=SRE -bench=SREProfile

