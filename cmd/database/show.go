package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

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
		dsn := username + ":" + password + "@tcp(" + host + ":3306)/?charset=utf8mb4&parseTime=True&loc=Local"
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			color.Infoln("连接数据库失败", err)
		} else {
			res, err := db.Query("SHOW DATABASES")
			if err != nil {
				color.Infoln("Error %s when creating DB\n", err)
				return
			}

			defer res.Close()

			color.Infoln("数据库列表：")
			name := ""
			for res.Next() {
				err := res.Scan(&name)
				if err != nil {
					color.Errorf(err.Error())
				}
				color.Println("  " + name)
			}
		}
	},
}
