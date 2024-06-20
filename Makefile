.PHONY: run
run:
	go install
	(cd testdata && mocka)

.PHONY: fmt
fmt:
	goimports -w -local "github.com/karamaru-alpha/mocka" ./internal ./testdata
	go fmt ./...
