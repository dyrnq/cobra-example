tidy:
	GOPROXY=https://goproxy.io,direct go mod tidy -v
build:
	GOPROXY=https://goproxy.io,direct CGO_ENABLED=0 go build -v -o main main.go

run: tidy build
	./main

