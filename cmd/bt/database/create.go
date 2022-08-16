package database

import (
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
	Short: "创建数据库",
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		key, _ := cmd.Flags().GetString("key")

		name, _ := cmd.Flags().GetString("name")
		user, _ := cmd.Flags().GetString("user")
		password, _ := cmd.Flags().GetString("password")
		link := host + "/database?action=AddDatabase"

		result := utils.Post(link, utils.PatchSign(key, url.Values{
			"name":           {name},
			"db_user":        {user},
			"password":       {password},
			"databaseAccess": {"127.0.0.1"},
			"address":        {"127.0.0.1"},
			"ps":             {name},
			"dtype":          {"MySQL"},
			"codeing":        {"utf8"},
		}))
		color.Infoln(result)
	},
}

func init() {
	Create.Flags().String("name", "", color.Blue.Render("数据库名称"))
	Create.Flags().String("user", "", color.Blue.Render("用户名"))
	Create.Flags().String("password", "", color.Blue.Render("密码"))
	Create.MarkFlagRequired("name")
}
