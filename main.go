package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// 创建应用实例
	app := NewApp()

	// 创建应用窗口
	err := wails.Run(&options.App{
		Title:            "YouTube视频下载器",
		Width:            800,
		Height:           600,
		Assets:           assets,
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		OnStartup:        app.Startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
