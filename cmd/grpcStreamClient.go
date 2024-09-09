/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	// "io"
	"log"
	"time"
	"math/big"
	"google.golang.org/grpc"
	pb "github.com/dyrnq/cobra-example/pkg/grpc/stream"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		conn, err := grpc.NewClient(viper.GetString("grpc.server"), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
	

	
		c := pb.NewStreamServiceClient(conn)
		
		stream, err := c.Channel(ctx)
		if err != nil {
			log.Printf("create stream: %w", err)
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
				}
				if resp != nil {
					log.Printf("Received message from server: %v", resp.Total)
				}
			}


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
	
	viper.BindPFlags(grpcStreamClientCmd.Flags())
}
