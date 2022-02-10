package site

import (
	"jarvis/cmd/bt/utils"
	"net/url"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
)

var show = &gcli.Command{
	Name: "show",
	Desc: "展示网站列表",
	Func: func(cmd *gcli.Command, args []string) error {
		cmd.AddArg("host", "宝塔地址", false)
		cmd.AddArg("key", "宝塔密钥", false)

		host := cmd.Arg("host").String()
		key := cmd.Arg("key").String()
		link := host + "/data?action=getData&table=sites"

		color.Infoln("地址：" + host)
		color.Infoln("密钥：" + key)

		result := utils.Post(link, utils.PatchSign(key, url.Values{}))
		color.Infoln(result)

		return nil
	},
}
