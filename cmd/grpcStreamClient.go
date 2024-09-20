/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"crypto/tls"
	"crypto/x509"
	"os"	
	// "io"
	"log"
	"time"
	"math/big"
	"google.golang.org/grpc"
	"runtime"
	pb "github.com/dyrnq/cobra-example/pkg/grpc/stream"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// grpcStreamClientCmd represents the grpcStreamClient command
var grpcStreamClientCmd = &cobra.Command{
	Use:   "client",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("grpcStreamClient called")

        certFile := viper.GetString("grpc.tls.cert")
        keyFile := viper.GetString("grpc.tls.key")
        caFile := viper.GetString("grpc.tls.ca") // 新增：读取CA证书路径
		var conn *grpc.ClientConn
		var err error
		for {


			if certFile!= "" && keyFile!= "" && caFile!= "" {
                cert, err := tls.LoadX509KeyPair(certFile, keyFile)
                if err!= nil {
                        log.Fatalf("Failed to load TLS certificate and key: %v", err)
                }

                // 创建TLS配置
                tlsConfig := &tls.Config{
                        Certificates: []tls.Certificate{cert},
                        // 使用CA证书进行验证
                        RootCAs:    x509.NewCertPool(),
                        InsecureSkipVerify: false, 
                }

                caCert, err := os.ReadFile(caFile)
                if err!= nil {
                        log.Fatalf("Failed to load CA certificate: %v", err)
                }

                if ok := tlsConfig.RootCAs.AppendCertsFromPEM(caCert);!ok {
                        log.Fatalf("Failed to append CA certificate")
                }

                // 使用TLS配置创建gRPC连接
                conn, err = grpc.Dial(viper.GetString("grpc.server"), grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))

        } else {
                // 如果没有提供TLS证书和密钥，则使用不安全的连接
                conn, err = grpc.Dial(viper.GetString("grpc.server"), grpc.WithTransportCredentials(insecure.NewCredentials()))

        }

		//conn, err := grpc.NewClient(viper.GetString("grpc.server"), grpc.WithTransportCredentials(insecure.NewCredentials()))



		if err != nil {
			log.Fatalf("did not connect: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}
		defer conn.Close()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
	

	
		c := pb.NewStreamServiceClient(conn)
		
		stream, err := c.Channel(ctx)
		if err != nil {
			log.Printf("create stream: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}


		// // I am hardcoding this here but you should not!
		// requests := []*pb.Request{
		// 	{Id: "id-1", Name: "name-1"},
		// 	{Id: "id-2", Name: "name-2"},
		// 	{Id: "id-3", Name: "name-3"},
		// 	{Id: "id-4", Name: "name-4"},
		// }
	
		// for _, request := range requests {
		// 	if err := stream.Send(request); err != nil {
		// 		log.Printf("send stream: %w", err)
		// 	}
		// }
		// response, err := stream.CloseAndRecv()
		// if err != nil {
		// 	log.Printf("close and receive: %w", err)
		// }
	
		// log.Printf("%+v\n", response)


		    // 定期发送消息给服务端
			go func() {

				// 创建一个新的 big.Int 并初始化为 0
				bigInt := big.NewInt(0)

				for {
					bigInt.Add(bigInt, big.NewInt(1))
					err := stream.Send(&pb.Request{
						Name: fmt.Sprintf("Name-%s", bigInt.String() ),
						Id: fmt.Sprintf("Id-%s", bigInt.String() ),
					})
					if err != nil {
						log.Printf("Failed to send message: %v", err)
						//return
						runtime.Goexit() // 立即终止当前协程
					}
					time.Sleep(1 * time.Second)
				}


			}()

	

	
			// 接收服务端的响应
			for {
				resp, err := stream.Recv()
				// if err == io.EOF {
				// 	log.Printf("Received message from server: %v", resp.Total)
				// }
				if err != nil {
					log.Printf("Failed to send message: %v", err)
					break;
				}
				if resp != nil {
					log.Printf("Received message from server: %v", resp.Total)
				}
			}

		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("grpc.server", cmd.Flags().Lookup("grpc.server"))
		viper.BindPFlag("grpc.tls.cert", cmd.Flags().Lookup("grpc.tls.cert"))
		viper.BindPFlag("grpc.tls.key", cmd.Flags().Lookup("grpc.tls.key"))
		viper.BindPFlag("grpc.tls.ca", cmd.Flags().Lookup("grpc.tls.ca"))	
   },
}

func init() {
	grpcStreamCmd.AddCommand(grpcStreamClientCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// grpcStreamClientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// grpcStreamClientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	grpcStreamClientCmd.Flags().StringP("grpc.server", "", "127.0.0.1:50053", "Server address")
	grpcStreamClientCmd.Flags().StringP("grpc.tls.cert", "", "", "grpc.tls.cert")
	grpcStreamClientCmd.Flags().StringP("grpc.tls.key", "", "", "grpc.tls.key")
	grpcStreamClientCmd.Flags().StringP("grpc.tls.ca", "", "", "grpc.tls.ca")

	// viper.BindPFlags(grpcStreamClientCmd.Flags())
}
