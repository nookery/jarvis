package database

import (
	"errors"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var DatabaseCmd = &cobra.Command{
	Use:   "database",
	Short: color.Blue.Render("数据库相关操作"),
	Long:  color.Success.Render("\r\n数据库相关操作"),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		host, _ := cmd.Flags().GetString("host")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")

		color.Infoln("地址：" + host)
		color.Infoln("用户：" + username)
		color.Infoln("密码：" + password)

		if host == "" {
			return errors.New(color.Error.Renderln("数据库地址") + "\r\n")
		}

		if username == "" {
			return errors.New(color.Error.Renderln("请输入用户名") + "\r\n")
		}

		if password == "" {
			return errors.New(color.Error.Renderln("请输入密码") + "\r\n")
		}

		return nil
	},
}

func init() {
	DatabaseCmd.AddCommand(create)
	DatabaseCmd.AddCommand(show)
	DatabaseCmd.PersistentFlags().String("host", "http://127.0.0.1:3306", "数据库地址")
	DatabaseCmd.PersistentFlags().StringP("username", "u", "root", "数据库用户名")
	DatabaseCmd.PersistentFlags().StringP("password", "p", "root", "数据库密码")
}
