package site

import (
	"encoding/json"
	"fmt"
	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"jarvis/cmd/bt/utils"
	"net/url"
)

type Webname struct {
	Domain     string `json:"domain"`
	Domainlist string `json:"domainlist"`
	Count      int    `json:"count"`
}

var create = &gcli.Command{
	Name: "create",
	Desc: "创建网站",
	Config: func(cmd *gcli.Command) {
		cmd.AddArg("key", "", true)
		cmd.AddArg("host", "宝塔地址", false)
		cmd.AddArg("domain", "网站域名")
	},
	Func: func(cmd *gcli.Command, args []string) error {
		host := cmd.Arg("host").String()
		key := cmd.Arg("key").String()
		domain := cmd.Arg("domain").String()
		link := host + "/site?action=AddSite"

		color.Infoln("地址：" + host)
		color.Infoln("密钥：" + key)

		webname, err := json.Marshal(Webname{
			Domain:     domain,
			Domainlist: "[]",
			Count:      0,
		})
		if err != nil {
			fmt.Println("生成json失败", err)
		}

		result := utils.Post(link, utils.PatchSign(key, url.Values{
			"webname": {string(webname)},
			"path":    {"/www/wwwroot/test"},
			"type_id": {"0"},
			"type":    {"PHP"},
			"version": {"80"},
			"port":    {"80"},
			"ps":      {"仅用于测试"},
			"ftp":     {"false"},
			"sql":     {"false"},
		}))
		color.Infoln(result)

		return nil
	},
}
