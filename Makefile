BINARY_NAME=translate

build:
	go build -o ${BINARY_NAME} translate.go specs.go
	go build -o ${BINARY_NAME}-cli translate-cli.go specs.go

run: build
	./${BINARY_NAME}

clean:
	 go clean