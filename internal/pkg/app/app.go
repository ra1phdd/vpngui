package app

import (
	"context"
	"embed"
	"fmt"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"vpngui/config"
	xray_api "vpngui/internal/app/xray-api"
	embedded "vpngui/pkg/embed"
	"vpngui/pkg/logger"
)

func New(assets embed.FS) error {
	err := config.LoadConfig()
	if err != nil {
		return err
	}

	logger.Init("info")

	err = embedded.Init()
	if err != nil {
		return err
	}

	app := NewApp()
	cfg := &config.Config{}

	err = wails.Run(&options.App{
		Title:  "VPN-GUI",
		Width:  400,
		Height: 550,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 30, G: 30, B: 80, A: 255},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
			xray_api.New(),
			cfg,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}

	return nil
}

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
