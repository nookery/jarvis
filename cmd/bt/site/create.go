package site

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"jarvis/cmd/bt/utils"
)

type Webname struct {
	Domain     string `json:"domain"`
	Domainlist string `json:"domainlist"`
	Count      int    `json:"count"`
}

var Create = &cobra.Command{
	Use:   "create",
	Short: "创建网站",
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		key, _ := cmd.Flags().GetString("key")
		domain, _ := cmd.Flags().GetString("domain")
		link := host + "/site?action=AddSite"

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
	},
}

func init() {
	Create.Flags().String("domain", "", "要新建的网站的域名")
	Create.MarkFlagRequired("domain")
}
