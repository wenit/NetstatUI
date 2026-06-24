package main

import (
	"embed"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wenit/NetstatUI/services/geo"
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
	application.RegisterEvent[monitor.GeoStatus]("geo:status")
}

func main() {
	switch runtime.GOOS {
	case "windows", "darwin", "linux":
	default:
		log.Fatalf("unsupported platform: %s", runtime.GOOS)
	}
	netstat.SetProvider(netstat.NewPlatformProvider())
	cache := process.NewCache()

	geoResolver, err := geo.New(geoDataDir())
	if err != nil {
		log.Printf("geo: init failed (geo lookup disabled): %v", err)
	} else {
		defer geoResolver.Close()
	}

	app := application.New(application.Options{
		Name:        "NetstatUI",
		Description: "NetstatUI",
		Services: []application.Service{
			application.NewService(NewAppService(cache)),
			application.NewService(monitor.New(cache, geoResolver)),
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

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

// geoDataDir returns <exe-dir>/data, the directory that holds the
// ip2region xdb files. The xdb files are no longer embedded in the
// binary; they are expected to sit next to the executable so users can
// swap or update them without rebuilding.
func geoDataDir() string {
	exe, err := os.Executable()
	if err != nil {
		return "data"
	}
	return filepath.Join(filepath.Dir(exe), "data")
}
