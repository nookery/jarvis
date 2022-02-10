package bt

import (
	"errors"
	"jarvis/cmd/bt/site"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var BtCmd = &cobra.Command{
	Use:   "bt",
	Short: "宝塔相关操作",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		host, _ := cmd.Flags().GetString("host")
		key, _ := cmd.Flags().GetString("key")

		color.Infoln("地址：" + host)
		color.Infoln("密钥：" + key + "\r\n")

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
	BtCmd.PersistentFlags().String("host", "http://127.0.0.1", "宝塔地址")
	BtCmd.PersistentFlags().StringP("key", "k", "", "宝塔密钥")
}
