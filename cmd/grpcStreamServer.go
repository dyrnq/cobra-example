/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	//"context"

	"log"
	"net"
	"io"
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


func (s streamServer) Channel(stream pb.StreamService_ChannelServer) error {
	var total int32
 
	for {
		port, err := stream.Recv()

		if err == io.EOF {
			return stream.Send(&pb.Response{Total: total,});
		}
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

		lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d",viper.GetString("grpc.address") , viper.GetInt("grpc.port")))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		pb.RegisterStreamServiceServer(s, &streamServer{})
		log.Printf("server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}

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
	viper.BindPFlags(grpcStreamServerCmd.Flags())



}













