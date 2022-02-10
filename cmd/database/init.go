package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gookit/gcli/v3"
)

var DatabaseCmd = &gcli.Command{
	Name: "database",
	Desc: "数据库相关操作",
	Subs: []*gcli.Command{createCmd, showCmd},
}
