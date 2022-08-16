package site

import (
	"encoding/json"
	"fmt"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"jarvis/cmd/bt/utils"
	"net/url"
)

type Webname struct {
	Domain     string `json:"domain"`
	Domainlist string `json:"domainlist"`
	Count      int    `json:"count"`
}

var Create = &cobra.Command{
	Use:   "create",
	Short: "创建网站",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		domain, _ := cmd.Flags().GetString("domain")
		path, _ := cmd.Flags().GetString("path")

		color.Infoln("域名：" + domain)
		color.Infoln("路径：" + path)

		// if domain == "" {
		// 	return errors.New(color.Error.Renderln("请输入网站域名") + "\r\n")
		// }
		// if path == "" {
		// 	return errors.New(color.Error.Renderln("请输入网站路径") + "\r\n")
		// }

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		key, _ := cmd.Flags().GetString("key")
		domain, _ := cmd.Flags().GetString("domain")
		comment, _ := cmd.Flags().GetString("comment")
		path, _ := cmd.Flags().GetString("path")
		link := host + "/site?action=AddSite"

		if path == "" {
			path = "/www/wwwroot/" + domain
		}

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
			"path":    {path},
			"type_id": {"0"},
			"type":    {"PHP"},
			"version": {"80"},
			"port":    {"80"},
			"ps":      {comment},
			"ftp":     {"false"},
			"sql":     {"false"},
		}))
		color.Infoln(result)
	},
}

func init() {
	Create.Flags().String("domain", "", color.Blue.Render("要新建的网站的域名"))
	Create.Flags().String("comment", "", color.Blue.Render("要新建的网站的备注"))
	Create.Flags().String("path", "", color.Blue.Render("要新建的网站的路径"))
	Create.MarkFlagRequired("domain")
}
