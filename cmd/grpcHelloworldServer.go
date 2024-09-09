/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"context"

	"log"
	"net"

	"google.golang.org/grpc"
	pb "github.com/dyrnq/cobra-example/pkg/grpc/helloworld"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(_ context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

// grpcHelloworldServerCmd represents the server command
var grpcHelloworldServerCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("server called")


		lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d",viper.GetString("grpc.address") , viper.GetInt("grpc.port")))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		pb.RegisterGreeterServer(s, &server{})
		log.Printf("server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}


	},
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("grpc.address", cmd.Flags().Lookup("grpc.address"))
		viper.BindPFlag("grpc.port", cmd.Flags().Lookup("grpc.port"))
   },
}

func init() {
	grpcHelloworldCmd.AddCommand(grpcHelloworldServerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// grpcHelloworldServerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// grpcHelloworldServerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	grpcHelloworldServerCmd.Flags().StringP("grpc.address", "", "", "Server address")
	grpcHelloworldServerCmd.Flags().IntP("grpc.port", "", 50051, "Server port")

	
	// viper.BindPFlags(grpcHelloworldServerCmd.Flags())

}
