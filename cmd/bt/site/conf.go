package site

import (
	"errors"
	"io/ioutil"
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
		file, _ := cmd.Flags().GetString("file")

		color.Infoln("名称：" + name + "\r\n")
		color.Infoln("文件：" + file + "\r\n")

		if name == "" {
			return errors.New(color.Red.Renderln("请输入网站名称") + "\r\n")
		}

		if file == "" {
			return errors.New(color.Red.Renderln("请输入配置文件的路径") + "\r\n")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		key, _ := cmd.Flags().GetString("key")
		link := host + "/files?action=SaveFileBody"

		name, _ := cmd.Flags().GetString("name")
		file, _ := cmd.Flags().GetString("file")

		path := "/www/server/panel/vhost/nginx/" + name + ".conf"

		data, err := ioutil.ReadFile(file)
		if err != nil {
			color.Red.Println("读取配置文件失败")
			return
		}

		result := utils.Post(link, utils.PatchSign(key, url.Values{
			"path":     {path},
			"data":     {string(data)},
			"encoding": {"utf-8"},
		}))
		color.Infoln(result)
	},
}

func init() {
	Conf.Flags().StringP("name", "n", "", "网站名称")
	Conf.Flags().StringP("file", "f", "", "配置文件路径")
	Conf.MarkFlagRequired("file")
	Conf.MarkFlagRequired("name")
}
