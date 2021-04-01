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

benchmarkCBEMarshal:
	$(info To see profile: go tool pprof cpuprofile.out)
	go test -run BenchmarkCBEMarshal -bench BenchmarkCBEMarshal -benchmem -memprofile memprofile.out -cpuprofile cpuprofile.out

benchmarkCTEMarshal:
	$(info To see profile: go tool pprof cpuprofile.out)
	go test -run BenchmarkCTEMarshal -bench BenchmarkCTEMarshal -benchmem -memprofile memprofile.out -cpuprofile cpuprofile.out

benchmarkJSONMarshal:
	$(info To see profile: go tool pprof cpuprofile.out)
	go test -run BenchmarkJSONMarshal -bench BenchmarkJSONMarshal -benchmem -memprofile memprofile.out -cpuprofile cpuprofile.out

benchmarkCBEUnmarshal:
	$(info To see profile: go tool pprof cpuprofile.out)
	go test -run BenchmarkCBEUnmarshalNoRules -bench BenchmarkCBEUnmarshalNoRules -benchmem -memprofile memprofile.out -cpuprofile cpuprofile.out

benchmarkCTEUnmarshal:
	$(info To see profile: go tool pprof cpuprofile.out)
	go test -run BenchmarkCTEUnmarshalNoRules -bench BenchmarkCTEUnmarshalNoRules -benchmem -memprofile memprofile.out -cpuprofile cpuprofile.out

benchmarkCBEUnmarshalRules:
	$(info To see profile: go tool pprof cpuprofile.out)
	go test -run BenchmarkCBEUnmarshalRules -bench BenchmarkCBEUnmarshalRules -benchmem -memprofile memprofile.out -cpuprofile cpuprofile.out

benchmarkJSONUnmarshal:
	$(info To see profile: go tool pprof cpuprofile.out)
	go test -run BenchmarkJSONUnmarshalNoRules -bench BenchmarkJSONUnmarshalNoRules -benchmem -memprofile memprofile.out -cpuprofile cpuprofile.out

benchmarkCTEDecode:
	$(info To see profile: go tool pprof cpuprofile.out)
	go test -run BenchmarkCTEDecode -bench BenchmarkCTEDecode -benchmem -memprofile memprofile.out -cpuprofile cpuprofile.out

benchmarkRules:
	$(info To see profile: go tool pprof cpuprofile.out)
	go test -run BenchmarkRules -bench BenchmarkRules -benchmem -memprofile memprofile.out -cpuprofile cpuprofile.out

benchmarkBuilder:
	$(info To see profile: go tool pprof cpuprofile.out)
	go test -run BenchmarkBuilder -bench BenchmarkBuilder -benchmem -memprofile memprofile.out -cpuprofile cpuprofile.out

benchmarkIterator:
	$(info To see profile: go tool pprof cpuprofile.out)
	go test -run BenchmarkIterator -bench BenchmarkIterator -benchmem -memprofile memprofile.out -cpuprofile cpuprofile.out

compare:
	go test -run Benchmark* -bench Benchmark*

lint:
	golint | grep -v "don't use an underscore in package name" | grep -v "imported but not used" || true
