all: build

build:
	go build -o bin/ ./cmd/...

clean:
	$(RM) -r bin

.PHONY: all build clean
