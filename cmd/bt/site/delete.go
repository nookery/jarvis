package site

import (
	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"jarvis/cmd/bt/utils"
	"net/url"
)

var delete = &gcli.Command{
	Name: "delete",
	Desc: "删除网站",
	Config: func(cmd *gcli.Command) {
		cmd.AddArg("id", "Id")
	},
	Func: func(cmd *gcli.Command, args []string) error {
		host := cmd.Arg("host").String()
		key := cmd.Arg("key").String()
		link := host + "/site?action=DeleteSite"

		color.Infoln("地址：" + host)
		color.Infoln("密钥：" + key)

		result := utils.Post(link, utils.PatchSign(key, url.Values{
			"webname": {"test.api4.top"},
			"id":      {"10"},
		}))
		color.Infoln(result)

		return nil
	},
}
