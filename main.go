package main

import (
	"jarvis/cmd"
	"jarvis/cmd/bt"
	"jarvis/cmd/database"

	"github.com/gookit/gcli/v3"
)

func main() {
	app := gcli.NewApp()
	app.Version = "1.0.3"
	app.Desc = "我是Jarvis，你的得力助理。"

	app.Add(cmd.Ping)
	app.Add(cmd.JokeCmd)
	app.Add(database.DatabaseCmd)
	app.Add(bt.BtCmd)

	app.Run(nil)
}
