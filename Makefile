BINARY_NAME=bin/manager

build:
	go build -o ${BINARY_NAME} ./cmd

run: build
	./${BINARY_NAME} $(CONFIG)

test:
	go test -v -coverprofile cover.out ./... && go tool cover -func cover.out

clean:
	rm -r bin

