export COVERAGE_PACKAGES=aggregator config fileParser flagsParser outputPrinter predictor processor

coverage:
	echo "mode: count" > coverage-all.out
	$(foreach pkg, $(COVERAGE_PACKAGES),\
					go test -cover -coverprofile=coverage.out -covermode=count ./$(pkg);\
					tail -n +2 coverage.out >> coverage-all.out;)
	go tool cover -func=coverage-all.out