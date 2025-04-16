/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package net

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)


var urlPath string
// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var client = http.Client {
			Transport:&http.Transport{
				Dial:net.Dialer{Timeout:2 * }
			}
		}
	},
}

func init() {
	pingCmd.Flags().StringVarP(&urlPath,"url","u","","The url to ping")
	if err := pingCmd.MarkFlagRequired("url"); err != nil{
		fmt.Println(err)
	}
	NetCmd.AddCommand(pingCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
