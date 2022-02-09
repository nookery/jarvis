package cmd

import (
	"github.com/gookit/gcli/v3"
)

var Ping = &gcli.Command{
	Name:    "ping",
	Desc:    "输出<info>pang</>，用于测试。",
	Aliases: []string{"pi"},
	Func: func(cmd *gcli.Command, args []string) error {
		gcli.Println("pang")
		return nil
	},
}
