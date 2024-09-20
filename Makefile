tidy:
	GOPROXY=https://goproxy.io,direct go mod tidy -v
build:
	GOPROXY=https://goproxy.io,direct CGO_ENABLED=0 go build -v -o main main.go

run: tidy build
	./main

tls-server:
	./main grpc-stream server --grpc.tls.ca certs/ca.crt --grpc.tls.cert certs/server.crt --grpc.tls.key certs/server.key --grpc.tls.client.verify

tls-client:
	./main grpc-stream client --grpc.tls.ca certs/ca.crt --grpc.tls.cert certs/client.crt --grpc.tls.key certs/client.key