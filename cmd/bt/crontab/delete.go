package crontab

import (
	"fmt"
	"jarvis/cmd/bt/utils"
	"net/url"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var delete = &cobra.Command{
	Use:   "delete",
	Short: "删除crontab",
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		key, _ := cmd.Flags().GetString("key")
		name, _ := cmd.Flags().GetString("name")
		link := host + "/crontab?action=DelCrontab"
		var id int

		items, err := Get(host, key)

		if err != nil {
			color.Errorln(err.Error())

			return
		}

		for _, item := range items {
			if item.Name == name {
				id = item.Id

				break
			}
		}

		if id == 0 {
			color.Errorln("找不到相关crontab")
			return
		}

		color.Blueln("相关Crontab的ID是：", id)

		result := utils.Post(link, utils.PatchSign(key, url.Values{
			"id": {fmt.Sprint(id)},
		}))
		color.Infoln(result)
	},
}

func init() {
	delete.Flags().String("name", "", color.Blue.Render("名称"))
	delete.MarkFlagRequired("name")
}
