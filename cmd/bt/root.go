package bt

import (
	"errors"
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

		color.Infoln("地址：" + host)

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
	BtCmd.AddCommand(site.SiteCmd)
	BtCmd.PersistentFlags().String("host", "http://127.0.0.1", color.Blue.Render("宝塔地址"))
	BtCmd.PersistentFlags().StringP("key", "k", "", color.Blue.Render("宝塔密钥"))
}
