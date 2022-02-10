package bt

import (
	"jarvis/cmd/bt/site"

	"github.com/gookit/gcli/v3"
)

var BtCmd = &gcli.Command{
	Name: "bt",
	Desc: "宝塔相关操作",
	Subs: []*gcli.Command{site.SiteCmd},
}
