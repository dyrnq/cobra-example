/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"math"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"	
)

// pythagoreanTripleCmd represents the pythagoreanTriple command
var pythagoreanTripleCmd = &cobra.Command{
	Use:   "pt",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		impl1(viper.GetInt("max"))
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("max", cmd.Flags().Lookup("max"))
   },
}

func impl1(max int){
	var count int
	for a := 1; a < max*2; a++ {
		for b := a + 1; b < max*2; b++ {
			c := int(math.Sqrt(float64(a*a + b*b)))
			if a*a+b*b == c*c {

				fmt.Printf("%d^2 + %d^2 = %d^2\n", a, b, c )
				count++
			}
			if count >= 100 {
				return
			}
		}
	}
	//fmt.Printf("%d\n", count)
}





func init() {
	rootCmd.AddCommand(pythagoreanTripleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pythagoreanTripleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pythagoreanTripleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	pythagoreanTripleCmd.Flags().IntP("max", "", 100, "max")
}
