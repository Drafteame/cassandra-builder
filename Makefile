install:
	go mod download

test:
	go test -v -race -cover ./... -coverprofile=coverage.out

html_cover:
	go tool cover -html=coverage.out

func_cover:
	go tool cover -func=coverage.out

lint:
	golangci-lint run ./...