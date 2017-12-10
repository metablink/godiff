run:
	go run main.go

build-deps:
	godep save ./...

test:
	ginkgo ./test/...
