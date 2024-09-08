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

SERVER_ADDRESS=0.0.0.0 ./main serve --server.address 127.0.0.1 --server.port=8899


SERVER_ADDRESS=0.0.0.0 ./main serve --server.address 127.0.0.1 --server.port=35566 --config ./config-demo.yaml

```

Priority

```bash
--args > ENV > config.yaml
```