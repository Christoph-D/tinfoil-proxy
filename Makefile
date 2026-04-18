.PHONY: build clean

BINARY := tinfoil-proxy

build:
	go build -ldflags="-extldflags=-static" -o $(BINARY) .

clean:
	rm -f $(BINARY)
