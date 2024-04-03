/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net"
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "build-your-own-curl",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		u, err := url.Parse(args[0])
		if err != nil {
			panic(err)
		}

		host := u.Hostname()
		port := u.Port()
		path := u.Path

		if port == "" {
			port = "80"
		}

		println("Host:", host)
		println("Port:", port)
		println("Path:", path)

		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
		if err != nil {
			panic(err)
		}

		defer conn.Close()

		fmt.Fprintf(conn, "GET %s HTTP/1.0\r\nHost: %s\r\n\r\n", path, host)

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(buf[:n]))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.build-your-own-curl.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
