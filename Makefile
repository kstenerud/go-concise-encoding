.PHONY: clean test testnormal testpurego test32bit cover benchmark compare lint

clean:
	go clean
	rm -f *.out *.test

test: testnormal testpurego test32bit

testnormal:
	$(info ** Running Tests (normal))
	go test -coverprofile=coverage.out ./...

testpurego:
	$(info ** Running Tests (purego))
	go test -tags purego ./...

test32bit:
	$(info ** Running Tests (32-bit))
	GOARCH=386 go test ./...

cover:
	$(info To see coverage: go tool cover -html=coverage.out)
	go test -coverprofile=coverage.out ./...

benchmark:
	$(info To see profile: go tool pprof cpuprofile.out)
	go test -run BenchmarkCBE* -bench BenchmarkCBE* -benchmem -memprofile memprofile.out -cpuprofile cpuprofile.out

benchmarkMarshal:
	$(info To see profile: go tool pprof cpuprofile.out)
	go test -run BenchmarkCBEMarshal -bench BenchmarkCBEMarshal -benchmem -memprofile memprofile.out -cpuprofile cpuprofile.out

benchmarkUnmarshal:
	$(info To see profile: go tool pprof cpuprofile.out)
	go test -run BenchmarkCBEUnmarshal -bench BenchmarkCBEUnmarshal -benchmem -memprofile memprofile.out -cpuprofile cpuprofile.out

benchmarkRules:
	$(info To see profile: go tool pprof cpuprofile.out)
	go test -run BenchmarkRules -bench BenchmarkRules -benchmem -memprofile memprofile.out -cpuprofile cpuprofile.out

compare:
	go test -run Benchmark* -bench Benchmark*

lint:
	golint | grep -v "don't use an underscore in package name" | grep -v "imported but not used" || true
