package main

import (
	"github.com/gookit/gcli/v3"
	"jarvis/cmd"
)

func main() {
	app := gcli.NewApp()
	app.Version = "1.0.3"
	app.Desc = "我是Jarvis，你的得力助理。"

	app.Add(cmd.Ping)
	app.Add(cmd.JokeCmd)
	app.Add(cmd.DatabaseCmd)

	app.Run(nil)
}
