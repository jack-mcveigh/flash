all: build

build:
	go build -o bin/ ./cmd/...

clean:
	$(RM) -r bin *.out

coverage-report:
	go test -coverprofile cover.out ./...
	go tool cover -html=cover.out

test:
	go test -v -cover ./pkg/...

.PHONY: all build clean
