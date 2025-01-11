# cobra-example


this project is a POC using cobra, viper, pflag.


## init

```bash
go mod init github.com/dyrnq/cobra-example

cobra-cli init --viper
cobra-cli add serve --viper
cobra-cli add version --viper
GOPROXY=https://goproxy.io,direct go mod tidy
```

## test


```bash
make run
```

```bash
## test case
./main serve

./main serve --server.address 127.0.0.1

SERVER_ADDRESS=0.0.0.0 ./main serve

SERVER_ADDRESS=0.0.0.0 SERVER_PORT=7777 ./main serve --config ./config-demo.yaml

SERVER_ADDRESS=0.0.0.0 ./main serve --server.address 127.0.0.1 --server.port=8899


SERVER_ADDRESS=0.0.0.0 ./main serve --server.address 127.0.0.1 --server.port=35566 --config ./config-demo.yaml

```

Priority

```bash
--args > ENV > config.yaml
```

## grpc

Add <https://github.com/grpc/grpc-go/tree/master/examples/helloworld> to subcomand

```bash
protoc \
--go_out=. \
--go_opt=paths=source_relative \
--go-grpc_out=. \
--go-grpc_opt=paths=source_relative \
pkg/grpc/helloworld/*.proto
```

> bash: protoc: command not found

```bash
# https://github.com/protocolbuffers/protobuf/releases/download/v28.0/protoc-28.0-linux-x86_64.zip
PROTOC_VER=28.0
if [ "$(uname -m)" = "x86_64" ]; then arch="amd64"; fi
if [ "$(uname -m)" = "aarch64" ]; then arch="arm64"; fi
PROTOC_ZIP="protoc-${PROTOC_VER}-linux-${arch}.zip"
url="https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VER}/${PROTOC_ZIP}";
url=${url/github.com/mirror.ghproxy.com/github.com}
echo "${url}"
curl -O -fL -# --retry 10 ${url}
unzip -o $PROTOC_ZIP -d /usr/local bin/protoc
unzip -o $PROTOC_ZIP -d /usr/local 'include/*'
rm -f $PROTOC_ZIP
```

> protoc-gen-go: program not found or is not executable
> Please specify a program using absolute path or make sure the program is available in your PATH system variable
> --go_out: protoc-gen-go: Plugin failed with status code 1.

```bash

$ GOPROXY=https://goproxy.io,direct go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go: downloading google.golang.org/protobuf v1.34.2

$ GOPROXY=https://goproxy.io,direct go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go: downloading google.golang.org/grpc v1.66.0
go: downloading google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.5.1
go: downloading google.golang.org/protobuf v1.34.1

```

> pkg/grpc/helloworld/helloworld_grpc.pb.go:33:16: undefined: grpc.SupportPackageIsVersion9

```bash
$ GOPROXY=https://goproxy.io,direct go get -u google.golang.org/grpc@latest
go: downloading golang.org/x/net v0.26.0
go: downloading golang.org/x/sys v0.21.0
go: downloading google.golang.org/genproto/googleapis/rpc v0.0.0-20240604185151-ef581f913117
go: downloading golang.org/x/net v0.29.0
go: downloading google.golang.org/genproto/googleapis/rpc v0.0.0-20240903143218-8af14fe29dc1
go: downloading google.golang.org/genproto v0.0.0-20240903143218-8af14fe29dc1
go: downloading golang.org/x/text v0.16.0
go: downloading golang.org/x/text v0.18.0
go: upgraded golang.org/x/net v0.23.0 => v0.29.0
go: upgraded golang.org/x/sys v0.18.0 => v0.25.0
go: upgraded golang.org/x/text v0.14.0 => v0.18.0
go: upgraded google.golang.org/genproto/googleapis/rpc v0.0.0-20240314234333-6e1732d8331c => v0.0.0-20240903143218-8af14fe29dc1
go: upgraded google.golang.org/grpc v1.62.1 => v1.66.0
go: upgraded google.golang.org/protobuf v1.33.0 => v1.34.2
```

```bash
cobra-cli add grpc-helloworld --parent rootCmd --viper
cobra-cli add server --parent grpcHelloworldCmd --viper
cobra-cli add client --parent grpcHelloworldCmd --viper
```

```bash
./main grpc-helloworld server --grpc.address 0.0.0.0 --grpc.port 50051

./main grpc-helloworld client --grpc.server 127.0.0.1:50051 --msg 你好
```


## grpc-stream

- <https://github.com/pramonow/go-grpc-server-streaming-example>

```bash
cobra-cli add grpc-stream           --parent rootCmd        --viper
cobra-cli add grpc-stream-server    --parent grpcStreamCmd  --viper
cobra-cli add grpc-stream-client    --parent grpcStreamCmd  --viper


./main grpc-stream server

./main grpc-stream client
```