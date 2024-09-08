/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("serve called")


		http.HandleFunc("/healthz", healthzHandler)
		http.HandleFunc("/", defaultHandler)
		info := fmt.Sprintf("%s:%d", ConfigVar.Server.Address, ConfigVar.Server.Port)
		// fmt.Println(info)
		// err := http.ListenAndServe(info, nil)
		// if err != nil {
		// 	fmt.Println("Failed to start server:", err)
		// }
		srv := &http.Server{Addr: info, Handler: nil}

		go func() {
			if err := srv.ListenAndServe(); err != nil {
				log.Println("Server error:", err)
			}
		}()
	
		log.Println("Server is starting on", srv.Addr)
	
		// 等待系统信号
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
	
		log.Println("Shutting down server...")
		srv.Close() // 优雅关闭服务器

	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	serveCmd.Flags().StringP("server.address", "", "", "Server address")
	serveCmd.Flags().IntP("server.port", "p", 8080, "Server port")

	
	viper.BindPFlags(serveCmd.Flags())


}



func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("hello"))
}