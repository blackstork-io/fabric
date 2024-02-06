/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func data(target string) {
	fmt.Println("data called on", target)
}

// dataCmd represents the data command
var dataCmd = &cobra.Command{
	Use:   "data TARGET",
	Short: "Execute a single data block",
	Long:  `Execute the data block and print out prettified JSON to stdout`,
	Run: func(cmd *cobra.Command, args []string) {
		data(args[0])
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(dataCmd)

	dataCmd.SetUsageTemplate(UsageTemplate(
		[2]string{"TARGET", "a path to the data block to be executed. Data block must be inside of a document, so the path would look lile 'document.<doc-name>.data.<plugin-name>.<data-name>'"},
	))
}
