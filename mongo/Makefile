check:
	@go vet ./...

test:
	@go clean -testcache
	@go test -race ./...
	@go test -timeout 10s -run TestLeak
	@go test -coverpkg=./... -coverprofile=coverage.data -timeout=5s ./...
	@go tool cover -html=coverage.data -o coverage.html

