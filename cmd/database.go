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
		host, _ := cmd.Flags().GetString("host")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		if host == "" {
			host = "127.0.0.1"
		}
		if username == "" {
			username = "root"
		}
		if password == "" {
			password = "root"
		}

		fmt.Println("数据库地址是：" + host)
		fmt.Println("用户名是：" + username)
		fmt.Println("密码是：" + password)
	},
}

func init() {
	rootCmd.AddCommand(databaseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	databaseCmd.PersistentFlags().String("host", "", "数据库地址")

	databaseCmd.PersistentFlags().String("username", "", "用户名")

	databaseCmd.PersistentFlags().String("password", "", "密码")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// databaseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
