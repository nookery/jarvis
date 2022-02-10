package site

import (
	"jarvis/cmd/bt/utils"
	"net/url"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
)

var types = &gcli.Command{
	Name: "types",
	Desc: "展示网站分类",
	Func: func(cmd *gcli.Command, args []string) error {
		host := cmd.Arg("host").String()
		key := cmd.Arg("key").String()
		link := host + "/site?action=get_site_types"

		color.Infoln("地址：" + host)
		color.Infoln("密钥：" + key)

		result := utils.Post(link, utils.PatchSign(key, url.Values{}))
		color.Infoln(result)

		return nil
	},
}
