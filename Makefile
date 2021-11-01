tidy:
	go mod tidy

test:
	go test -race -v -count=1 ./...