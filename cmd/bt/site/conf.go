package site

import (
	"errors"
	"jarvis/cmd/bt/utils"
	"net/url"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var Conf = &cobra.Command{
	Use:   "conf",
	Short: color.Blue.Render("保存网站配置"),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		content, _ := cmd.Flags().GetString("content")

		color.Infoln("名称：" + name + "\r\n")

		if name == "" {
			return errors.New(color.Red.Renderln("请输入网站名称") + "\r\n")
		}

		if content == "" {
			return errors.New(color.Red.Renderln("请提供配置文件的内容") + "\r\n")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		key, _ := cmd.Flags().GetString("key")
		link := host + "/files?action=SaveFileBody"

		name, _ := cmd.Flags().GetString("name")
		content, _ := cmd.Flags().GetString("content")

		path := "/www/server/panel/vhost/nginx/" + name + ".conf"

		result := utils.Post(link, utils.PatchSign(key, url.Values{
			"path":     {path},
			"data":     {content},
			"encoding": {"utf-8"},
		}))
		color.Infoln(result)
	},
}

func init() {
	Conf.Flags().StringP("name", "n", "", color.Blue.Render("网站名称"))
	Conf.Flags().StringP("content", "c", "", color.Blue.Render("配置文件内容，内容较多时可以这样写：-c=\"$(cat sample.conf)\""))
	Conf.MarkFlagRequired("name")
	Conf.MarkFlagRequired("content")
}
