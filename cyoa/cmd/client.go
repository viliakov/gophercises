/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: client,
}

func init() {
	rootCmd.AddCommand(clientCmd)

	clientCmd.Flags().String("server", "127.0.0.1:8080", "Server to get stories from")
}

func client(cmd *cobra.Command, args []string) {
	fmt.Println("client called")
}
