package crontab

import (
	"encoding/json"
	"errors"
	"jarvis/cmd/bt/utils"
	"net/url"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

type ErrorResponse struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
}

type CrontabItem struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func Get(host string, key string) ([]CrontabItem, error) {
	link := host + "/crontab?action=GetCrontab"

	var error ErrorResponse
	var items []CrontabItem

	result := utils.Post(link, utils.PatchSign(key, url.Values{}))

	json.Unmarshal([]byte(result), &error)
	json.Unmarshal([]byte(result), &items)

	if error.Status != "" {
		return nil, errors.New(color.Error.Renderln(error.Msg) + "\r\n")
	} else {
		return items, nil
	}
}

var get = &cobra.Command{
	Use:   "get",
	Short: "展示Crontab列表",
	Long:  color.Success.Render("展示Crontab列表"),
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		key, _ := cmd.Flags().GetString("key")

		items, err := Get(host, key)

		if err != nil {
			color.Errorln(err.Error())
		} else {
			for _, item := range items {
				color.Infoln(item.Id, utils.StrPadRight(item.Type, 16, " "), item.Name)
			}
		}
	},
}
