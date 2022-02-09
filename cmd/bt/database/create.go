package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
)

var createCmd = &gcli.Command{
	Name: "create",
	Desc: "新建数据库",
	Func: create,
}

func create(cmd *gcli.Command, args []string) error {
	cmd.AddArg("host", "host")
	cmd.AddArg("username", "username")
	cmd.AddArg("password", "password")
	cmd.AddArg("name", "name")

	host := cmd.Arg("host").String()
	username := cmd.Arg("username").String()
	password := cmd.Arg("password").String()
	name := cmd.Arg("name").String()
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
		return nil
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
			return nil
		}
		color.Infoln("成功")
	}

	return nil
}
