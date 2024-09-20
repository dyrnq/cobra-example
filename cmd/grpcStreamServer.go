/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"crypto/tls"
	"crypto/x509"
	//"context"

	"log"
	"net"
	"os"
	"google.golang.org/grpc"
	pb "github.com/dyrnq/cobra-example/pkg/grpc/stream"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)


type streamServer struct{
	pb.UnimplementedStreamServiceServer
}


// func (s streamServer) Channel(stream pb.StreamService_ChannelServer) error {
// 	var total int32
 
// 	for {
// 		port, err := stream.Recv()
// 		if err == io.EOF {
// 			return stream.SendAndClose(&pb.Response{
// 				Total: total,
// 			})
// 		}
// 		if err != nil {
// 			return err
// 		}
 
// 		total++			

// 		log.Printf("%+v\n", port)
// 	}
// }

var total int32
func (s streamServer) Channel(stream pb.StreamService_ChannelServer) error {
	//var total int32
 
	for {
		port, err := stream.Recv()
		stream.Send(&pb.Response{Total: total,});
		// if err == io.EOF {
		// 	return stream.Send(&pb.Response{Total: total,});
		// }
		if err != nil {
			return err
		}
 
		total++
		log.Printf("%+v\n", port)
	}
}



// func (s server) FetchResponse(in *pb.Request, srv pb.StreamService_FetchResponseServer) error {

// 	log.Printf("fetch response for id : %d", in.Id)

// 	var wg sync.WaitGroup
// 	for i := 0; i < 5; i++ {
// 		wg.Add(1)
// 		go func(count int64) {
// 			defer wg.Done()
// 			time.Sleep(time.Duration(count) * time.Second)
// 			resp := pb.Response{Result: fmt.Sprintf("Request #%d For Id:%d", count, in.Id)}
// 			if err := srv.Send(&resp); err != nil {
// 				log.Printf("send error %v", err)
// 			}
// 			log.Printf("finishing request number : %d", count)
// 		}(int64(i))
// 	}

// 	wg.Wait()
// 	return nil
// }

// grpcStreamServerCmd represents the grpcStreamServer command
var grpcStreamServerCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("grpcStreamServer called")
		var lis net.Listener
		var err error
		var msg string

		certFile := viper.GetString("grpc.tls.cert")
		keyFile := viper.GetString("grpc.tls.key")
		if certFile!= "" && keyFile!= "" {

			// 加载 CA 证书
			caCertFile := viper.GetString("grpc.tls.ca")

			// 检查是否需要客户端认证
			clientAuth := viper.GetBool("grpc.tls.client.verify")

			cert, err := tls.LoadX509KeyPair(certFile, keyFile)
			if err!= nil {
				log.Fatalf("failed to load TLS certificate and key: %v", err)
			}

			tlsConfig := &tls.Config{
				Certificates: []tls.Certificate{cert},
				// 根据命令行参数决定是否需要客户端认证
				ClientAuth: tls.NoClientCert,
			}

			// 如果需要客户端认证，加载 CA 证书并设置 ClientAuth
			if clientAuth {
				caCert, err := os.ReadFile(caCertFile)
				if err!= nil {
					log.Fatalf("failed to load CA certificate: %v", err)
				}
				caCertPool := x509.NewCertPool()
				caCertPool.AppendCertsFromPEM(caCert)
				tlsConfig.ClientCAs = caCertPool
				tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
			}
			// 使用 tls.Listen 创建安全的监听器
			lis, err = tls.Listen("tcp", fmt.Sprintf("%s:%d", viper.GetString("grpc.address"), viper.GetInt("grpc.port")), tlsConfig)
			msg = "TLS"
		}else{
			lis, err = net.Listen("tcp", fmt.Sprintf("%s:%d",viper.GetString("grpc.address") , viper.GetInt("grpc.port")))
			msg = "plain TCP"
		}
		
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		pb.RegisterStreamServiceServer(s, &streamServer{})
		log.Printf("server listening at %s %v", msg, lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}

	},
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("grpc.address", cmd.Flags().Lookup("grpc.address"))
		viper.BindPFlag("grpc.port", cmd.Flags().Lookup("grpc.port"))
		viper.BindPFlag("grpc.tls.cert", cmd.Flags().Lookup("grpc.tls.cert"))
		viper.BindPFlag("grpc.tls.key", cmd.Flags().Lookup("grpc.tls.key"))
		// 绑定新添加的命令行参数
		viper.BindPFlag("grpc.tls.ca", cmd.Flags().Lookup("grpc.tls.ca"))
		viper.BindPFlag("grpc.tls.client.verify", cmd.Flags().Lookup("grpc.tls.client.verify"))		
   },
}

func init() {
	grpcStreamCmd.AddCommand(grpcStreamServerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// grpcStreamServerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// grpcStreamServerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	grpcStreamServerCmd.Flags().StringP("grpc.address", "", "", "Server address")
	grpcStreamServerCmd.Flags().IntP("grpc.port", "", 50053, "Server port")

	// 在不同的子命令中使用相同的参数名称确实会导致 Viper 混淆。这是因为 Viper 使用参数名称作为键来存储和读取配置值。
	//viper.BindPFlags(grpcStreamServerCmd.Flags())

	// 添加新的命令行参数
	grpcStreamServerCmd.Flags().String("grpc.tls.cert", "", "Path to the TLS certificate file")
	grpcStreamServerCmd.Flags().String("grpc.tls.key", "", "Path to the TLS key file")
	grpcStreamServerCmd.Flags().String("grpc.tls.ca", "", "Path to the CA certificate file for client authentication")
	grpcStreamServerCmd.Flags().Bool("grpc.tls.client.verify", false, "Enable client verify")

}













