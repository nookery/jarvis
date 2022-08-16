package site

import (
	"jarvis/cmd/bt/utils"
	"net/url"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var types = &cobra.Command{
	Use:   "types",
	Short: "展示网站分类",
	Long:  color.Success.Render("展示网站分类"),
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		key, _ := cmd.Flags().GetString("key")
		link := host + "/site?action=get_site_types"

		result := utils.Post(link, utils.PatchSign(key, url.Values{}))
		color.Infoln(result)
	},
}
