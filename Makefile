BINARY_NAME := guesswidth
SRCS := $(shell git ls-files '*.go')
LDFLAGS := "-X github.com/noborus/guesswidth.version=$(shell git describe --tags --abbrev=0 --always) -X github.com/noborus/guesswidth.revision=$(shell git rev-parse --short HEAD)"

all: build

test: $(SRCS)
	go test ./...

build: $(BINARY_NAME)

$(BINARY_NAME): $(SRCS)
	go build -ldflags $(LDFLAGS) -o $(BINARY_NAME) ./cmd/guesswidth

install:
	go install -ldflags $(LDFLAGS) ./cmd/guesswidth

clean:
	rm -f $(BINARY_NAME)

.PHONY: all test build install clean
