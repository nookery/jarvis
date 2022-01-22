/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// databaseCmd represents the database command
var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "数据库相关操作",
	Long:  `数据库相关操作`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("database called")
		fmt.Println(args)
	},
}

func init() {
	rootCmd.AddCommand(databaseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	databaseCmd.PersistentFlags().String("host", "", "数据库地址")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	databaseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
