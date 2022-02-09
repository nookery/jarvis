package site

import (
	"net/url"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"jarvis/cmd/bt/utils"
)

var delete = &cobra.Command{
	Use:   "delete",
	Short: "删除网站",
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		key, _ := cmd.Flags().GetString("key")
		link := host + "/site?action=DeleteSite"

		color.Infoln("地址：" + host)
		color.Infoln("密钥：" + key)

		result := utils.Post(link, utils.PatchSign(key, url.Values{
			"webname": {"test.api4.top"},
			"id":      {"10"},
		}))
		color.Infoln(result)
	},
	Args: nil,
}
