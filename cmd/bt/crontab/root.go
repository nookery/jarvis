package crontab

import (
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var CrontabCmd = &cobra.Command{
	Use:   "crontab",
	Short: color.Blue.Render("crontab相关操作"),
	Long:  color.Success.Render("\r\ncrontab相关操作"),
}

func init() {
	CrontabCmd.AddCommand(get)
	CrontabCmd.AddCommand(create)
	CrontabCmd.AddCommand(delete)
}
