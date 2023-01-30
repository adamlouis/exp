
build: *.go
	mkdir -p build
	go build -o build/lox

run:
	go run .

lint:
	go fmt

gen:
	go run genast/main.go .
	go fmt *.go
