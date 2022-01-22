/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "数据库相关操作",
}

var create = &cobra.Command{
	Use:   "create",
	Short: "新建数据库",
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		name, _ := cmd.Flags().GetString("name")
		if host == "" {
			host = "127.0.0.1"
		}
		if username == "" {
			username = "root"
		}
		if password == "" {
			password = "root"
		}

		color.Infoln("地址：" + host)
		color.Infoln("用户：" + username)
		color.Infoln("密码：" + password)
		color.Infoln("新建：" + name + "\r\n")

		// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
		dsn := username + ":" + password + "@tcp(" + host + ":3306)/www_api4_top?charset=utf8mb4&parseTime=True&loc=Local"
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			color.Infoln("连接数据库失败")
		} else {
			_, err := db.Exec("CREATE DATABASE IF NOT EXISTS " + name)
			if err != nil {
				color.Errorf("Error %s when creating DB\n", err)
				return
			}
			color.Infoln("成功")
		}
	},
}

var show = &cobra.Command{
	Use:   "show",
	Short: "展示数据库列表",
	Long:  `展示数据库列表`,
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

		color.Infoln("数据库地址是：" + host)
		color.Infoln("用户名是：" + username)
		color.Infoln("密码是：" + password)

		// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
		dsn := username + ":" + password + "@tcp(" + host + ":3306)?charset=utf8mb4&parseTime=True&loc=Local"
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			color.Infoln("连接数据库失败")
		} else {
			color.Infoln("连接数据库成功")
			res, err := db.Exec("SHOW DATABASES")
			if err != nil {
				color.Infoln("Error %s when creating DB\n", err)
				return
			}
			no, err := res.RowsAffected()
			if err != nil {
				color.Infoln("Error %s when fetching rows", err)
				return
			}
			color.Infoln("rows affected %d\n", no)
		}
	},
}

func init() {
	rootCmd.AddCommand(databaseCmd)
	databaseCmd.AddCommand(create)
	databaseCmd.AddCommand(show)

	databaseCmd.PersistentFlags().String("host", "", "数据库地址")
	databaseCmd.PersistentFlags().String("username", "", "用户名")
	databaseCmd.PersistentFlags().String("password", "", "密码")

	create.Flags().String("name", "", "要新建的数据库名称")
}
