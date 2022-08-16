package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var create = &cobra.Command{
	Use:   "create",
	Short: "创建数据库",
	Long:  color.Success.Render("\r\n创建数据库。"),
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

		if name == "" {
			color.Warnln("请输入要新建的数据库名称")
			return
		}

		color.Infoln("地址：" + host)
		color.Infoln("用户：" + username)
		color.Infoln("密码：" + password)
		color.Infoln("新建：" + name + "\r\n")

		// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
		dsn := username + ":" + password + "@tcp(" + host + ":3306)/?charset=utf8mb4&parseTime=True&loc=Local"
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

func init() {
	create.Flags().String("name", "", "要新建的数据库名称")
}
