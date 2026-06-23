package main

import (
	"embed"
	"log"
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wenit/NetstatUI/services/monitor"
	"github.com/wenit/NetstatUI/services/netstat"
	"github.com/wenit/NetstatUI/services/process"
)

//go:embed all:frontend/dist
var assets embed.FS

func init() {
	application.RegisterEvent[[]netstat.ConnInfo]("conn:full")
	application.RegisterEvent[monitor.Diff]("conn:diff")
	application.RegisterEvent[monitor.Stats]("conn:stats")
	application.RegisterEvent[string]("conn:error")
}

func main() {
	if runtime.GOOS != "windows" && runtime.GOOS != "darwin" {
		log.Fatalf("unsupported platform: %s", runtime.GOOS)
	}
	netstat.SetProvider(netstat.NewPlatformProvider())
	cache := process.NewCache()

	app := application.New(application.Options{
		Name:        "NetstatUI",
		Description: "A graphical user interface for the netstat command",
		Services: []application.Service{
			application.NewService(NewAppService(cache)),
			application.NewService(monitor.New(cache)),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:          "NetstatUI",
		Frameless:      true,
		Width:          1200,
		Height:         760,
		MinWidth:       860,
		MinHeight:      520,
		BackgroundType: application.BackgroundTypeTranslucent,
		Windows: application.WindowsWindow{
			BackdropType: application.Mica,
		},
		URL: "/",
	})

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
