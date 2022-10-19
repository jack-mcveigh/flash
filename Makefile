all: build

build:
	go build -o bin/ ./cmd/...

clean:
	$(RM) -r bin

test:
	go test ./...

.PHONY: all build clean
