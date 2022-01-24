package bt

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

type Webname struct {
	Domain     string `json:"domain"`
	Domainlist string `json:"domainlist"`
	Count      int    `json:"count"`
}

var SiteCmd = &cobra.Command{
	Use:   "site",
	Short: "网站相关操作",
}

var show = &cobra.Command{
	Use:   "show",
	Short: "展示网站列表",
	Long:  `展示网站列表`,
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		key, _ := cmd.Flags().GetString("key")
		link := host + "/data?action=getData&table=sites"

		color.Infoln("地址：" + host)
		color.Infoln("密钥：" + key)

		result := Post(link, patchSign(key, url.Values{}))
		color.Infoln(result)
	},
}

var create = &cobra.Command{
	Use:   "create",
	Short: "创建网站",
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		key, _ := cmd.Flags().GetString("key")
		link := host + "/site?action=AddSite"

		color.Infoln("地址：" + host)
		color.Infoln("密钥：" + key)

		webname, err := json.Marshal(Webname{
			Domain:     "test.api4.top",
			Domainlist: "[]",
			Count:      0,
		})
		if err != nil {
			fmt.Println("生成json失败", err)
		}

		result := Post(link, patchSign(key, url.Values{
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

var delete = &cobra.Command{
	Use:   "delete",
	Short: "删除网站",
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		key, _ := cmd.Flags().GetString("key")
		link := host + "/site?action=DeleteSite"

		color.Infoln("地址：" + host)
		color.Infoln("密钥：" + key)

		result := Post(link, patchSign(key, url.Values{
			"webname": {"test.api4.top"},
			"id":      {"10"},
		}))
		color.Infoln(result)
	},
	Args: nil,
}

var types = &cobra.Command{
	Use:   "types",
	Short: "展示网站分类",
	Long:  `展示网站分类`,
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		key, _ := cmd.Flags().GetString("key")
		link := host + "/site?action=get_site_types"

		color.Infoln("地址：" + host)
		color.Infoln("密钥：" + key)

		result := Post(link, patchSign(key, url.Values{}))
		color.Infoln(result)
	},
}

var php = &cobra.Command{
	Use:   "php",
	Short: "展示PHP版本列表",
	Long:  `展示PHP版本列表`,
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		key, _ := cmd.Flags().GetString("key")
		link := host + "/site?action=GetPHPVersion"

		color.Infoln("地址：" + host)
		color.Infoln("密钥：" + key)

		result := Post(link, patchSign(key, url.Values{}))
		color.Infoln(result)
	},
}

func init() {
	SiteCmd.AddCommand(show)
	SiteCmd.AddCommand(types)
	SiteCmd.AddCommand(php)
	SiteCmd.AddCommand(create)
	SiteCmd.AddCommand(delete)

	SiteCmd.PersistentFlags().String("host", "http://127.0.0.1", "宝塔地址")
	SiteCmd.PersistentFlags().String("key", "", "宝塔密钥")
}
