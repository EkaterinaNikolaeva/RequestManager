BINARY_NAME=bin/manager

build:
	go build -o ${BINARY_NAME} ./cmd

run: build
	./${BINARY_NAME} $(CONFIG)

test:
	go test ./... -cover

clean:
	rm -r bin

