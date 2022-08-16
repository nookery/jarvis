package database

import (
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var DatabaseCmd = &cobra.Command{
	Use:   "database",
	Short: color.Blue.Render("数据库相关操作"),
	Long:  color.Success.Render("\r\n数据库相关操作"),
}

func init() {
	DatabaseCmd.AddCommand(Create)
}
