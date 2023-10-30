check:
	@go vet ./...

test:
	@go test -race ./...
	@go test -timeout 10s -run TestLeak

