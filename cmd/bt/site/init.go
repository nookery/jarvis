package site

import (
	"github.com/gookit/gcli/v3"
)

var SiteCmd = &gcli.Command{
	Name: "site",
	Desc: "网站相关操作",
	Subs: []*gcli.Command{
		show, types, php, delete, create,
	},
}
