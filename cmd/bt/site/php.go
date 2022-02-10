package site

import (
	// "jarvis/cmd/bt/utils"
	// "net/url"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
)

var php = &gcli.Command{
	Name: "php",
	Desc: "展示PHP版本列表",
	// Config: func(cmd *gcli.Command) {
	// 	cmd.AddArg("key", "", true)
	// 	cmd.AddArg("host", "http://127.0.0.1", false)
	// },
	Func: func(cmd *gcli.Command, args []string) error {
		host := cmd.Parent().Parent().Arg("host").String()
		key := cmd.Arg("key").String()
		// link := host + "/site?action=GetPHPVersion"

		color.Infoln("地址：" + host)
		color.Infoln("密钥：" + key)

		// result := utils.Post(link, utils.PatchSign(key, url.Values{}))
		// color.Infoln(result)

		return nil
	},
}
