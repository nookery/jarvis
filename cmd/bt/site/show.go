package site

import (
	"jarvis/cmd/bt/utils"
	"net/url"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var show = &cobra.Command{
	Use:   "show",
	Short: "展示网站列表",
	Long:  color.Success.Render("展示网站列表"),
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		key, _ := cmd.Flags().GetString("key")
		link := host + "/data?action=getData&table=sites"

		result := utils.Post(link, utils.PatchSign(key, url.Values{}))
		color.Infoln(result)
	},
}
