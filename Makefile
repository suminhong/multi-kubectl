BINARY_NAME=mk

build:
	go build -o $(BINARY_NAME) main.go

install: build
	mv $(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

test:
	go test ./...

clean:
	rm -f $(BINARY_NAME)
