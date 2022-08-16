package bt

import (
	"errors"
	"jarvis/cmd/bt/crontab"
	"jarvis/cmd/bt/database"
	"jarvis/cmd/bt/site"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var BtCmd = &cobra.Command{
	Use:   "bt",
	Long:  color.Success.Render("\r\n宝塔管理工具。"),
	Short: color.Blue.Render("宝塔相关操作"),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		host, _ := cmd.Flags().GetString("host")
		key, _ := cmd.Flags().GetString("key")

		color.Blueln("\r\n宝塔API地址：" + host + "\r\n")

		if host == "" {
			return errors.New(color.Error.Renderln("请输入宝塔地址") + "\r\n")
		}

		if key == "" {
			return errors.New(color.Red.Renderln("请输入宝塔密钥") + "\r\n")
		}

		return nil
	},
}

func init() {
	BtCmd.AddCommand(Http)
	BtCmd.AddCommand(crontab.CrontabCmd)
	BtCmd.AddCommand(site.SiteCmd)
	BtCmd.AddCommand(database.DatabaseCmd)
	BtCmd.PersistentFlags().StringP("host", "s", "http://127.0.0.1:8888", color.Blue.Render("宝塔地址"))
	BtCmd.PersistentFlags().StringP("key", "k", "", color.Blue.Render("宝塔密钥"))
}
